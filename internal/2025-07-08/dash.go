package main

import (
   "encoding/xml"
   "fmt"
   "io"
   "net/url"
   "os"
   "path/filepath" // Added for filepath.Abs and filepath.Dir
   "strconv"
   "strings"
)

// MPD represents the root element of a DASH MPD file.
type MPD struct {
   XMLName           xml.Name `xml:"MPD"`
   BaseURL           string   `xml:"BaseURL"`
   Periods           []Period `xml:"Period"`
   MediaPresentationDuration string `xml:"mediaPresentationDuration,attr"` // e.g., PT0H0M30.0S
}

// Period represents a Period element within the MPD.
type Period struct {
   XMLName        xml.Name         `xml:"Period"`
   AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
   Duration       string           `xml:"duration,attr"` // e.g., PT0H0M10.0S
   BaseURL        string           `xml:"BaseURL"` // BaseURL can also be at Period level
}

// AdaptationSet represents an AdaptationSet element.
type AdaptationSet struct {
   XMLName        xml.Name         `xml:"AdaptationSet"`
   MimeType       string           `xml:"mimeType,attr"`
   // SegmentTemplate can be a child of AdaptationSet
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"` // Use pointer as it's optional
   Representations []Representation `xml:"Representation"`
   BaseURL        string           `xml:"BaseURL"` // BaseURL can also be at AdaptationSet level
}

// Representation represents a Representation element.
type Representation struct {
   XMLName        xml.Name        `xml:"Representation"`
   ID             string          `xml:"id,attr"`
   Bandwidth      int             `xml:"bandwidth,attr"`
   // SegmentTemplate can also be a child of Representation (overrides AdaptationSet's if present)
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"` // Use pointer as it's optional
   BaseURL        string           `xml:"BaseURL"` // BaseURL can also be at Representation level
}

// SegmentTemplate represents a SegmentTemplate element.
type SegmentTemplate struct {
   XMLName        xml.Name `xml:"SegmentTemplate"`
   Media          string   `xml:"media,attr"`
   Initialization string   `xml:"initialization,attr"`
   StartNumber    int      `xml:"startNumber,attr"`
   Timescale      int      `xml:"timescale,attr"`
   Duration       int      `xml:"duration,attr"` // Duration of each segment in timescale units
   SegmentTimeline *SegmentTimeline `xml:"SegmentTimeline"` // Optional, for variable duration segments
}

// SegmentTimeline represents a SegmentTimeline element for variable duration segments.
type SegmentTimeline struct {
   XMLName xml.Name `xml:"SegmentTimeline"`
   Ss      []S      `xml:"S"`
}

// S represents an S element within SegmentTimeline.
type S struct {
   XMLName xml.Name `xml:"S"`
   T       int      `xml:"t,attr"` // Start time of the segment (optional, if omitted, it's continuous from previous)
   D       int      `xml:"d,attr"` // Duration of the segment
   R       int      `xml:"r,attr"` // Repeat count (optional, default 0 means 1 segment)
}


func main() {
   if len(os.Args) < 2 {
      fmt.Println("Usage: go run mpd_parser.go <path_to_mpd_file> [mpd_url]")
      os.Exit(1)
   }

   mpdFilePath := os.Args[1]
   var mpdURL string
   if len(os.Args) > 2 {
      mpdURL = os.Args[2]
   }

   // Read the XML file
   xmlFile, err := os.Open(mpdFilePath)
   if err != nil {
      fmt.Printf("Error opening file: %v\n", err)
      return
   }
   defer xmlFile.Close()

   byteValue, _ := io.ReadAll(xmlFile)

   var mpd MPD
   err = xml.Unmarshal(byteValue, &mpd)
   if err != nil {
      fmt.Printf("Error unmarshalling XML: %v\n", err)
      return
   }

   // Determine the initial effective BaseURL
   var initialEffectiveBaseURL string
   if mpdURL != "" {
      // If an MPD URL is provided, use its directory as the initial base
      parsedMPDURL, err := url.Parse(mpdURL)
      if err != nil {
         fmt.Printf("Warning: Could not parse provided MPD URL '%s': %v. Falling back to local file path.\n", mpdURL, err)
         // Fallback to local file path if MPD URL is invalid
         absPath, _ := filepath.Abs(mpdFilePath)
         initialEffectiveBaseURL = "file://" + filepath.ToSlash(filepath.Dir(absPath)) + "/"
      } else {
         // Get the directory of the MPD URL. ResolveReference with "." ensures we get the directory path.
         initialEffectiveBaseURL = parsedMPDURL.ResolveReference(&url.URL{Path: "."}).String()
         // Ensure it ends with a slash if it's a directory
         if !strings.HasSuffix(initialEffectiveBaseURL, "/") {
            initialEffectiveBaseURL += "/"
         }
      }
   } else {
      // If no MPD URL, use the local file's directory as the base, made absolute with "file://"
      absPath, err := filepath.Abs(mpdFilePath)
      if err != nil {
         fmt.Printf("Error getting absolute path for MPD file: %v\n", err)
         return
      }
      initialEffectiveBaseURL = "file://" + filepath.ToSlash(filepath.Dir(absPath)) + "/"
   }

   effectiveBaseURL := initialEffectiveBaseURL

   fmt.Println("--- Extracted Segment URLs ---")

   for _, period := range mpd.Periods {
      // Update effective BaseURL if defined at Period level
      if period.BaseURL != "" {
         effectiveBaseURL = resolveURL(effectiveBaseURL, period.BaseURL)
      }

      periodDurationSeconds := parseDuration(period.Duration)
      if periodDurationSeconds == 0 {
         periodDurationSeconds = parseDuration(mpd.MediaPresentationDuration)
      }

      for _, as := range period.AdaptationSets {
         // Update effective BaseURL if defined at AdaptationSet level
         if as.BaseURL != "" {
            effectiveBaseURL = resolveURL(effectiveBaseURL, as.BaseURL)
         }

         // Get the SegmentTemplate defined at the AdaptationSet level
         asSegmentTemplate := as.SegmentTemplate

         for _, rep := range as.Representations {
            // Update effective BaseURL if defined at Representation level
            if rep.BaseURL != "" {
               effectiveBaseURL = resolveURL(effectiveBaseURL, rep.BaseURL)
            }

            // Determine the effective SegmentTemplate for this Representation.
            // A SegmentTemplate at the Representation level overrides one at the AdaptationSet level.
            var currentST *SegmentTemplate
            if rep.SegmentTemplate != nil {
               currentST = rep.SegmentTemplate
            } else {
               currentST = asSegmentTemplate
            }

            // If no SegmentTemplate is found at either level, skip this representation
            if currentST == nil {
               fmt.Printf("Warning: Representation ID %s: No SegmentTemplate found at AdaptationSet or Representation level. Skipping.\n", rep.ID)
               continue
            }

            // Construct initialization URL specific to this representation
            if currentST.Initialization != "" {
               initURL := strings.ReplaceAll(currentST.Initialization, "$RepresentationID$", rep.ID)
               // No $Time$ or $Number$ in initialization URL
               fullInitURL := resolveURL(effectiveBaseURL, initURL)
               fmt.Printf("Initialization URL (Rep ID: %s): %s\n", rep.ID, fullInitURL)
            }

            if currentST.SegmentTimeline != nil && len(currentST.SegmentTimeline.Ss) > 0 {
               // Handle SegmentTimeline (variable duration segments)
               fmt.Printf("Segments for Representation ID: %s (Bandwidth: %d):\n", rep.ID, rep.Bandwidth)
               
               var currentTime int // Time in timescale units
               // If the first S element has a 't' attribute, use it as the starting time.
               // Otherwise, start from 0 for the beginning of the period.
               if currentST.SegmentTimeline.Ss[0].T != 0 {
                  currentTime = currentST.SegmentTimeline.Ss[0].T
               } else {
                  currentTime = 0
               }

               segmentCounter := currentST.StartNumber // For $Number$ placeholder, if used
               for _, s := range currentST.SegmentTimeline.Ss {
                  // If 't' is specified for this S element, it overrides the running currentTime
                  if s.T != 0 {
                     currentTime = s.T
                  }

                  // Calculate the actual number of segments in this run (1 + r)
                  runSegments := 1
                  if s.R > 0 {
                     runSegments = 1 + s.R
                  }
                  
                  for i := 0; i < runSegments; i++ {
                     mediaURLTemplate := strings.ReplaceAll(currentST.Media, "$RepresentationID$", rep.ID)
                     mediaURLTemplate = strings.ReplaceAll(mediaURLTemplate, "$Bandwidth$", strconv.Itoa(rep.Bandwidth))
                     
                     // Replace $Time$ with current time
                     mediaURLTemplate = strings.ReplaceAll(mediaURLTemplate, "$Time$", strconv.Itoa(currentTime))
                     
                     // Replace $Number$ if it exists (less common with $Time$, but for robustness)
                     mediaURLTemplate = strings.ReplaceAll(mediaURLTemplate, "$Number$", strconv.Itoa(segmentCounter))

                     fullSegmentURL := resolveURL(effectiveBaseURL, mediaURLTemplate)
                     fmt.Printf("  Segment (Time: %d, Num: %d): %s\n", currentTime, segmentCounter, fullSegmentURL)
                     
                     currentTime += s.D // Advance time for the next segment in this run
                     segmentCounter++
                  }
               }
            } else {
               // Handle fixed duration segments using StartNumber and Duration
               if currentST.Duration == 0 || currentST.Timescale == 0 {
                  fmt.Printf("Warning: Representation ID %s: Effective SegmentTemplate has no SegmentTimeline and is missing Duration or Timescale. Cannot calculate fixed-duration segments.\n", rep.ID)
                  continue
               }

               segmentDurationSeconds := float64(currentST.Duration) / float64(currentST.Timescale)
               numSegments := 0
               if periodDurationSeconds > 0 {
                  numSegments = int(periodDurationSeconds / segmentDurationSeconds)
               } else {
                  fmt.Printf("Warning: Period duration not found for Representation ID: %s. Assuming 10 segments.\n", rep.ID)
                  numSegments = 10 // Arbitrary default
               }

               fmt.Printf("Segments for Representation ID: %s (Bandwidth: %d):\n", rep.ID, rep.Bandwidth)
               for i := 0; i < numSegments; i++ {
                  segmentNumber := currentST.StartNumber + i
                  mediaURL := strings.ReplaceAll(currentST.Media, "$RepresentationID$", rep.ID)
                  mediaURL = strings.ReplaceAll(mediaURL, "$Bandwidth$", strconv.Itoa(rep.Bandwidth))
                  mediaURL = strings.ReplaceAll(mediaURL, "$Number$", strconv.Itoa(segmentNumber))
                  
                  fullSegmentURL := resolveURL(effectiveBaseURL, mediaURL)
                  fmt.Printf("  Segment %d: %s\n", segmentNumber, fullSegmentURL)
               }
            }
         }
      }
   }
}

// resolveURL combines a base URL (which should be absolute) and a relative path to form a full URL.
func resolveURL(base, relative string) string {
   // If the relative path is already an absolute URL (has a scheme), return it directly.
   if strings.Contains(relative, "://") {
      return relative
   }
   baseURL, err := url.Parse(base)
   if err != nil {
      fmt.Printf("Warning: Could not parse base URL '%s': %v. Returning relative path.\n", base, err)
      return relative
   }
   relativeURL, err := url.Parse(relative)
   if err != nil {
      fmt.Printf("Warning: Could not parse relative URL '%s': %v. Returning relative path.\n", relative, err)
      return relative
   }
   // Resolve the reference. net/url.ResolveReference is designed to handle this.
   return baseURL.ResolveReference(relativeURL).String()
}

// parseDuration parses an ISO 8601 duration string (e.g., PT0H0M30.0S) into seconds.
// This is a simplified parser and might not handle all edge cases of ISO 8601 durations.
func parseDuration(duration string) float64 {
   if duration == "" {
      return 0
   }
   
   duration = strings.TrimPrefix(duration, "PT")
   var totalSeconds float64

   if strings.Contains(duration, "H") {
      parts := strings.Split(duration, "H")
      hours, _ := strconv.ParseFloat(parts[0], 64)
      totalSeconds += hours * 3600
      duration = parts[1]
   }
   if strings.Contains(duration, "M") {
      parts := strings.Split(duration, "M")
      minutes, _ := strconv.ParseFloat(parts[0], 64)
      totalSeconds += minutes * 60
      duration = parts[1]
   }
   if strings.Contains(duration, "S") {
      parts := strings.Split(duration, "S")
      seconds, _ := strconv.ParseFloat(parts[0], 64)
      totalSeconds += seconds
   }
   return totalSeconds
}

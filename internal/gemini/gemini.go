package main

import (
   "encoding/json"
   "encoding/xml"
   "fmt"
   "io/ioutil"
   "math"
   "net/url"
   "os"
   "strconv"
   "strings"
)

// Define structs to parse the MPD XML
type MPD struct {
   XMLName xml.Name `xml:"MPD"`
   BaseURL string   `xml:"BaseURL,omitempty"`
   Periods []Period `xml:"Period"`
}

type Period struct {
   XMLName        xml.Name        `xml:"Period"`
   BaseURL        string          `xml:"BaseURL,omitempty"`
   Duration       string          `xml:"duration,attr"` // Period duration in ISO 8601 format
   AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
}

type AdaptationSet struct {
   XMLName         xml.Name         `xml:"AdaptationSet"`
   BaseURL         string           `xml:"BaseURL,omitempty"`
   Representations []Representation `xml:"Representation"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
}

type Representation struct {
   XMLName         xml.Name         `xml:"Representation"`
   ID              string           `xml:"id,attr"`
   BaseURL         string           `xml:"BaseURL,omitempty"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
}

type SegmentTemplate struct {
   XMLName                xml.Name         `xml:"SegmentTemplate"`
   Timescale              uint64           `xml:"timescale,attr"`
   Initialization         string           `xml:"initialization,attr"`
   Media                  string           `xml:"media,attr"`
   StartNumber            uint64           `xml:"startNumber,attr"`
   Duration               uint64           `xml:"duration,attr"`
   EndNumber              uint64           `xml:"endNumber,attr"`
   PresentationTimeOffset uint64           `xml:"presentationTimeOffset,attr"`
   SegmentTimeline        *SegmentTimeline `xml:"SegmentTimeline"`
}

type SegmentTimeline struct {
   XMLName xml.Name `xml:"SegmentTimeline"` // Corrected: Should match the struct name for the XML tag
   S       []S      `xml:"S"`
}

type S struct {
   XMLName xml.Name `xml:"S"`
   T       uint64   `xml:"t,attr"` // Start time
   D       uint64   `xml:"d,attr"` // Duration
   R       int64    `xml:"r,attr"` // Repeat count
}

// parseDuration parses an ISO 8601 duration string (e.g., PT1H51M39.193S) into seconds.
func parseDuration(durationStr string) (float64, error) {
   if !strings.HasPrefix(durationStr, "PT") {
      return 0, fmt.Errorf("invalid ISO 8601 duration format: %s", durationStr)
   }

   durationStr = durationStr[2:]
   var totalSeconds float64

   // Simplified manual parse for ISO 8601 duration without regex.
   // This approach expects H, M, S components to be clearly delimited by their letters.
   parts := make(map[string]float64)
   var currentNum string
   for _, r := range durationStr {
      if r >= '0' && r <= '9' || r == '.' {
         currentNum += string(r)
      } else {
         if currentNum != "" {
            val, _ := strconv.ParseFloat(currentNum, 64)
            parts[string(r)] = val
            currentNum = ""
         }
      }
   }
   if currentNum != "" { // Last part (seconds)
      val, _ := strconv.ParseFloat(currentNum, 64)
      parts["S"] = val
   }

   totalSeconds += parts["H"] * 3600
   totalSeconds += parts["M"] * 60
   totalSeconds += parts["S"]

   return totalSeconds, nil
}

func main() {
   if len(os.Args) < 2 {
      fmt.Println("Usage: go run main.go <path_to_mpd_file>")
      os.Exit(1)
   }

   mpdFilePath := os.Args[1]
   mpdContent, err := ioutil.ReadFile(mpdFilePath)
   if err != nil {
      fmt.Printf("Error reading MPD file: %v\n", err)
      os.Exit(1)
   }

   var mpd MPD
   err = xml.Unmarshal(mpdContent, &mpd)
   if err != nil {
      fmt.Printf("Error unmarshalling MPD XML: %v\n", err)
      os.Exit(1)
   }

   mpdBaseURL, err := url.Parse("http://test.test/test.mpd")
   if err != nil {
      fmt.Printf("Error parsing base MPD URL: %v\n", err)
      os.Exit(1)
   }

   representationSegments := make(map[string][]string)

   currentBaseURL := mpdBaseURL
   if mpd.BaseURL != "" {
      parsedBaseURL, err := url.Parse(mpd.BaseURL)
      if err != nil {
         fmt.Printf("Warning: Could not parse MPD BaseURL '%s': %v\n", mpd.BaseURL, err)
      } else {
         currentBaseURL = mpdBaseURL.ResolveReference(parsedBaseURL)
      }
   }

   for _, period := range mpd.Periods {
      periodBaseURL := currentBaseURL
      if period.BaseURL != "" {
         parsedBaseURL, err := url.Parse(period.BaseURL)
         if err != nil {
            fmt.Printf("Warning: Could not parse Period BaseURL '%s': %v\n", period.BaseURL, err)
         } else {
            periodBaseURL = currentBaseURL.ResolveReference(parsedBaseURL)
         }
      }

      for _, as := range period.AdaptationSets {
         asBaseURL := periodBaseURL
         if as.BaseURL != "" {
            parsedBaseURL, err := url.Parse(as.BaseURL)
            if err != nil {
               fmt.Printf("Warning: Could not parse AdaptationSet BaseURL '%s': %v\n", as.BaseURL, err)
            } else {
               asBaseURL = periodBaseURL.ResolveReference(parsedBaseURL)
            }
         }

         for _, rep := range as.Representations {
            // Added warning if Representation ID is empty
            if rep.ID == "" {
               fmt.Printf("Warning: Encountered a Representation with an empty 'id' attribute (or missing). Placeholder $RepresentationID$ will be replaced by an empty string, leading to a double hyphen in the URL.\n")
            }

            repBaseURL := asBaseURL
            if rep.BaseURL != "" {
               parsedBaseURL, err := url.Parse(rep.BaseURL)
               if err != nil {
                  fmt.Printf("Warning: Could not parse Representation BaseURL '%s': %v\n", rep.BaseURL, err)
               } else {
                  repBaseURL = asBaseURL.ResolveReference(parsedBaseURL)
               }
            }

            var currentRepSegments []string

            var effectiveSegmentTemplate *SegmentTemplate
            if rep.SegmentTemplate != nil {
               effectiveSegmentTemplate = rep.SegmentTemplate
            } else {
               effectiveSegmentTemplate = as.SegmentTemplate
            }

            if rep.BaseURL != "" && effectiveSegmentTemplate == nil {
               currentRepSegments = append(currentRepSegments, repBaseURL.String())
            } else if effectiveSegmentTemplate != nil {
               st := effectiveSegmentTemplate

               if st.SegmentTimeline != nil {
                  currentTime := st.PresentationTimeOffset
                  if st.Timescale == 0 {
                     st.Timescale = 1
                  }

                  currentSegmentNumber := st.StartNumber
                  if currentSegmentNumber == 0 {
                     currentSegmentNumber = 1
                  }

                  for _, s := range st.SegmentTimeline.S {
                     if s.T > 0 {
                        currentTime = s.T
                     }

                     numRepeats := s.R + 1
                     for i := 0; i < int(numRepeats); i++ {
                        mediaPath := st.Media

                        // Explicit hardcoded replacements for $Number%0xd$ (no loops or switches)
                        mediaPath = strings.Replace(mediaPath, "$Number%09d$", fmt.Sprintf("%09d", currentSegmentNumber), -1)
                        mediaPath = strings.Replace(mediaPath, "$Number%08d$", fmt.Sprintf("%08d", currentSegmentNumber), -1)
                        mediaPath = strings.Replace(mediaPath, "$Number%07d$", fmt.Sprintf("%07d", currentSegmentNumber), -1)
                        mediaPath = strings.Replace(mediaPath, "$Number%06d$", fmt.Sprintf("%06d", currentSegmentNumber), -1)
                        mediaPath = strings.Replace(mediaPath, "$Number%05d$", fmt.Sprintf("%05d", currentSegmentNumber), -1)
                        mediaPath = strings.Replace(mediaPath, "$Number%04d$", fmt.Sprintf("%04d", currentSegmentNumber), -1)
                        mediaPath = strings.Replace(mediaPath, "$Number%03d$", fmt.Sprintf("%03d", currentSegmentNumber), -1)
                        mediaPath = strings.Replace(mediaPath, "$Number%02d$", fmt.Sprintf("%02d", currentSegmentNumber), -1)

                        // Handle generic $Number$
                        mediaPath = strings.Replace(mediaPath, "$Number$", fmt.Sprintf("%d", currentSegmentNumber), -1)

                        // Handle generic $Time$
                        mediaPath = strings.Replace(mediaPath, "$Time$", fmt.Sprintf("%d", currentTime), -1)

                        // Handle $RepresentationID$
                        mediaPath = strings.Replace(mediaPath, "$RepresentationID$", rep.ID, -1)

                        mediaURL, err := url.Parse(mediaPath)
                        if err != nil {
                           fmt.Printf("Warning: Could not parse media path '%s': %v\n", mediaPath, err)
                           continue
                        }
                        currentRepSegments = append(currentRepSegments, repBaseURL.ResolveReference(mediaURL).String())
                        currentTime += s.D
                        currentSegmentNumber++
                     }
                  }
               } else { // SegmentTemplate without SegmentTimeline
                  startNumber := st.StartNumber
                  if startNumber == 0 {
                     startNumber = 1
                  }

                  var endNumber uint64
                  if st.EndNumber > 0 {
                     endNumber = st.EndNumber
                  } else {
                     periodDurationSeconds := 0.0
                     if period.Duration != "" {
                        parsedPeriodDuration, err := parseDuration(period.Duration)
                        if err == nil {
                           periodDurationSeconds = parsedPeriodDuration
                        } else {
                           fmt.Printf("Warning: Could not parse Period duration '%s': %v. Falling back to default segment count for %s.\n", period.Duration, err, rep.ID)
                        }
                     }

                     segmentDurationInSeconds := float64(st.Duration) / float64(st.Timescale)

                     if segmentDurationInSeconds > 0 && periodDurationSeconds > 0 {
                        segmentCountFloat := math.Ceil(periodDurationSeconds / segmentDurationInSeconds)
                        segmentCount := uint64(segmentCountFloat)
                        endNumber = startNumber + segmentCount - 1
                     } else {
                        endNumber = startNumber + 10
                     }
                  }

                  if st.Timescale == 0 {
                     st.Timescale = 1
                  }

                  for i := startNumber; i <= endNumber; i++ {
                     mediaPath := st.Media

                     // Explicit hardcoded replacements for $Number%0xd$ (no loops or switches)
                     mediaPath = strings.Replace(mediaPath, "$Number%09d$", fmt.Sprintf("%09d", i), -1)
                     mediaPath = strings.Replace(mediaPath, "$Number%08d$", fmt.Sprintf("%08d", i), -1)
                     mediaPath = strings.Replace(mediaPath, "$Number%07d$", fmt.Sprintf("%07d", i), -1)
                     mediaPath = strings.Replace(mediaPath, "$Number%06d$", fmt.Sprintf("%06d", i), -1)
                     mediaPath = strings.Replace(mediaPath, "$Number%05d$", fmt.Sprintf("%05d", i), -1)
                     mediaPath = strings.Replace(mediaPath, "$Number%04d$", fmt.Sprintf("%04d", i), -1)
                     mediaPath = strings.Replace(mediaPath, "$Number%03d$", fmt.Sprintf("%03d", i), -1)
                     mediaPath = strings.Replace(mediaPath, "$Number%02d$", fmt.Sprintf("%02d", i), -1)

                     // Handle generic $Number$
                     mediaPath = strings.Replace(mediaPath, "$Number$", fmt.Sprintf("%d", i), -1)

                     // Handle generic $Time$
                     mediaPath = strings.Replace(mediaPath, "$Time$", fmt.Sprintf("%d", st.PresentationTimeOffset), -1)

                     // Handle $RepresentationID$
                     mediaPath = strings.Replace(mediaPath, "$RepresentationID$", rep.ID, -1)

                     mediaURL, err := url.Parse(mediaPath)
                     if err != nil {
                        fmt.Printf("Warning: Could not parse media path '%s': %v\n", mediaPath, err)
                        continue
                     }
                     currentRepSegments = append(currentRepSegments, repBaseURL.ResolveReference(mediaURL).String())
                  }
               }
            }

            if len(currentRepSegments) > 0 {
               representationSegments[rep.ID] = append(representationSegments[rep.ID], currentRepSegments...)
            }
         }
      }
   }

   jsonOutput, err := json.MarshalIndent(representationSegments, "", "  ")
   if err != nil {
      fmt.Printf("Error marshalling JSON: %v\n", err)
      os.Exit(1)
   }

   fmt.Println(string(jsonOutput))
}

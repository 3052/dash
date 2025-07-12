package main

import (
   "encoding/json"
   "encoding/xml"
   "flag"
   "fmt"
   "io"
   "net/url"
   "os"
   "path"
   "strings"
)

type MPD struct {
   XMLName xml.Name `xml:"MPD"`
   Period  []Period `xml:"Period"`
   BaseURL string   `xml:"BaseURL"`
}

type Period struct {
   AdaptationSet []AdaptationSet `xml:"AdaptationSet"`
   BaseURL       string          `xml:"BaseURL"`
}

type AdaptationSet struct {
   Representation  []Representation `xml:"Representation"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
}

type Representation struct {
   ID              string           `xml:"id,attr"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   BaseURL         string           `xml:"BaseURL"`
}

type SegmentTemplate struct {
   Media           string           `xml:"media,attr"`
   Initialization  string           `xml:"initialization,attr"`
   Duration        int              `xml:"duration,attr"`
   Timescale       int              `xml:"timescale,attr"`
   StartNumber     int              `xml:"startNumber,attr"`
   SegmentTimeline *SegmentTimeline `xml:"SegmentTimeline"`
}

type SegmentTimeline struct {
   S []Segment `xml:"S"`
}

type Segment struct {
   T uint64 `xml:"t,attr"` // Start time
   D uint64 `xml:"d,attr"` // Duration
   R int    `xml:"r,attr"` // Repeat count
}

type Output struct {
   Representations map[string][]string `json:"representations"`
}

func main() {
   mpdPath := flag.String("mpd", "", "Path to the MPD file")
   flag.Parse()

   if *mpdPath == "" {
      fmt.Fprintln(os.Stderr, "Please provide the path to the MPD file using -mpd flag")
      os.Exit(1)
   }

   mpdURL := "http://test.test/test.mpd"
   output, err := parseMPD(*mpdPath, mpdURL)
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error: %v\n", err)
      os.Exit(1)
   }

   jsonOutput, err := json.MarshalIndent(output, "", "  ")
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
      os.Exit(1)
   }

   fmt.Println(string(jsonOutput))
}

func parseMPD(mpdPath, mpdURL string) (Output, error) {
   file, err := os.Open(mpdPath)
   if err != nil {
      return Output{}, fmt.Errorf("failed to open MPD file: %w", err)
   }
   defer file.Close()

   data, err := io.ReadAll(file)
   if err != nil {
      return Output{}, fmt.Errorf("failed to read MPD file: %w", err)
   }

   var mpd MPD
   if err := xml.Unmarshal(data, &mpd); err != nil {
      return Output{}, fmt.Errorf("failed to parse MPD XML: %w", err)
   }

   fmt.Fprintf(os.Stderr, "Debug: Parsed MPD with %d periods\n", len(mpd.Period))

   baseURL, err := url.Parse(mpdURL)
   if err != nil {
      return Output{}, fmt.Errorf("failed to parse MPD URL: %w", err)
   }
   if mpd.BaseURL != "" {
      relativeBase, err := url.Parse(mpd.BaseURL)
      if err != nil {
         return Output{}, fmt.Errorf("failed to parse MPD BaseURL: %w", err)
      }
      baseURL = baseURL.ResolveReference(relativeBase)
      fmt.Fprintf(os.Stderr, "Debug: MPD BaseURL: %s\n", baseURL.String())
   }

   output := Output{Representations: make(map[string][]string)}
   for _, period := range mpd.Period {
      periodBaseURL := baseURL
      if period.BaseURL != "" {
         relativeBase, err := url.Parse(period.BaseURL)
         if err != nil {
            fmt.Fprintf(os.Stderr, "Warning: failed to parse Period BaseURL: %v\n", err)
            continue
         }
         periodBaseURL = baseURL.ResolveReference(relativeBase)
         fmt.Fprintf(os.Stderr, "Debug: Period BaseURL: %s\n", periodBaseURL.String())
      }

      for _, adaptationSet := range period.AdaptationSet {
         fmt.Fprintf(os.Stderr, "Debug: Processing AdaptationSet with %d representations\n", len(adaptationSet.Representation))
         if adaptationSet.SegmentTemplate != nil {
            fmt.Fprintf(os.Stderr, "Debug: Found SegmentTemplate at AdaptationSet level\n")
         }

         for _, rep := range adaptationSet.Representation {
            segmentURLs := []string{}
            repBaseURL := periodBaseURL
            if rep.BaseURL != "" && rep.SegmentTemplate == nil {
               // For representations with BaseURL only (e.g., subtitles), resolve directly against MPD base
               resolvedURL, err := resolveURL(baseURL, rep.BaseURL)
               if err != nil {
                  fmt.Fprintf(os.Stderr, "Warning: failed to resolve BaseURL for %s: %v\n", rep.ID, err)
                  continue
               }
               fmt.Fprintf(os.Stderr, "Debug: Generated BaseURL for %s: %s\n", rep.ID, resolvedURL)
               segmentURLs = append(segmentURLs, resolvedURL)
            } else {
               if rep.BaseURL != "" {
                  relativeBase, err := url.Parse(rep.BaseURL)
                  if err != nil {
                     fmt.Fprintf(os.Stderr, "Warning: failed to parse Representation BaseURL for %s: %v\n", rep.ID, err)
                     continue
                  }
                  repBaseURL = periodBaseURL.ResolveReference(relativeBase)
                  fmt.Fprintf(os.Stderr, "Debug: Representation %s BaseURL: %s\n", rep.ID, repBaseURL.String())
               }

               // Use Representation SegmentTemplate if available, else fall back to AdaptationSet
               segTemplate := rep.SegmentTemplate
               if segTemplate == nil {
                  segTemplate = adaptationSet.SegmentTemplate
               }

               if segTemplate != nil {
                  fmt.Fprintf(os.Stderr, "Debug: Found SegmentTemplate for %s\n", rep.ID)
                  // Handle initialization segment
                  if segTemplate.Initialization != "" {
                     initTemplate := strings.ReplaceAll(segTemplate.Initialization, "$RepresentationID$", rep.ID)
                     initURL, err := resolveURL(repBaseURL, initTemplate)
                     if err != nil {
                        fmt.Fprintf(os.Stderr, "Warning: failed to resolve initialization URL for %s: %v\n", rep.ID, err)
                        continue
                     }
                     fmt.Fprintf(os.Stderr, "Debug: Generated init URL for %s: %s\n", rep.ID, initURL)
                     segmentURLs = append(segmentURLs, initURL)
                  }
                  // Handle media segments
                  if segTemplate.Media != "" {
                     mediaTemplate := strings.ReplaceAll(segTemplate.Media, "$RepresentationID$", rep.ID)
                     if strings.Contains(mediaTemplate, "$Time$") && segTemplate.SegmentTimeline != nil {
                        fmt.Fprintf(os.Stderr, "Debug: Processing $Time$ segments for %s\n", rep.ID)
                        currentTime := uint64(0)
                        for _, segment := range segTemplate.SegmentTimeline.S {
                           count := 1
                           if segment.R > 0 {
                              count = segment.R + 1
                           }
                           for i := 0; i < count; i++ {
                              mediaURL := strings.ReplaceAll(mediaTemplate, "$Time$", fmt.Sprintf("%d", currentTime))
                              resolvedURL, err := resolveURL(repBaseURL, mediaURL)
                              if err != nil {
                                 fmt.Fprintf(os.Stderr, "Warning: failed to resolve media URL for %s at time %d: %v\n", rep.ID, currentTime, err)
                                 continue
                              }
                              fmt.Fprintf(os.Stderr, "Debug: Generated segment URL for %s: %s\n", rep.ID, resolvedURL)
                              segmentURLs = append(segmentURLs, resolvedURL)
                              currentTime += segment.D
                           }
                        }
                     } else if strings.Contains(mediaTemplate, "$Number$") && segTemplate.SegmentTimeline != nil {
                        fmt.Fprintf(os.Stderr, "Debug: Processing $Number$ segments for %s\n", rep.ID)
                        // Calculate total segments from SegmentTimeline
                        totalSegments := 0
                        for _, segment := range segTemplate.SegmentTimeline.S {
                           count := 1
                           if segment.R > 0 {
                              count = segment.R + 1
                           }
                           totalSegments += count
                        }
                        startNumber := segTemplate.StartNumber
                        if startNumber == 0 {
                           startNumber = 1
                        }
                        for i := 0; i < totalSegments; i++ {
                           mediaURL := strings.ReplaceAll(mediaTemplate, "$Number$", fmt.Sprintf("%d", startNumber+i))
                           resolvedURL, err := resolveURL(repBaseURL, mediaURL)
                           if err != nil {
                              fmt.Fprintf(os.Stderr, "Warning: failed to resolve media URL for %s at number %d: %v\n", rep.ID, startNumber+i, err)
                              continue
                           }
                           fmt.Fprintf(os.Stderr, "Debug: Generated segment URL for %s: %s\n", rep.ID, resolvedURL)
                           segmentURLs = append(segmentURLs, resolvedURL)
                        }
                     } else if strings.Contains(mediaTemplate, "$Number$") {
                        // Fallback for $Number$ without SegmentTimeline
                        fmt.Fprintf(os.Stderr, "Debug: Processing $Number$ segments (no SegmentTimeline) for %s\n", rep.ID)
                        startNumber := segTemplate.StartNumber
                        if startNumber == 0 {
                           startNumber = 1
                        }
                        endNumber := startNumber + 79 // Fallback for thumbnails
                        for i := startNumber; i <= endNumber; i++ {
                           mediaURL := strings.ReplaceAll(mediaTemplate, "$Number$", fmt.Sprintf("%d", i))
                           resolvedURL, err := resolveURL(repBaseURL, mediaURL)
                           if err != nil {
                              fmt.Fprintf(os.Stderr, "Warning: failed to resolve media URL for %s at number %d: %v\n", rep.ID, i, err)
                              continue
                           }
                           fmt.Fprintf(os.Stderr, "Debug: Generated segment URL for %s: %s\n", rep.ID, resolvedURL)
                           segmentURLs = append(segmentURLs, resolvedURL)
                        }
                     } else {
                        fmt.Fprintf(os.Stderr, "Warning: no $Time$ or $Number$ in media template for %s: %s\n", rep.ID, mediaTemplate)
                     }
                  }
               } else {
                  fmt.Fprintf(os.Stderr, "Warning: no SegmentTemplate for representation %s\n", rep.ID)
               }
            }

            output.Representations[rep.ID] = segmentURLs
         }
      }
   }

   return output, nil
}

func resolveURL(base *url.URL, ref string) (string, error) {
   relative, err := url.Parse(ref)
   if err != nil {
      return "", fmt.Errorf("failed to parse reference URL %s: %w", ref, err)
   }
   resolved := base.ResolveReference(relative)
   // Normalize the path to avoid duplication
   resolved.Path = path.Clean(resolved.Path)
   return resolved.String(), nil
}

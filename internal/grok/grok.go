package main

import (
   "encoding/json"
   "encoding/xml"
   "fmt"
   "os"
   "path"
   "strings"
)

type MPD struct {
   XMLName xml.Name `xml:"MPD"`
   Period  Period   `xml:"Period"`
   BaseURL string   `xml:"BaseURL"`
}

type Period struct {
   AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
   BaseURL        string          `xml:"BaseURL"`
}

type AdaptationSet struct {
   Representations []Representation `xml:"Representation"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
}

type Representation struct {
   ID              string           `xml:"id,attr"`
   SegmentBase     *SegmentBase     `xml:"SegmentBase"`
   SegmentList     *SegmentList     `xml:"SegmentList"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   BaseURL         string           `xml:"BaseURL"`
}

type SegmentBase struct {
   Initialization string `xml:"initialization,attr"`
}

type SegmentList struct {
   Initialization string       `xml:"initialization,attr"`
   SegmentURLs    []SegmentURL `xml:"SegmentURL"`
}

type SegmentURL struct {
   Media string `xml:"media,attr"`
}

type SegmentTemplate struct {
   Initialization  string           `xml:"initialization,attr"`
   Media           string           `xml:"media,attr"`
   SegmentTimeline *SegmentTimeline `xml:"SegmentTimeline"`
   Timescale       int              `xml:"timescale,attr"`
   Duration        int              `xml:"duration,attr"`
   StartNumber     int              `xml:"startNumber,attr"`
   EndNumber       int              `xml:"endNumber,attr"`
}

type SegmentTimeline struct {
   Segments []Segment `xml:"S"`
}

type Segment struct {
   Time     uint64 `xml:"t,attr"`
   Duration uint64 `xml:"d,attr"`
   Repeat   int    `xml:"r,attr"`
}

type Output struct {
   RepresentationURLs map[string][]string `json:"representation_urls"`
}

func main() {
   if len(os.Args) != 2 {
      fmt.Fprintln(os.Stderr, "Usage: go run mpd_parser.go <mpd_file>")
      os.Exit(1)
   }

   mpdFile := os.Args[1]
   data, err := os.ReadFile(mpdFile)
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error reading MPD file: %v\n", err)
      os.Exit(1)
   }

   var mpd MPD
   err = xml.Unmarshal(data, &mpd)
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error parsing MPD file: %v\n", err)
      os.Exit(1)
   }

   output := Output{RepresentationURLs: make(map[string][]string)}
   mpdBaseURL := mpd.BaseURL
   fmt.Fprintf(os.Stderr, "MPD BaseURL: %s\n", mpdBaseURL)

   for i, period := range []Period{mpd.Period} {
      fmt.Fprintf(os.Stderr, "Processing Period %d, BaseURL: %s\n", i, period.BaseURL)
      for j, adaptationSet := range period.AdaptationSets {
         fmt.Fprintf(os.Stderr, "  Processing AdaptationSet %d\n", j)
         for k, rep := range adaptationSet.Representations {
            fmt.Fprintf(os.Stderr, "    Processing Representation %d: %s\n", k, rep.ID)
            urls := []string{}
            // Combine BaseURLs: MPD > Period > Representation
            baseURL := mpdBaseURL
            if period.BaseURL != "" {
               baseURL = joinURLs(baseURL, period.BaseURL)
               fmt.Fprintf(os.Stderr, "      After Period BaseURL: %s\n", baseURL)
            }
            if rep.BaseURL != "" {
               baseURL = joinURLs(baseURL, rep.BaseURL)
               fmt.Fprintf(os.Stderr, "      After Representation BaseURL: %s\n", baseURL)
               // Add Representation BaseURL directly if present
               urls = append(urls, baseURL)
               fmt.Fprintf(os.Stderr, "      Added Representation BaseURL: %s\n", baseURL)
               if strings.Contains(baseURL, "..") {
                  fmt.Fprintf(os.Stderr, "      Warning: Representation BaseURL contains '..': %s\n", baseURL)
               }
            }
            fmt.Fprintf(os.Stderr, "      Combined BaseURL: %s\n", baseURL)
            if strings.Contains(baseURL, "..") {
               fmt.Fprintf(os.Stderr, "      Warning: Combined BaseURL contains '..': %s\n", baseURL)
            }

            // Check for SegmentTemplate at Representation or AdaptationSet level
            segTemplate := rep.SegmentTemplate
            if segTemplate == nil {
               segTemplate = adaptationSet.SegmentTemplate
            }

            if rep.SegmentBase != nil && rep.SegmentBase.Initialization != "" {
               url := strings.ReplaceAll(rep.SegmentBase.Initialization, "$RepresentationID$", rep.ID)
               cleanURL := joinURLs(baseURL, url)
               urls = append(urls, cleanURL)
               fmt.Fprintf(os.Stderr, "      Added SegmentBase URL: %s\n", cleanURL)
               if strings.Contains(cleanURL, "..") {
                  fmt.Fprintf(os.Stderr, "      Warning: SegmentBase URL contains '..': %s\n", cleanURL)
               }
            }

            if rep.SegmentList != nil {
               if rep.SegmentList.Initialization != "" {
                  url := strings.ReplaceAll(rep.SegmentList.Initialization, "$RepresentationID$", rep.ID)
                  cleanURL := joinURLs(baseURL, url)
                  urls = append(urls, cleanURL)
                  fmt.Fprintf(os.Stderr, "      Added SegmentList Initialization URL: %s\n", cleanURL)
                  if strings.Contains(cleanURL, "..") {
                     fmt.Fprintf(os.Stderr, "      Warning: SegmentList Initialization URL contains '..': %s\n", cleanURL)
                  }
               }
               for _, segURL := range rep.SegmentList.SegmentURLs {
                  cleanURL := joinURLs(baseURL, segURL.Media)
                  urls = append(urls, cleanURL)
                  fmt.Fprintf(os.Stderr, "      Added SegmentList URL: %s\n", cleanURL)
                  if strings.Contains(cleanURL, "..") {
                     fmt.Fprintf(os.Stderr, "      Warning: SegmentList URL contains '..': %s\n", cleanURL)
                  }
               }
            }

            if segTemplate != nil {
               fmt.Fprintf(os.Stderr, "      SegmentTemplate found: Media=%s, Initialization=%s, Timescale=%d\n", segTemplate.Media, segTemplate.Initialization, segTemplate.Timescale)
               // Replace $RepresentationID$ in initialization and media templates
               initURL := segTemplate.Initialization
               mediaTemplate := segTemplate.Media
               if initURL != "" {
                  initURL = strings.ReplaceAll(initURL, "$RepresentationID$", rep.ID)
                  cleanURL := joinURLs(baseURL, initURL)
                  urls = append(urls, cleanURL)
                  fmt.Fprintf(os.Stderr, "      Added Initialization URL: %s\n", cleanURL)
                  if strings.Contains(cleanURL, "..") {
                     fmt.Fprintf(os.Stderr, "      Warning: Initialization URL contains '..': %s\n", cleanURL)
                  }
               }

               if mediaTemplate != "" {
                  mediaTemplate = strings.ReplaceAll(mediaTemplate, "$RepresentationID$", rep.ID)
                  if segTemplate.SegmentTimeline != nil {
                     fmt.Fprintf(os.Stderr, "      SegmentTimeline found with %d segments\n", len(segTemplate.SegmentTimeline.Segments))
                     if strings.Contains(mediaTemplate, "$Time$") {
                        // Handle $Time$ placeholder
                        currentTime := uint64(0)
                        for k, seg := range segTemplate.SegmentTimeline.Segments {
                           if seg.Time != 0 {
                              currentTime = seg.Time
                           }
                           mediaURL := strings.ReplaceAll(mediaTemplate, "$Time$", fmt.Sprintf("%d", currentTime))
                           cleanURL := joinURLs(baseURL, mediaURL)
                           urls = append(urls, cleanURL)
                           fmt.Fprintf(os.Stderr, "        Added Time-based URL %d (t=%d, d=%d, r=%d): %s\n", k, currentTime, seg.Duration, seg.Repeat, cleanURL)
                           if strings.Contains(cleanURL, "..") {
                              fmt.Fprintf(os.Stderr, "        Warning: Time-based URL contains '..': %s\n", cleanURL)
                           }
                           for j := 0; j < seg.Repeat; j++ {
                              currentTime += seg.Duration
                              mediaURL = strings.ReplaceAll(mediaTemplate, "$Time$", fmt.Sprintf("%d", currentTime))
                              cleanURL := joinURLs(baseURL, mediaURL)
                              urls = append(urls, cleanURL)
                              fmt.Fprintf(os.Stderr, "        Added Repeated Time-based URL %d (t=%d): %s\n", j+1, currentTime, cleanURL)
                              if strings.Contains(cleanURL, "..") {
                                 fmt.Fprintf(os.Stderr, "        Warning: Repeated Time-based URL contains '..': %s\n", cleanURL)
                              }
                           }
                           currentTime += seg.Duration
                        }
                     } else if strings.Contains(mediaTemplate, "$Number$") {
                        // Handle $Number$ placeholder in SegmentTimeline
                        currentNumber := 1
                        for k, seg := range segTemplate.SegmentTimeline.Segments {
                           mediaURL := strings.ReplaceAll(mediaTemplate, "$Number$", fmt.Sprintf("%d", currentNumber))
                           cleanURL := joinURLs(baseURL, mediaURL)
                           urls = append(urls, cleanURL)
                           fmt.Fprintf(os.Stderr, "        Added Number-based URL %d (n=%d, r=%d): %s\n", k, currentNumber, seg.Repeat, cleanURL)
                           if strings.Contains(cleanURL, "..") {
                              fmt.Fprintf(os.Stderr, "        Warning: Number-based URL contains '..': %s\n", cleanURL)
                           }
                           for j := 0; j < seg.Repeat; j++ {
                              currentNumber++
                              mediaURL = strings.ReplaceAll(mediaTemplate, "$Number$", fmt.Sprintf("%d", currentNumber))
                              cleanURL = joinURLs(baseURL, mediaURL)
                              urls = append(urls, cleanURL)
                              fmt.Fprintf(os.Stderr, "        Added Repeated Number-based URL %d (n=%d): %s\n", j+1, currentNumber, cleanURL)
                              if strings.Contains(cleanURL, "..") {
                                 fmt.Fprintf(os.Stderr, "        Warning: Repeated Number-based URL contains '..': %s\n", cleanURL)
                              }
                           }
                           currentNumber++
                        }
                     }
                  } else if segTemplate.Duration > 0 && segTemplate.EndNumber > 0 {
                     // Handle SegmentTemplate without SegmentTimeline
                     start := segTemplate.StartNumber
                     if start == 0 {
                        start = 1 // Default to 1 if not specified
                     }
                     end := segTemplate.EndNumber
                     fmt.Fprintf(os.Stderr, "      Non-timeline SegmentTemplate: startNumber=%d, endNumber=%d\n", start, end)
                     for i := start; i <= end; i++ {
                        mediaURL := strings.ReplaceAll(mediaTemplate, "$Number$", fmt.Sprintf("%d", i))
                        cleanURL := joinURLs(baseURL, mediaURL)
                        urls = append(urls, cleanURL)
                        fmt.Fprintf(os.Stderr, "        Added URL %d: %s\n", i, cleanURL)
                        if strings.Contains(cleanURL, "..") {
                           fmt.Fprintf(os.Stderr, "        Warning: Non-timeline URL contains '..': %s\n", cleanURL)
                        }
                     }
                  } else {
                     fmt.Fprintf(os.Stderr, "      Warning: SegmentTemplate has no SegmentTimeline or valid duration/endNumber\n")
                  }
               }
            } else {
               fmt.Fprintf(os.Stderr, "      No SegmentTemplate found for representation %s\n", rep.ID)
            }

            output.RepresentationURLs[rep.ID] = urls
            fmt.Fprintf(os.Stderr, "      Total URLs for %s: %d\n", rep.ID, len(urls))
         }
      }
   }

   jsonOutput, err := json.MarshalIndent(output, "", "  ")
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error generating JSON output: %v\n", err)
      os.Exit(1)
   }

   fmt.Println(string(jsonOutput))
}

func joinURLs(base, relative string) string {
   // Split query parameters to preserve them
   var basePath, baseQuery, relPath, relQuery string
   if idx := strings.Index(base, "?"); idx != -1 {
      basePath = base[:idx]
      baseQuery = base[idx:]
   } else {
      basePath = base
   }
   if idx := strings.Index(relative, "?"); idx != -1 {
      relPath = relative[:idx]
      relQuery = relative[idx:]
   } else {
      relPath = relative
   }

   // Ensure forward slashes
   basePath = strings.ReplaceAll(basePath, "\\", "/")
   relPath = strings.ReplaceAll(relPath, "\\", "/")

   // Strip excessive parent directories from basePath
   for strings.HasPrefix(basePath, "../") {
      basePath = strings.TrimPrefix(basePath, "../")
   }
   if basePath == "" || basePath == strings.Repeat("../", 13) {
      basePath = "/"
   }

   // Combine paths
   var combinedPath string
   if basePath == "" || basePath == "/" {
      combinedPath = path.Clean("/" + relPath)
   } else {
      combinedPath = path.Clean(basePath + "/" + relPath)
   }

   // Combine query parameters (relative query takes precedence if present)
   query := baseQuery
   if relQuery != "" {
      query = relQuery
   }

   // Construct final URL
   cleanURL := combinedPath
   if query != "" {
      cleanURL += query
   }

   // Log for debugging
   fmt.Fprintf(os.Stderr, "      joinURLs: basePath=%s, baseQuery=%s, relPath=%s, relQuery=%s, combinedPath=%s, cleanURL=%s\n",
      basePath, baseQuery, relPath, relQuery, combinedPath, cleanURL)

   // Final check for '..'
   if strings.Contains(cleanURL, "..") {
      fmt.Fprintf(os.Stderr, "      Warning: Final URL contains '..': %s\n", cleanURL)
   }

   return cleanURL
}

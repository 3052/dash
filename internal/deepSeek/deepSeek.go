package main

import (
   "encoding/json"
   "encoding/xml"
   "fmt"
   "io/ioutil"
   "net/url"
   "os"
   "path"
   "strconv"
   "strings"
)

type MPD struct {
   XMLName xml.Name `xml:"MPD"`
   Periods []Period `xml:"Period"`
   BaseURL string   `xml:"BaseURL"`
}

type Period struct {
   AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
   BaseURL        string          `xml:"BaseURL"`
}

type AdaptationSet struct {
   Representations []Representation `xml:"Representation"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   BaseURL         string           `xml:"BaseURL"`
   ContentType     string           `xml:"contentType,attr"`
}

type Representation struct {
   ID              string           `xml:"id,attr"`
   BaseURL         string           `xml:"BaseURL"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   SegmentList     *SegmentList     `xml:"SegmentList"`
   SegmentBase     *SegmentBase     `xml:"SegmentBase"`
}

type SegmentTemplate struct {
   Initialization  string           `xml:"initialization,attr"`
   Media           string           `xml:"media,attr"`
   StartNumber     int              `xml:"startNumber,attr"`
   Duration        int              `xml:"duration,attr"`
   Timescale       int              `xml:"timescale,attr"`
   EndNumber       int              `xml:"endNumber,attr"`
   SegmentTimeline *SegmentTimeline `xml:"SegmentTimeline"`
}

type SegmentTimeline struct {
   Segments []TimelineSegment `xml:"S"`
}

type TimelineSegment struct {
   T int64 `xml:"t,attr"`
   D int64 `xml:"d,attr"`
   R int   `xml:"r,attr"`
}

type SegmentList struct {
   Initialization *URL  `xml:"Initialization"`
   SegmentURLs    []URL `xml:"SegmentURL"`
}

type SegmentBase struct {
   Initialization *URL `xml:"Initialization"`
}

type URL struct {
   SourceURL string `xml:"sourceURL,attr"`
}

type Result struct {
   Representations map[string][]string `json:"representations"`
}

func main() {
   if len(os.Args) < 2 {
      fmt.Println("Usage: go run dash_parser.go <path_to_mpd_file>")
      os.Exit(1)
   }

   mpdPath := os.Args[1]
   mpdURL := "http://test.test/test.mpd"

   data, err := ioutil.ReadFile(mpdPath)
   if err != nil {
      fmt.Printf("Error reading file: %v\n", err)
      os.Exit(1)
   }

   var mpd MPD
   err = xml.Unmarshal(data, &mpd)
   if err != nil {
      fmt.Printf("Error parsing MPD: %v\n", err)
      os.Exit(1)
   }

   result := Result{
      Representations: make(map[string][]string),
   }

   for _, period := range mpd.Periods {
      baseURL := mpd.BaseURL
      if period.BaseURL != "" {
         baseURL = period.BaseURL
      }

      for _, adaptSet := range period.AdaptationSets {
         adaptBaseURL := baseURL
         if adaptSet.BaseURL != "" {
            adaptBaseURL = adaptSet.BaseURL
         }

         for _, rep := range adaptSet.Representations {
            var segments []string
            repBaseURL := adaptBaseURL
            if rep.BaseURL != "" {
               repBaseURL = rep.BaseURL
            }

            // Handle simple BaseURL case
            if rep.BaseURL != "" && rep.SegmentTemplate == nil && rep.SegmentList == nil && rep.SegmentBase == nil {
               segments = append(segments, resolveURL(mpdURL, repBaseURL, rep.BaseURL))
            }

            // Handle SegmentTemplate
            if rep.SegmentTemplate != nil || adaptSet.SegmentTemplate != nil {
               segTemplate := rep.SegmentTemplate
               if segTemplate == nil {
                  segTemplate = adaptSet.SegmentTemplate
               }

               startNumber := 1
               if segTemplate.StartNumber > 0 {
                  startNumber = segTemplate.StartNumber
               }

               // Initialization segment
               if segTemplate.Initialization != "" {
                  initURL := strings.Replace(segTemplate.Initialization, "$RepresentationID$", rep.ID, -1)
                  initURL = strings.Replace(initURL, "$Number$", "0", -1)
                  segments = append(segments, resolveURL(mpdURL, repBaseURL, initURL))
               }

               // Media segments
               if segTemplate.Media != "" {
                  if segTemplate.SegmentTimeline != nil {
                     // Timeline-based segments
                     segmentNumber := startNumber
                     for _, seg := range segTemplate.SegmentTimeline.Segments {
                        repeat := 1
                        if seg.R > 0 {
                           repeat += seg.R
                        }

                        for i := 0; i < repeat; i++ {
                           time := seg.T
                           if i > 0 {
                              time += seg.D * int64(i)
                           }

                           mediaURL := strings.Replace(segTemplate.Media, "$RepresentationID$", rep.ID, -1)
                           mediaURL = strings.Replace(mediaURL, "$Number$", strconv.Itoa(segmentNumber), -1)
                           mediaURL = strings.Replace(mediaURL, "$Time$", strconv.FormatInt(time, 10), -1)
                           segments = append(segments, resolveURL(mpdURL, repBaseURL, mediaURL))
                           segmentNumber++
                        }
                     }
                  } else {
                     // Number-based segments
                     endNumber := startNumber + 5 // Default to 5 segments
                     if segTemplate.EndNumber > 0 {
                        endNumber = segTemplate.EndNumber
                     }

                     for i := startNumber; i <= endNumber; i++ {
                        mediaURL := strings.Replace(segTemplate.Media, "$RepresentationID$", rep.ID, -1)
                        mediaURL = strings.Replace(mediaURL, "$Number$", strconv.Itoa(i), -1)
                        mediaURL = strings.Replace(mediaURL, "$Time$", strconv.Itoa(i*segTemplate.Duration), -1)
                        segments = append(segments, resolveURL(mpdURL, repBaseURL, mediaURL))
                     }
                  }
               }
            }

            // Handle SegmentList
            if rep.SegmentList != nil {
               if rep.SegmentList.Initialization != nil {
                  segments = append(segments, resolveURL(mpdURL, repBaseURL, rep.SegmentList.Initialization.SourceURL))
               }
               for _, seg := range rep.SegmentList.SegmentURLs {
                  segments = append(segments, resolveURL(mpdURL, repBaseURL, seg.SourceURL))
               }
            }

            // Handle SegmentBase
            if rep.SegmentBase != nil && rep.SegmentBase.Initialization != nil {
               segments = append(segments, resolveURL(mpdURL, repBaseURL, rep.SegmentBase.Initialization.SourceURL))
            }

            if len(segments) > 0 {
               result.Representations[rep.ID] = segments
            }
         }
      }
   }

   jsonData, err := json.MarshalIndent(result, "", "  ")
   if err != nil {
      fmt.Printf("Error generating JSON: %v\n", err)
      os.Exit(1)
   }

   fmt.Println(string(jsonData))
}

// resolveURL properly resolves URLs with scheme and host
func resolveURL(mpdURL, baseURL, relativeURL string) string {

   out, _ := url.Parse(mpdURL)
   out, _ = out.Parse(baseURL)
   if relativeURL != baseURL {
      out, _ = out.Parse(relativeURL)
   }
   return out.String()

   // Parse the base MPD URL
   base, err := url.Parse(mpdURL)
   if err != nil {
      return relativeURL
   }

   // Handle absolute URLs
   if strings.HasPrefix(relativeURL, "http://") || strings.HasPrefix(relativeURL, "https://") {
      return relativeURL
   }

   // Create a base URL for resolution
   resolutionBase := base

   // If we have a baseURL, parse it and resolve against the MPD URL
   if baseURL != "" {
      relBase, err := url.Parse(baseURL)
      if err == nil {
         // Resolve the baseURL against the MPD URL
         resolutionBase = base.ResolveReference(relBase)
      }
   }

   // Parse the relative URL
   relURL, err := url.Parse(relativeURL)
   if err != nil {
      return relativeURL
   }

   // Resolve against the proper base URL
   resolved := resolutionBase.ResolveReference(relURL)

   // Clean the path to remove any . or .. sequences
   resolved.Path = path.Clean(resolved.Path)

   // Return the string representation
   return resolved.String()
}

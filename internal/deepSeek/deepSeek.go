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

// MPD struct represents the root MPD element
type MPD struct {
   XMLName                   xml.Name `xml:"MPD"`
   Periods                   []Period `xml:"Period"`
   BaseURL                   string   `xml:"BaseURL"`
   Type                      string   `xml:"type,attr"`
   MediaPresentationDuration string   `xml:"mediaPresentationDuration,attr"`
}

// Period struct represents a Period element in MPD
type Period struct {
   AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
}

// AdaptationSet struct represents an AdaptationSet element
type AdaptationSet struct {
   Representations []Representation `xml:"Representation"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   BaseURL         string           `xml:"BaseURL"`
   MimeType        string           `xml:"mimeType,attr"`
}

// Representation struct represents a Representation element
type Representation struct {
   ID              string           `xml:"id,attr"`
   Bandwidth       int              `xml:"bandwidth,attr,omitempty"`
   Width           int              `xml:"width,attr,omitempty"`
   Height          int              `xml:"height,attr,omitempty"`
   BaseURL         string           `xml:"BaseURL"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   SegmentList     *SegmentList     `xml:"SegmentList"`
   SegmentBase     *SegmentBase     `xml:"SegmentBase"`
}

// [Other struct definitions remain the same...]

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
      for _, adaptSet := range period.AdaptationSets {
         for _, rep := range adaptSet.Representations {
            var segments []string

            baseURL := mpd.BaseURL
            if adaptSet.BaseURL != "" {
               baseURL = adaptSet.BaseURL
            }
            if rep.BaseURL != "" {
               baseURL = rep.BaseURL
            }

            // Handle simple BaseURL case (like for subtitles)
            if rep.BaseURL != "" && rep.SegmentTemplate == nil && rep.SegmentList == nil && rep.SegmentBase == nil {
               segments = append(segments, resolveURL(mpdURL, baseURL, rep.BaseURL))
            }

            // Handle SegmentTemplate
            if rep.SegmentTemplate != nil {
               startNumber := 1
               if rep.SegmentTemplate.StartNumber > 0 {
                  startNumber = rep.SegmentTemplate.StartNumber
               }

               if rep.SegmentTemplate.Initialization != "" {
                  initURL := strings.Replace(rep.SegmentTemplate.Initialization, "$Number$", "0", -1)
                  initURL = strings.Replace(initURL, "$Time$", "0", -1)
                  segments = append(segments, resolveURL(mpdURL, baseURL, initURL))
               }

               if rep.SegmentTemplate.Media != "" && rep.SegmentTemplate.SegmentTimeline != nil {
                  segments = append(segments, generateSegmentURLsFromTimeline(
                     mpdURL,
                     baseURL,
                     rep.SegmentTemplate.Media,
                     rep.SegmentTemplate.SegmentTimeline,
                     startNumber,
                  )...)
               } else if rep.SegmentTemplate.Media != "" {
                  for i := startNumber; i < startNumber+5; i++ {
                     mediaURL := strings.Replace(rep.SegmentTemplate.Media, "$Number$", fmt.Sprintf("%d", i), -1)
                     mediaURL = strings.Replace(mediaURL, "$Time$", fmt.Sprintf("%d", i*rep.SegmentTemplate.Duration), -1)
                     segments = append(segments, resolveURL(mpdURL, baseURL, mediaURL))
                  }
               }
            }

            // [Other segment handling code remains the same...]

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

// SegmentTemplate struct represents a SegmentTemplate element
type SegmentTemplate struct {
   Initialization  string           `xml:"initialization,attr"`
   Media           string           `xml:"media,attr"`
   StartNumber     int              `xml:"startNumber,attr,omitempty"`
   Duration        int              `xml:"duration,attr,omitempty"`
   Timescale       int              `xml:"timescale,attr,omitempty"`
   SegmentTimeline *SegmentTimeline `xml:"SegmentTimeline"`
}

// SegmentTimeline represents the SegmentTimeline element
type SegmentTimeline struct {
   Segments []TimelineSegment `xml:"S"`
}

// TimelineSegment represents an S element in SegmentTimeline
type TimelineSegment struct {
   T int64 `xml:"t,attr"` // start time
   D int64 `xml:"d,attr"` // duration
   R int   `xml:"r,attr"` // repeat count
}

// SegmentList struct represents a SegmentList element
type SegmentList struct {
   Initialization *URL  `xml:"Initialization"`
   SegmentURLs    []URL `xml:"SegmentURL"`
}

// SegmentBase struct represents a SegmentBase element
type SegmentBase struct {
   Initialization *URL `xml:"Initialization"`
}

// URL struct represents URL elements in SegmentList/SegmentBase
type URL struct {
   SourceURL string `xml:"sourceURL,attr"`
}

// Result struct for JSON output
type Result struct {
   Representations map[string][]string `json:"representations"`
}

// generateSegmentURLsFromTimeline generates segment URLs based on SegmentTimeline
func generateSegmentURLsFromTimeline(mpdURL, baseURL, mediaTemplate string, timeline *SegmentTimeline, startNumber int) []string {
   var segments []string
   segmentNumber := startNumber

   for _, seg := range timeline.Segments {
      repeat := 1
      if seg.R > 0 {
         repeat += seg.R
      }

      for i := 0; i < repeat; i++ {
         time := seg.T
         if i > 0 {
            time += seg.D * int64(i)
         }

         mediaURL := strings.Replace(mediaTemplate, "$Number$", strconv.Itoa(segmentNumber), -1)
         mediaURL = strings.Replace(mediaURL, "$Time$", strconv.FormatInt(time, 10), -1)
         segments = append(segments, resolveURL(mpdURL, baseURL, mediaURL))
         segmentNumber++
      }
   }

   return segments
}

// resolveURL handles the URL resolution correctly without any duplication
func resolveURL(mpdURL, baseURL, relativeURL string) string {
   // If URL is already absolute, return as-is
   if strings.HasPrefix(relativeURL, "http://") || strings.HasPrefix(relativeURL, "https://") {
      return relativeURL
   }

   // Parse MPD URL to get scheme and host
   mpd, err := url.Parse(mpdURL)
   if err != nil {
      return relativeURL
   }

   // For relative URLs that already start with the path component,
   // just prepend the scheme and host
   if strings.HasPrefix(relativeURL, "/") {
      return fmt.Sprintf("%s://%s%s", mpd.Scheme, mpd.Host, relativeURL)
   }

   // For relative paths, use them as-is with the MPD's directory
   mpdDir := path.Dir(mpd.Path)
   if mpdDir == "." {
      mpdDir = ""
   }

   // Join with MPD directory (only if relativeURL doesn't already contain full path)
   finalPath := path.Join(mpdDir, relativeURL)

   // Rebuild URL
   return fmt.Sprintf("%s://%s/%s", mpd.Scheme, mpd.Host, strings.TrimPrefix(finalPath, "/"))
}

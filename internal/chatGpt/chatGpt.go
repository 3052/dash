package main

import (
   "encoding/json"
   "encoding/xml"
   "fmt"
   "math"
   "net/url"
   "os"
   "path"
   "strconv"
   "strings"
   "time"
)

type MPD struct {
   XMLName                   xml.Name `xml:"MPD"`
   BaseURL                   string   `xml:"BaseURL"`
   MediaPresentationDuration string   `xml:"mediaPresentationDuration,attr"`
   Periods                   []Period `xml:"Period"`
}

type Period struct {
   ID       string          `xml:"id,attr"`
   Duration string          `xml:"duration,attr"`
   BaseURL  string          `xml:"BaseURL"`
   Sets     []AdaptationSet `xml:"AdaptationSet"`
}

type AdaptationSet struct {
   ID              string           `xml:"id,attr"`
   BaseURL         string           `xml:"BaseURL"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   Representations []Representation `xml:"Representation"`
}

type Representation struct {
   ID              string           `xml:"id,attr"`
   BaseURL         string           `xml:"BaseURL"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
}

type SegmentTemplate struct {
   Timescale       int              `xml:"timescale,attr"`
   Duration        int64            `xml:"duration,attr"`
   StartNumber     int              `xml:"startNumber,attr"`
   EndNumber       int              `xml:"endNumber,attr"`
   Initialization  string           `xml:"initialization,attr"`
   Media           string           `xml:"media,attr"`
   SegmentTimeline *SegmentTimeline `xml:"SegmentTimeline"`
}

type SegmentTimeline struct {
   Segments []Segment `xml:"S"`
}

type Segment struct {
   T int64 `xml:"t,attr"`
   D int64 `xml:"d,attr"`
   R int   `xml:"r,attr"`
}

func resolveURL(base, rel string) string {
   baseURL, _ := url.Parse(base)
   relURL, _ := url.Parse(rel)
   return baseURL.ResolveReference(relURL).String()
}

func joinBaseURLs(parts ...string) string {
   joined := path.Join(parts...)
   if !strings.Contains(path.Base(joined), ".") {
      joined += "/"
   }
   return joined
}

func parseDuration(d string) time.Duration {
   if d == "" {
      return 0
   }
   // Basic ISO8601 PT{seconds}S support
   if strings.HasPrefix(d, "PT") && strings.HasSuffix(d, "S") {
      secStr := strings.TrimSuffix(strings.TrimPrefix(d, "PT"), "S")
      if f, err := strconv.ParseFloat(secStr, 64); err == nil {
         return time.Duration(f * float64(time.Second))
      }
   }
   return 0
}

func contains(list []string, item string) bool {
   for _, v := range list {
      if v == item {
         return true
      }
   }
   return false
}

func main() {
   if len(os.Args) < 2 {
      fmt.Println("Usage: go run extract.go <path-to-mpd>")
      return
   }

   filePath := os.Args[1]
   data, err := os.ReadFile(filePath)
   if err != nil {
      panic(err)
   }

   var mpd MPD
   if err := xml.Unmarshal(data, &mpd); err != nil {
      panic(err)
   }

   baseMPD := "http://test.test/test.mpd"
   result := make(map[string][]string)

   for _, period := range mpd.Periods {
      periodDuration := parseDuration(period.Duration)
      if periodDuration == 0 {
         periodDuration = parseDuration(mpd.MediaPresentationDuration)
      }

      for _, set := range period.Sets {
         for _, rep := range set.Representations {
            repID := rep.ID
            segmentURLs := result[repID] // accumulate across periods

            fullBase := joinBaseURLs(mpd.BaseURL, period.BaseURL, set.BaseURL, rep.BaseURL)
            fullBaseURL := resolveURL(baseMPD, fullBase)

            st := rep.SegmentTemplate
            if st == nil {
               st = set.SegmentTemplate
            }

            if st == nil {
               if rep.BaseURL != "" {
                  segmentURLs = append(segmentURLs, fullBaseURL)
               }
               if len(segmentURLs) > 0 {
                  result[repID] = segmentURLs
               }
               continue
            }

            // Deduplicate init segment
            if st.Initialization != "" {
               init := strings.ReplaceAll(st.Initialization, "$RepresentationID$", repID)
               initURL := resolveURL(fullBaseURL, init)
               if !contains(segmentURLs, initURL) {
                  segmentURLs = append(segmentURLs, initURL)
               }
            }

            // SegmentTimeline
            if st.SegmentTimeline != nil {
               number := st.StartNumber
               if number == 0 {
                  number = 1
               }
               for _, seg := range st.SegmentTimeline.Segments {
                  repeat := seg.R
                  if repeat < 0 {
                     repeat = 0
                  }
                  count := 1 + repeat
                  for i := 0; i < count; i++ {
                     media := strings.ReplaceAll(st.Media, "$RepresentationID$", repID)
                     media = strings.ReplaceAll(media, "$Number$", strconv.Itoa(number))
                     segmentURLs = append(segmentURLs, resolveURL(fullBaseURL, media))
                     number++
                  }
               }
            } else if st.EndNumber > 0 {
               number := st.StartNumber
               if number == 0 {
                  number = 1
               }
               for i := number; i <= st.EndNumber; i++ {
                  media := strings.ReplaceAll(st.Media, "$RepresentationID$", repID)
                  media = strings.ReplaceAll(media, "$Number$", strconv.Itoa(i))
                  segmentURLs = append(segmentURLs, resolveURL(fullBaseURL, media))
               }
            } else if st.Duration > 0 && periodDuration > 0 {
               timescale := st.Timescale
               if timescale == 0 {
                  timescale = 1
               }
               durationSec := float64(st.Duration) / float64(timescale)
               numSegments := int(math.Ceil(periodDuration.Seconds() / durationSec))
               start := st.StartNumber
               if start == 0 {
                  start = 1
               }
               for i := 0; i < numSegments; i++ {
                  num := start + i
                  media := strings.ReplaceAll(st.Media, "$RepresentationID$", repID)
                  media = strings.ReplaceAll(media, "$Number$", strconv.Itoa(num))
                  segmentURLs = append(segmentURLs, resolveURL(fullBaseURL, media))
               }
            }

            if len(segmentURLs) > 0 {
               result[repID] = segmentURLs
            }
         }
      }
   }

   out, _ := json.MarshalIndent(result, "", "  ")
   fmt.Println(string(out))
}

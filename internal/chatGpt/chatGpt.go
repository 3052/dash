package main

import (
   "encoding/json"
   "encoding/xml"
   "fmt"
   "math"
   "net/url"
   "os"
   "regexp"
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
   SegmentBase     *SegmentBase     `xml:"SegmentBase"`
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

type SegmentBase struct {
   IndexRange     string          `xml:"indexRange,attr"`
   Initialization *Initialization `xml:"Initialization"`
}

type Initialization struct {
   Range string `xml:"range,attr"`
}

var numberPattern = regexp.MustCompile(`\$Number(?:%0(\d+)d)?\$`)

func resolveURL(base, rel string) string {
   baseURL, _ := url.Parse(base)
   relURL, _ := url.Parse(rel)
   return baseURL.ResolveReference(relURL).String()
}

func parseDuration(d string) time.Duration {
   if d == "" {
      return 0
   }
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

func replaceNumberPattern(template string, number int) string {
   return numberPattern.ReplaceAllStringFunc(template, func(match string) string {
      sub := numberPattern.FindStringSubmatch(match)
      if len(sub) > 1 && sub[1] != "" {
         width, _ := strconv.Atoi(sub[1])
         return fmt.Sprintf("%0*d", width, number)
      }
      return strconv.Itoa(number)
   })
}

func replaceSegmentTemplate(template, repID string, number int, timeVal int64) string {
   t := strings.ReplaceAll(template, "$RepresentationID$", repID)
   t = replaceNumberPattern(t, number)
   t = strings.ReplaceAll(t, "$Time$", strconv.FormatInt(timeVal, 10))
   return t
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
            segmentURLs := result[repID]

            // Resolve base
            base := baseMPD
            if mpd.BaseURL != "" {
               base = resolveURL(base, mpd.BaseURL)
            }
            if period.BaseURL != "" {
               base = resolveURL(base, period.BaseURL)
            }
            if set.BaseURL != "" {
               base = resolveURL(base, set.BaseURL)
            }
            if rep.BaseURL != "" {
               base = resolveURL(base, rep.BaseURL)
            }

            // SegmentBase
            if rep.SegmentBase != nil && rep.BaseURL != "" {
               mediaURL := base
               if !contains(segmentURLs, mediaURL) {
                  segmentURLs = append(segmentURLs, mediaURL)
               }
               if rep.SegmentBase.Initialization != nil {
                  if !contains(segmentURLs, mediaURL) {
                     segmentURLs = append(segmentURLs, mediaURL)
                  }
               }
               result[repID] = segmentURLs
               continue
            }

            // SegmentTemplate
            st := rep.SegmentTemplate
            if st == nil {
               st = set.SegmentTemplate
            }
            if st == nil {
               if rep.BaseURL != "" {
                  segmentURLs = append(segmentURLs, base)
               }
               if len(segmentURLs) > 0 {
                  result[repID] = segmentURLs
               }
               continue
            }

            // Initialization
            if st.Initialization != "" {
               init := replaceSegmentTemplate(st.Initialization, repID, 0, 0)
               initURL := resolveURL(base, init)
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
               t := int64(0)
               for _, seg := range st.SegmentTimeline.Segments {
                  repeat := seg.R
                  if repeat < 0 {
                     repeat = 0
                  }
                  startTime := seg.T
                  if startTime != 0 {
                     t = startTime
                  }
                  for i := 0; i <= repeat; i++ {
                     media := replaceSegmentTemplate(st.Media, repID, number, t)
                     segmentURLs = append(segmentURLs, resolveURL(base, media))
                     t += seg.D
                     number++
                  }
               }
            } else if st.EndNumber > 0 {
               number := st.StartNumber
               if number == 0 {
                  number = 1
               }
               for i := number; i <= st.EndNumber; i++ {
                  media := replaceSegmentTemplate(st.Media, repID, i, 0)
                  segmentURLs = append(segmentURLs, resolveURL(base, media))
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
               t := int64(0)
               for i := 0; i < numSegments; i++ {
                  num := start + i
                  media := replaceSegmentTemplate(st.Media, repID, num, t)
                  segmentURLs = append(segmentURLs, resolveURL(base, media))
                  t += st.Duration
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

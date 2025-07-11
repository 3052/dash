package main

import (
   "encoding/xml"
   "fmt"
   "net/url"
   "os"
   "path"
   "strconv"
   "strings"
)

// ---------- minimal MPD structs ----------

type MPD struct {
   XMLName                   xml.Name  `xml:"MPD"`
   BaseURL                   []BaseURL `xml:"BaseURL"`
   Periods                   []Period  `xml:"Period"`
   MediaPresentationDuration string    `xml:"mediaPresentationDuration,attr"`
}

type BaseURL struct {
   URL string `xml:",chardata"`
}

type Period struct {
   AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
}

type AdaptationSet struct {
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   Representations []Representation `xml:"Representation"`
}

type Representation struct {
   ID              string           `xml:"id,attr"`
   Bandwidth       string           `xml:"bandwidth,attr"`
   BaseURL         string           `xml:"BaseURL"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   SegmentList     *SegmentList     `xml:"SegmentList"`
}

type SegmentTemplate struct {
   Media           string           `xml:"media,attr"`
   Initialization  string           `xml:"initialization,attr"`
   StartNumber     int              `xml:"startNumber,attr"`
   Timescale       int              `xml:"timescale,attr"`
   Duration        int              `xml:"duration,attr"` // @duration
   SegmentTimeline *SegmentTimeline `xml:"SegmentTimeline"`
}

type SegmentTimeline struct {
   S []S `xml:"S"`
}

type S struct {
   T int `xml:"t,attr"` // media time
   D int `xml:"d,attr"` // duration
   R int `xml:"r,attr"` // repeat count
}

type SegmentList struct {
   SegmentURLs []struct {
      Media string `xml:"media,attr"`
   } `xml:"SegmentURL"`
}

// ---------- helpers ----------

func resolveBase(mpdPath string, mpd *MPD, rep *Representation) string {
   // 1. Representation-level BaseURL
   if rep.BaseURL != "" {
      if strings.HasPrefix(rep.BaseURL, "http") {
         return ensureSlash(rep.BaseURL)
      }
      return ensureSlash(rep.BaseURL)
   }
   // 2. MPD-level BaseURL
   if len(mpd.BaseURL) > 0 && mpd.BaseURL[0].URL != "" {
      return ensureSlash(strings.TrimSpace(mpd.BaseURL[0].URL))
   }
   // 3. derive from MPD file location
   u, _ := url.Parse(mpdPath)
   return u.ResolveReference(&url.URL{Path: path.Dir(u.Path) + "/"}).String()
}

func ensureSlash(s string) string {
   if !strings.HasSuffix(s, "/") {
      return s + "/"
   }
   return s
}

// expand substitutes $Number$, $Time$, $Bandwidth$, $RepresentationID$
func expand(tpl string, vars map[string]string) string {
   for k, v := range vars {
      tpl = strings.ReplaceAll(tpl, "$"+k+"$", v)
   }
   return tpl
}

// segmentTimes returns (time,nextTime) for every segment
func segmentTimes(st *SegmentTemplate) []int {
   var times []int
   if st.SegmentTimeline != nil && len(st.SegmentTimeline.S) > 0 {
      t := 0
      for _, s := range st.SegmentTimeline.S {
         if s.T != 0 { // explicit @t overrides inherited time
            t = s.T
         }
         repeats := s.R
         if repeats == 0 {
            repeats = 1
         }
         for i := 0; i < repeats; i++ {
            times = append(times, t)
            t += s.D
         }
      }
   } else {
      // uniform template
      start := st.StartNumber
      if start == 0 {
         start = 1
      }
      ts := st.Timescale
      if ts == 0 {
         ts = 1
      }
      d := st.Duration
      for i := 0; ; i++ {
         t := i * d
         // crude stop when we are past the presentation duration (optional)
         total := parseDuration(mpd.MediaPresentationDuration) * float64(ts)
         if float64(t) > total*1.2 { // generous exit
            break
         }
         times = append(times, t)
      }
   }
   return times
}

// very small ISO-8601 parser (PT1H2M30.5S)
func parseDuration(d string) float64 {
   if !strings.HasPrefix(d, "PT") {
      return 0
   }
   d = strings.TrimPrefix(d, "PT")
   h, m, s := 0.0, 0.0, 0.0
   if idx := strings.Index(d, "H"); idx != -1 {
      h, _ = strconv.ParseFloat(d[:idx], 64)
      d = d[idx+1:]
   }
   if idx := strings.Index(d, "M"); idx != -1 {
      m, _ = strconv.ParseFloat(d[:idx], 64)
      d = d[idx+1:]
   }
   if idx := strings.Index(d, "S"); idx != -1 {
      s, _ = strconv.ParseFloat(d[:idx], 64)
   }
   return h*3600 + m*60 + s
}

// ---------- main ----------

func main() {
   if len(os.Args) != 2 {
      fmt.Fprintf(os.Stderr, "usage: %s manifest.mpd\n", os.Args[0])
      os.Exit(1)
   }
   mpdFile := os.Args[1]

   data, err := os.ReadFile(mpdFile)
   if err != nil {
      fmt.Fprintf(os.Stderr, "read MPD: %v\n", err)
      os.Exit(1)
   }

   var mpd MPD
   if err := xml.Unmarshal(data, &mpd); err != nil {
      fmt.Fprintf(os.Stderr, "parse MPD: %v\n", err)
      os.Exit(1)
   }

   for _, period := range mpd.Periods {
      for _, as := range period.AdaptationSets {
         for _, rep := range as.Representations {
            base := resolveBase(mpdFile, &mpd, &rep)

            // effective template
            st := as.SegmentTemplate
            if rep.SegmentTemplate != nil {
               st = rep.SegmentTemplate
            }

            if st != nil {
               times := segmentTimes(st)
               start := st.StartNumber
               if start == 0 {
                  start = 1
               }
               for idx, t := range times {
                  vars := map[string]string{
                     "RepresentationID": rep.ID,
                     "Bandwidth":        rep.Bandwidth,
                     "Number":           strconv.Itoa(start + idx),
                     "Time":             strconv.Itoa(t),
                  }
                  seg := expand(st.Media, vars)
                  fmt.Println(base + seg)
               }
               continue
            }

            // SegmentList
            if rep.SegmentList != nil {
               for _, su := range rep.SegmentList.SegmentURLs {
                  fmt.Println(base + su.Media)
               }
               continue
            }
         }
      }
   }
}

package main

import (
   "encoding/json"
   "encoding/xml"
   "fmt"
   "net/url"
   "os"
   "path"
   "strconv"
   "strings"
)

// crude ISO-8601 → seconds
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

// ---------- minimal DASH structs ----------
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
   BaseURL        []BaseURL       `xml:"BaseURL"`
   AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
}
type AdaptationSet struct {
   BaseURL         []BaseURL        `xml:"BaseURL"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   Representations []Representation `xml:"Representation"`
}
type Representation struct {
   ID              string           `xml:"id,attr"`
   Bandwidth       string           `xml:"bandwidth,attr"`
   BaseURL         []BaseURL        `xml:"BaseURL"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   SegmentList     *SegmentList     `xml:"SegmentList"`
}
type SegmentTemplate struct {
   Media           string           `xml:"media,attr"`
   Initialization  string           `xml:"initialization,attr"`
   StartNumber     int              `xml:"startNumber,attr"`
   Timescale       int              `xml:"timescale,attr"`
   Duration        int              `xml:"duration,attr"`
   SegmentTimeline *SegmentTimeline `xml:"SegmentTimeline"`
   EndNumber       int              `xml:"endNumber,attr"`
}
type SegmentTimeline struct {
   S []S `xml:"S"`
}
type S struct {
   T int `xml:"t,attr"`
   D int `xml:"d,attr"`
   R int `xml:"r,attr"`
}
type SegmentList struct {
   SegmentURLs []struct {
      Media string `xml:"media,attr"`
   } `xml:"SegmentURL"`
}

// ---------- helpers ----------
func pickBaseURL(list []BaseURL) string {
   if len(list) > 0 {
      return strings.TrimSpace(list[0].URL)
   }
   return ""
}
func buildBase(mpdPath string, mpd *MPD, period *Period, as *AdaptationSet, rep *Representation) string {
   b := pickBaseURL(rep.BaseURL)
   if b == "" {
      b = pickBaseURL(as.BaseURL)
   }
   if b == "" {
      b = pickBaseURL(period.BaseURL)
   }
   if b == "" {
      b = pickBaseURL(mpd.BaseURL)
   }
   if b == "" {
      u, _ := url.Parse(mpdPath)
      b = u.ResolveReference(&url.URL{Path: path.Dir(u.Path) + "/"}).String()
   }
   if !strings.HasSuffix(b, "/") {
      b += "/"
   }
   if !strings.HasPrefix(b, "http") {
      u, _ := url.Parse(mpdPath)
      b = u.ResolveReference(&url.URL{Path: b}).String()
   }
   return b
}
func expand(tpl string, vars map[string]string) string {
   for k, v := range vars {
      tpl = strings.ReplaceAll(tpl, "$"+k+"$", v)
   }
   return tpl
}

func segmentTimes(st *SegmentTemplate, totalDur float64) []int {
   // 1. SegmentTimeline present
   if st.SegmentTimeline != nil && len(st.SegmentTimeline.S) > 0 {
      var times []int
      t := 0
      for _, s := range st.SegmentTimeline.S {
         if s.T != 0 {
            t = s.T
         }
         for i := 0; i < s.R+1; i++ {
            times = append(times, t)
            t += s.D
         }
      }
      return times
   }

   // 2. Uniform template
   ts := st.Timescale
   if ts == 0 {
      ts = 1
   }
   d := st.Duration
   if d == 0 {
      return nil
   }
   start := st.StartNumber
   if start == 0 {
      start = 1
   }
   // If endNumber is given, use it; otherwise compute from presentation duration
   segments := int(totalDur * float64(ts) / float64(d))
   if st.EndNumber > 0 {
      segments = st.EndNumber - start + 1
   }
   var times []int
   for i := 0; i < segments; i++ {
      times = append(times, i*d)
   }
   return times
}

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

   totalDur := parseDuration(mpd.MediaPresentationDuration)

   type Segment struct {
      URL    string `json:"url"`
      Time   int    `json:"time"`
      Number int    `json:"number"`
   }
   out := make(map[string][]Segment)

   for _, period := range mpd.Periods {
      for _, as := range period.AdaptationSets {
         for _, rep := range as.Representations {
            base := buildBase(mpdFile, &mpd, &period, &as, &rep)

            st := as.SegmentTemplate
            if rep.SegmentTemplate != nil {
               st = rep.SegmentTemplate
            }

            var segs []Segment

            if st != nil {
               // timeline mode
               times := segmentTimes(st, totalDur)
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
                  segs = append(segs, Segment{
                     URL:    base + seg,
                     Time:   t,
                     Number: start + idx,
                  })
               }
            } else if len(rep.BaseURL) > 0 {
               // single-segment (e.g. thumbnail)
               segs = append(segs, Segment{
                  URL:    base + rep.BaseURL[0].URL,
                  Time:   0,
                  Number: 1,
               })
            } else if rep.SegmentList != nil {
               for idx, su := range rep.SegmentList.SegmentURLs {
                  segs = append(segs, Segment{
                     URL:    base + su.Media,
                     Time:   0,
                     Number: idx + 1,
                  })
               }
            }

            if len(segs) > 0 {
               out[rep.ID] = segs
            }
         }
      }
   }

   enc := json.NewEncoder(os.Stdout)
   enc.SetIndent("", "  ")
   _ = enc.Encode(out)
}

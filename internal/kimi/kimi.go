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

// ---------- MPD structs ----------

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
   StartNumber     int              `xml:"startNumber,attr"`
   Timescale       int              `xml:"timescale,attr"`
   Duration        int              `xml:"duration,attr"`
   SegmentTimeline *SegmentTimeline `xml:"SegmentTimeline"`
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

// pickBaseURL returns the first BaseURL found at the given level, else ""
func pickBaseURL(list []BaseURL) string {
   if len(list) > 0 {
      return strings.TrimSpace(list[0].URL)
   }
   return ""
}

// buildBase walks the hierarchy and returns the final absolute base
func buildBase(mpdPath string, mpd *MPD, period *Period, as *AdaptationSet, rep *Representation) string {
   // 1-5 hierarchy
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
   // 5. derive from manifest file location
   if b == "" {
      u, _ := url.Parse(mpdPath)
      b = u.ResolveReference(&url.URL{Path: path.Dir(u.Path) + "/"}).String()
   }
   // ensure trailing slash
   if !strings.HasSuffix(b, "/") {
      b += "/"
   }
   // if relative, resolve against manifest directory
   if !strings.HasPrefix(b, "http") {
      u, _ := url.Parse(mpdPath)
      b = u.ResolveReference(&url.URL{Path: b}).String()
   }
   return b
}

// expand substitutes $Time$, $Number$, $Bandwidth$, $RepresentationID$
func expand(tpl string, vars map[string]string) string {
   for k, v := range vars {
      tpl = strings.ReplaceAll(tpl, "$"+k+"$", v)
   }
   return tpl
}

func segmentTimes(st *SegmentTemplate) []int {
   var times []int
   if st.SegmentTimeline == nil || len(st.SegmentTimeline.S) == 0 {
      return times
   }

   t := 0
   for _, s := range st.SegmentTimeline.S {
      if s.T != 0 {
         t = s.T
      }
      repeats := s.R
      if repeats == 0 {
         repeats = 1
      }
      for i := 0; i <= repeats; i++ {
         times = append(times, t)
         t += s.D
      }
   }
   // remove the trailing “next” time that was added one step too far
   return times[:len(times)-1]
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

   for _, period := range mpd.Periods {
      for _, as := range period.AdaptationSets {
         for _, rep := range as.Representations {
            base := buildBase(mpdFile, &mpd, &period, &as, &rep)

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
                     "Time":             strconv.Itoa(t), // <— correct
                  }
                  seg := expand(st.Media, vars)
                  fmt.Println(base + seg)
               }

            }

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

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

// ---------- MPD structs (only what we need) ----------

type MPD struct {
   XMLName        xml.Name        `xml:"MPD"`
   BaseURL        []BaseURL       `xml:"BaseURL"`
   Periods        []Period        `xml:"Period"`
   XMLBase        string          `xml:"base,attr"`
   MinBufferTime  string          `xml:"minBufferTime,attr"`
   Type           string          `xml:"type,attr"`
   MediaPresentationDuration string `xml:"mediaPresentationDuration,attr"`
}

type BaseURL struct {
   URL string `xml:",chardata"`
}

type Period struct {
   Duration string          `xml:"duration,attr"`
   AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
}

type AdaptationSet struct {
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   Representations []Representation `xml:"Representation"`
}

type Representation struct {
   ID         string `xml:"id,attr"`
   Bandwidth  string `xml:"bandwidth,attr"`
   BaseURL    string `xml:"BaseURL"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   SegmentList *SegmentList `xml:"SegmentList"`
}

type SegmentTemplate struct {
   Media      string `xml:"media,attr"`
   Initialization string `xml:"initialization,attr"`
   StartNumber int    `xml:"startNumber,attr"`
   Timescale   int    `xml:"timescale,attr"`
   Duration    int    `xml:"duration,attr"`
   SegmentTimeline *SegmentTimeline `xml:"SegmentTimeline"`
}

type SegmentTimeline struct {
   S []S `xml:"S"`
}

type S struct {
   T string `xml:"t,attr"`
   D int    `xml:"d,attr"`
   R int    `xml:"r,attr"` // number of repeats
}

type SegmentList struct {
   SegmentURLs []SegmentURL `xml:"SegmentURL"`
}

type SegmentURL struct {
   Media string `xml:"media,attr"`
}

// ---------- helpers ----------

// resolveBase returns the effective BaseURL for a given representation
func resolveBase(mp *MPD, rep *Representation) string {
   // 1. Representation-level BaseURL overrides everything
   if rep.BaseURL != "" && !strings.HasPrefix(rep.BaseURL, "http") {
      return rep.BaseURL
   }
   if rep.BaseURL != "" {
      return rep.BaseURL
   }
   // 2. MPD-level BaseURL
   if len(mp.BaseURL) > 0 && mp.BaseURL[0].URL != "" {
      return strings.TrimSpace(mp.BaseURL[0].URL)
   }
   return ""
}

// expandTemplate substitutes $Number$, $Bandwidth$, etc.
func expandTemplate(tpl string, repID, bandwidth string, number int) string {
   tpl = strings.ReplaceAll(tpl, "$RepresentationID$", repID)
   tpl = strings.ReplaceAll(tpl, "$Bandwidth$", bandwidth)
   tpl = strings.ReplaceAll(tpl, "$Number$", fmt.Sprintf("%d", number))
   return tpl
}

// ---------- main ----------

func main() {
   if len(os.Args) != 2 {
      fmt.Fprintf(os.Stderr, "Usage: %s <manifest.mpd>\n", os.Args[0])
      os.Exit(1)
   }

   mpdPath := os.Args[1]
   data, err := os.ReadFile(mpdPath)
   if err != nil {
      fmt.Fprintf(os.Stderr, "cannot read MPD: %v\n", err)
      os.Exit(1)
   }

   var mpd MPD
   if err := xml.Unmarshal(data, &mpd); err != nil {
      fmt.Fprintf(os.Stderr, "cannot parse MPD: %v\n", err)
      os.Exit(1)
   }

   // Default baseURL is the directory of the MPD itself
   mpdURL, _ := url.Parse(mpdPath)
   defaultBase := mpdURL.ResolveReference(&url.URL{Path: path.Dir(mpdURL.Path) + "/"}).String()

   for _, period := range mpd.Periods {
      for _, as := range period.AdaptationSets {
         for _, rep := range as.Representations {
            // 1. Determine BaseURL
            base := resolveBase(&mpd, &rep)
            if base == "" {
               base = defaultBase
            }
            if !strings.HasSuffix(base, "/") {
               base += "/"
            }

            // 2. Pick the effective SegmentTemplate
            st := as.SegmentTemplate
            if rep.SegmentTemplate != nil {
               st = rep.SegmentTemplate
            }

            // 3. Handle SegmentTemplate
            if st != nil {
               start := st.StartNumber
               if start == 0 {
                  start = 1
               }
               mediaTpl := st.Media

               // If SegmentTimeline present, use it
               if st.SegmentTimeline != nil && len(st.SegmentTimeline.S) > 0 {
                  number := start
                  for _, s := range st.SegmentTimeline.S {
                     repeats := s.R
                     if repeats == 0 {
                        repeats = 1
                     }
                     for i := 0; i < repeats; i++ {
                        seg := expandTemplate(mediaTpl, rep.ID, rep.Bandwidth, number)
                        fmt.Println(base + seg)
                        number++
                     }
                  }
               } else {
                  // Uniform duration
                  dur := st.Duration
                  timescale := st.Timescale
                  if timescale == 0 {
                     timescale = 1
                  }
                  // How many segments?
                  totalDur := parseDuration(mpd.MediaPresentationDuration)
                  segments := int(totalDur * float64(timescale) / float64(dur))
                  for n := start; n < start+segments; n++ {
                     seg := expandTemplate(mediaTpl, rep.ID, rep.Bandwidth, n)
                     fmt.Println(base + seg)
                  }
               }
               continue
            }

            // 4. Handle SegmentList
            if rep.SegmentList != nil {
               for _, su := range rep.SegmentList.SegmentURLs {
                  fmt.Println(base + su.Media)
               }
               continue
            }

            // 5. Single-segment (BaseURL only) – not relevant for segment list
         }
      }
   }
}

// parseDuration converts ISO8601 duration to seconds (very naive)
func parseDuration(d string) float64 {
   // PT1H2M30.5S => 1*3600 + 2*60 + 30.5
   if !strings.HasPrefix(d, "PT") {
      return 0
   }
   d = strings.TrimPrefix(d, "PT")
   hours := 0.0
   minutes := 0.0
   seconds := 0.0
   if idx := strings.Index(d, "H"); idx != -1 {
      hours, _ = strconv.ParseFloat(d[:idx], 64)
      d = d[idx+1:]
   }
   if idx := strings.Index(d, "M"); idx != -1 {
      minutes, _ = strconv.ParseFloat(d[:idx], 64)
      d = d[idx+1:]
   }
   if idx := strings.Index(d, "S"); idx != -1 {
      seconds, _ = strconv.ParseFloat(d[:idx], 64)
   }
   return hours*3600 + minutes*60 + seconds
}

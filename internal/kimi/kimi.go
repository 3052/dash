package main

import (
   "encoding/json"
   "encoding/xml"
   "fmt"
   "net/url"
   "os"
   "strconv"
   "strings"
)

type MPD struct {
   BaseURL string `xml:"BaseURL"`
   Period  struct {
      BaseURL       string `xml:"BaseURL"`
      Duration      string `xml:"duration,attr"`
      AdaptationSet []struct {
         Representation []struct {
            ID              string `xml:"id,attr"`
            BaseURL         string `xml:"BaseURL"`
            SegmentTemplate *struct {
               Media           string `xml:"media,attr"`
               Timescale       int    `xml:"timescale,attr"`
               Duration        int    `xml:"duration,attr"`
               EndNumber       int    `xml:"endNumber,attr"`
               SegmentTimeline *struct {
                  S []struct {
                     T int `xml:"t,attr"`
                     D int `xml:"d,attr"`
                     R int `xml:"r,attr"`
                  } `xml:"S"`
               } `xml:"SegmentTimeline"`
            } `xml:"SegmentTemplate"`
         } `xml:"Representation"`
      } `xml:"AdaptationSet"`
   } `xml:"Period"`
}

func resolve(base, ref string) string {
   u, _ := url.Parse(base)
   r, _ := url.Parse(ref)
   return u.ResolveReference(r).String()
}

func parseDuration(d string) float64 {
   d = strings.TrimPrefix(d, "PT")
   var h, m, s float64
   fmt.Sscanf(d, "%fH%fM%fS", &h, &m, &s)
   return h*3600 + m*60 + s
}

func main() {
   if len(os.Args) != 2 {
      fmt.Fprintln(os.Stderr, "usage: dash_urls file.mpd")
      os.Exit(1)
   }

   data, err := os.ReadFile(os.Args[1])
   if err != nil {
      fmt.Fprintln(os.Stderr, err)
      os.Exit(1)
   }

   // Strip namespace so we can use plain struct tags
   clean := strings.ReplaceAll(string(data), `xmlns="urn:mpeg:dash:schema:mpd:2011"`, "")
   dec := xml.NewDecoder(strings.NewReader(clean))

   var mpd MPD
   if err := dec.Decode(&mpd); err != nil {
      fmt.Fprintln(os.Stderr, "decode:", err)
      os.Exit(1)
   }

   base := "http://test.test/test.mpd"
   if mpd.BaseURL != "" {
      base = resolve(base, mpd.BaseURL)
   }
   if mpd.Period.BaseURL != "" {
      base = resolve(base, mpd.Period.BaseURL)
   }

   out := make(map[string][]string)

   for _, as := range mpd.Period.AdaptationSet {
      for _, rep := range as.Representation {
         repBase := base
         if rep.BaseURL != "" {
            repBase = resolve(repBase, rep.BaseURL)
         }

         st := rep.SegmentTemplate
         if st == nil {
            continue
         }

         var urls []string

         // 1. SegmentTimeline
         if st.SegmentTimeline != nil {
            seg := 1
            for _, s := range st.SegmentTimeline.S {
               for i := 0; i < 1+s.R; i++ {
                  urls = append(urls, resolve(repBase,
                     strings.ReplaceAll(st.Media, "$Number$", strconv.Itoa(seg))))
                  seg++
               }
            }
         } else {
            // 2. Constant duration
            periodSec := parseDuration(mpd.Period.Duration)
            periodScaled := int(periodSec * float64(st.Timescale))
            total := (periodScaled + st.Duration - 1) / st.Duration
            if st.EndNumber > 0 && total > st.EndNumber {
               total = st.EndNumber
            }
            for n := 1; n <= total; n++ {
               urls = append(urls, resolve(repBase,
                  strings.ReplaceAll(st.Media, "$Number$", strconv.Itoa(n))))
            }
         }

         if len(urls) > 0 {
            out[rep.ID] = urls
         }
      }
   }

   enc := json.NewEncoder(os.Stdout)
   enc.SetIndent("", "  ")
   _ = enc.Encode(out)
}

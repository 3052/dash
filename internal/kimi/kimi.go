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

/* ---------- XML types ---------- */

type MPD struct {
   BaseURL string `xml:"BaseURL"`
   Period  struct {
      AdaptationSet []struct {
         Representation []Representation `xml:"Representation"`
      } `xml:"AdaptationSet"`
   } `xml:"Period"`
}

type Representation struct {
   ID              string `xml:"id,attr"`
   BaseURL         string `xml:"BaseURL"`
   SegmentTemplate *struct {
      Media           string `xml:"media,attr"`
      Init            string `xml:"initialization,attr"`
      TimeScale       int    `xml:"timescale,attr"`
      SegmentTimeline *struct {
         S []struct {
            T int `xml:"t,attr"` // start time (optional)
            D int `xml:"d,attr"` // duration
            R int `xml:"r,attr"` // repeat count (optional)
         } `xml:"S"`
      } `xml:"SegmentTimeline"`
   } `xml:"SegmentTemplate"`
}

/* ---------- helpers ---------- */

func resolve(base, ref string) string {
   u, _ := url.Parse(base)
   r, _ := url.Parse(ref)
   return u.ResolveReference(r).String()
}

/* ---------- main ---------- */

func main() {
   if len(os.Args) != 2 {
      fmt.Fprintf(os.Stderr, "usage: %s <mpd>\n", os.Args[0])
      os.Exit(1)
   }

   // Read file
   data, err := os.ReadFile(os.Args[1])
   if err != nil {
      fmt.Fprintln(os.Stderr, err)
      os.Exit(1)
   }

   // Parse XML
   var mpd MPD
   if err := xml.Unmarshal(data, &mpd); err != nil {
      fmt.Fprintln(os.Stderr, err)
      os.Exit(1)
   }

   // Base URL
   base := "http://test.test/test.mpd"
   if mpd.BaseURL != "" {
      base = resolve(base, mpd.BaseURL)
   }

   // Collect URLs by Representation ID
   out := make(map[string][]string)

   for _, as := range mpd.Period.AdaptationSet {
      for _, rep := range as.Representation {
         repBase := base
         if rep.BaseURL != "" {
            repBase = resolve(repBase, rep.BaseURL)
         }

         // Case 1: single-file representation (BaseURL only)
         if rep.SegmentTemplate == nil && rep.BaseURL != "" {
            out[rep.ID] = []string{resolve(base, rep.BaseURL)}
            continue
         }

         st := rep.SegmentTemplate
         if st == nil || st.SegmentTimeline == nil {
            continue
         }

         // Build segment list
         var urls []string
         number := 1 // segment numbers start at 1
         for _, s := range st.SegmentTimeline.S {
            count := 1 + s.R
            for i := 0; i < count; i++ {
               media := strings.ReplaceAll(st.Media, "$Number$", strconv.Itoa(number))
               urls = append(urls, resolve(repBase, media))
               number++
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

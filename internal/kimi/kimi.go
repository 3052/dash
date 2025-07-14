package main

import (
   "encoding/json"
   "encoding/xml"
   "fmt"
   "net/url"
   "os"
   "strings"
)

const base = "http://test.test/test.mpd"

// ---------- XML structs ------------------------------------------------------

type MPD struct {
   XMLName xml.Name `xml:"MPD"`
   BaseURL string   `xml:"BaseURL"`
   Period  struct {
      BaseURL       string          `xml:"BaseURL"`
      AdaptationSet []AdaptationSet `xml:"AdaptationSet"`
   } `xml:"Period"`
}

type AdaptationSet struct {
   BaseURL         string           `xml:"BaseURL"`
   Representation  []Representation `xml:"Representation"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
}

type Representation struct {
   ID              string           `xml:"id,attr"`
   BaseURL         string           `xml:"BaseURL"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
}

type SegmentTemplate struct {
   Media           string           `xml:"media,attr"`
   Initialization  string           `xml:"initialization,attr"`
   Duration        int              `xml:"duration,attr"`
   Timescale       int              `xml:"timescale,attr"`
   StartNumber     int              `xml:"startNumber,attr"`
   EndNumber       int              `xml:"endNumber,attr"`
   SegmentTimeline *SegmentTimeline `xml:"SegmentTimeline"`
}

type SegmentTimeline struct {
   S []S `xml:"S"`
}
type S struct {
   T uint64 `xml:"t,attr"` // start time
   D uint64 `xml:"d,attr"` // duration
   R int    `xml:"r,attr"` // repeat count (default 0)
}

// ---------- helpers ----------------------------------------------------------

func join(baseURL, rel string) string {
   u, _ := url.Parse(baseURL)
   relParsed, _ := url.Parse(rel)
   return u.ResolveReference(relParsed).String()
}

func segmentURLs(baseURL string, tpl *SegmentTemplate, repID string) []string {
   if tpl == nil {
      return nil
   }
   mediaTpl := strings.ReplaceAll(tpl.Media, "$RepresentationID$", repID)
   initTpl := strings.ReplaceAll(tpl.Initialization, "$RepresentationID$", repID)

   var urls []string
   if init := join(baseURL, initTpl); init != "" {
      urls = append(urls, init)
   }

   // Timeline case
   if tpl.SegmentTimeline != nil {
      num := 1
      if tpl.StartNumber > 0 {
         num = tpl.StartNumber
      }
      for _, s := range tpl.SegmentTimeline.S {
         for i := 0; i <= s.R; i++ {
            media := strings.ReplaceAll(mediaTpl, "$Number$", fmt.Sprintf("%d", num))
            media = strings.ReplaceAll(media, "$Time$", fmt.Sprintf("%d", s.T+uint64(i)*s.D))
            urls = append(urls, join(baseURL, media))
            num++
         }
      }
      return urls
   }

   // Simple @duration/@endNumber case
   start := 1
   if tpl.StartNumber > 0 {
      start = tpl.StartNumber
   }
   end := tpl.EndNumber
   if end == 0 && tpl.Duration > 0 && tpl.Timescale > 0 {
      // 1-hour fallback if nothing else
      end = (3600 * tpl.Timescale) / tpl.Duration
   }
   for n := start; n <= end; n++ {
      media := strings.ReplaceAll(mediaTpl, "$Number$", fmt.Sprintf("%d", n))
      urls = append(urls, join(baseURL, media))
   }
   return urls
}

// ---------- main -------------------------------------------------------------

func main() {
   if len(os.Args) < 2 {
      fmt.Fprintf(os.Stderr, "usage: %s <mpd-file>\n", os.Args[0])
      os.Exit(1)
   }
   data, err := os.ReadFile(os.Args[1])
   if err != nil {
      panic(err)
   }
   var mpd MPD
   if err := xml.Unmarshal(data, &mpd); err != nil {
      panic(err)
   }

   // Build base URL chain
   mpdBase := strings.TrimSpace(mpd.BaseURL)
   periodBase := join(base, mpdBase)
   if p := strings.TrimSpace(mpd.Period.BaseURL); p != "" {
      periodBase = join(periodBase, p)
   }

   out := make(map[string][]string)

   for _, as := range mpd.Period.AdaptationSet {
      asBase := periodBase
      if a := strings.TrimSpace(as.BaseURL); a != "" {
         asBase = join(asBase, a)
      }

      for _, rep := range as.Representation {
         repBase := asBase
         if r := strings.TrimSpace(rep.BaseURL); r != "" {
            repBase = join(asBase, r)
         }

         switch {
         case rep.SegmentTemplate != nil:
            out[rep.ID] = segmentURLs(repBase, rep.SegmentTemplate, rep.ID)
         case as.SegmentTemplate != nil:
            out[rep.ID] = segmentURLs(repBase, as.SegmentTemplate, rep.ID)
         default:
            // Single segment
            out[rep.ID] = []string{repBase}
         }
      }
   }

   enc := json.NewEncoder(os.Stdout)
   enc.SetEscapeHTML(false)
   enc.SetIndent("", "  ")
   _ = enc.Encode(out)
}

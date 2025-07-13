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
   XMLName xml.Name `xml:"MPD"`
   BaseURL string   `xml:"BaseURL"`
   Period  Period   `xml:"Period"`
}

type Period struct {
   BaseURL        string          `xml:"BaseURL"`
   AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
}

type AdaptationSet struct {
   BaseURL         string           `xml:"BaseURL"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   Representations []Representation `xml:"Representation"`
}

type SegmentTemplate struct {
   Timescale     int              `xml:"timescale,attr"`
   Duration      int              `xml:"duration,attr"`
   EndNumber     int              `xml:"endNumber,attr"`
   Initialization string           `xml:"initialization,attr"`
   Media         string           `xml:"media,attr"`
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

type Representation struct {
   ID              string           `xml:"id,attr"`
   BaseURL         string           `xml:"BaseURL"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
}

func resolveURL(base string, rel string) string {
   baseURL, _ := url.Parse(base)
   relURL, _ := url.Parse(rel)
   return baseURL.ResolveReference(relURL).String()
}


func main() {
   if len(os.Args) < 2 {
      fmt.Println("Usage: go run extract_segments.go <path-to-mpd>")
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
   rootBase := resolveURL(baseMPD, mpd.BaseURL)
   result := make(map[string][]string)

   for _, as := range mpd.Period.AdaptationSets {
      asBase := resolveURL(rootBase, as.BaseURL)
      for _, rep := range as.Representations {
         fullBase := resolveURL(asBase, rep.BaseURL)
         repID := rep.ID
         segmentURLs := []string{}

         // Prefer SegmentTemplate inside Representation
         st := rep.SegmentTemplate
         if st == nil {
            st = as.SegmentTemplate
         }

         // Case 1: No SegmentTemplate, just a BaseURL (single file)
         if st == nil {
            if rep.BaseURL != "" {
               segmentURLs = append(segmentURLs, resolveURL(asBase, rep.BaseURL))
            }
            result[repID] = segmentURLs
            continue
         }

         // Initialization segment
         if st.Initialization != "" {
            initURL := strings.ReplaceAll(st.Initialization, "$RepresentationID$", repID)
            segmentURLs = append(segmentURLs, resolveURL(fullBase, initURL))
         }

         // SegmentTimeline-based segments
         if st.SegmentTimeline != nil {
            number := 0
            for _, seg := range st.SegmentTimeline.Segments {
               repeat := seg.R
               if repeat < 0 {
                  repeat = 0
               }
               for i := 0; i <= repeat; i++ {
                  number++
                  media := st.Media
                  media = strings.ReplaceAll(media, "$RepresentationID$", repID)
                  media = strings.ReplaceAll(media, "$Number$", strconv.Itoa(number))
                  segmentURLs = append(segmentURLs, resolveURL(fullBase, media))
               }
            }
         } else {
            // Number-based fallback
            for i := 1; i <= st.EndNumber; i++ {
               media := st.Media
               media = strings.ReplaceAll(media, "$RepresentationID$", repID)
               media = strings.ReplaceAll(media, "$Number$", strconv.Itoa(i))
               segmentURLs = append(segmentURLs, resolveURL(fullBase, media))
            }
         }

         result[repID] = segmentURLs
      }
   }

   enc, _ := json.MarshalIndent(result, "", "  ")
   fmt.Println(string(enc))
}

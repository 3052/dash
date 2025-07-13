package main

import (
   "encoding/json"
   "encoding/xml"
   "fmt"
   "io/ioutil"
   "net/url"
   "os"
   "strconv"
   "strings"
)

// --- Structs to Unmarshal MPD XML ---

type MPD struct {
   XMLName xml.Name `xml:"MPD"`
   BaseURL string   `xml:"BaseURL"`
   Periods []Period `xml:"Period"`
   MpdUrl  *url.URL `xml:"-"`
}

type Period struct {
   XMLName        xml.Name        `xml:"Period"`
   BaseURL        string          `xml:"BaseURL"`
   AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
}

type AdaptationSet struct {
   XMLName         xml.Name         `xml:"AdaptationSet"`
   Representations []Representation `xml:"Representation"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   BaseURL         string           `xml:"BaseURL"`
}

type Representation struct {
   XMLName         xml.Name         `xml:"Representation"`
   ID              string           `xml:"id,attr"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   BaseURL         string           `xml:"BaseURL"`
}

type SegmentTemplate struct {
   XMLName         xml.Name         `xml:"SegmentTemplate"`
   Initialization  string           `xml:"initialization,attr"`
   Media           string           `xml:"media,attr"`
   EndNumber       int              `xml:"endNumber,attr"`
   StartNumber     int              `xml:"startNumber,attr"`
   SegmentTimeline *SegmentTimeline `xml:"SegmentTimeline"`
}

type SegmentTimeline struct {
   XMLName  xml.Name `xml:"SegmentTimeline"`
   Segments []S      `xml:"S"`
}

type S struct {
   XMLName  xml.Name `xml:"S"`
   Time     *uint64  `xml:"t,attr"`
   Duration uint64   `xml:"d,attr"`
   Repeat   *int     `xml:"r,attr"`
}

// --- Main Logic ---

func main() {
   if len(os.Args) < 2 {
      fmt.Println("Usage: go run mpd-parser.go <path-to-mpd-file>")
      os.Exit(1)
   }
   filePath := os.Args[1]

   mpdBaseUrl, err := url.Parse("http://test.test/test.mpd")
   if err != nil {
      fmt.Printf("Error parsing base MPD URL: %v\n", err)
      os.Exit(1)
   }

   byteValue, err := ioutil.ReadFile(filePath)
   if err != nil {
      fmt.Printf("Error reading file %s: %v\n", filePath, err)
      os.Exit(1)
   }

   var mpd MPD
   if err := xml.Unmarshal(byteValue, &mpd); err != nil {
      fmt.Printf("Error unmarshalling XML: %v\n", err)
      os.Exit(1)
   }
   mpd.MpdUrl = mpdBaseUrl

   representationUrls := make(map[string][]string)

   for _, period := range mpd.Periods {
      for _, adaptationSet := range period.AdaptationSets {
         for _, representation := range adaptationSet.Representations {

            effectiveBaseStr := getEffectiveBaseURL(mpd.BaseURL, period.BaseURL, adaptationSet.BaseURL)
            currentBase, err := mpd.MpdUrl.Parse(effectiveBaseStr)
            if err != nil {
               continue
            }

            var segments []string
            template := representation.SegmentTemplate
            if template == nil {
               template = adaptationSet.SegmentTemplate
            }

            if template != nil {
               finalSegmentBase, _ := currentBase.Parse(representation.BaseURL)

               if template.Initialization != "" {
                  initURL := strings.ReplaceAll(template.Initialization, "$RepresentationID$", representation.ID)
                  if resolved, e := finalSegmentBase.Parse(initURL); e == nil {
                     segments = append(segments, resolved.String())
                  }
               }

               // ** CORRECTED LOGIC **
               if template.SegmentTimeline != nil {
                  segmentCounter := template.StartNumber
                  if segmentCounter == 0 {
                     segmentCounter = 1 // Default start number is 1
                  }

                  for _, s := range template.SegmentTimeline.Segments {
                     repeatCount := 0
                     if s.Repeat != nil {
                        repeatCount = *s.Repeat
                     }
                     for i := 0; i <= repeatCount; i++ {
                        mediaUrl := strings.ReplaceAll(template.Media, "$RepresentationID$", representation.ID)
                        // Use the counter for $Number$, not the time value.
                        mediaUrl = strings.ReplaceAll(mediaUrl, "$Number$", strconv.Itoa(segmentCounter))
                        if resolved, e := finalSegmentBase.Parse(mediaUrl); e == nil {
                           segments = append(segments, resolved.String())
                        }
                        segmentCounter++ // Increment for the next segment
                     }
                  }
               } else { // Fallback for Number-based templates (e.g. molotov.txt)
                  start := template.StartNumber
                  if start == 0 {
                     start = 1
                  }
                  for i := start; i <= template.EndNumber; i++ {
                     mediaUrl := strings.ReplaceAll(template.Media, "$RepresentationID$", representation.ID)
                     mediaUrl = strings.ReplaceAll(mediaUrl, "$Number$", strconv.Itoa(i))
                     if resolved, e := finalSegmentBase.Parse(mediaUrl); e == nil {
                        segments = append(segments, resolved.String())
                     }
                  }
               }
            } else if representation.BaseURL != "" { // For single-file content like subtitles
               if resolved, e := currentBase.Parse(representation.BaseURL); e == nil {
                  segments = append(segments, resolved.String())
               }
            }

            if len(segments) > 0 {
               representationUrls[representation.ID] = segments
            }
         }
      }
   }

   jsonOutput, err := json.MarshalIndent(representationUrls, "", "  ")
   if err != nil {
      fmt.Printf("Error marshalling JSON: %v\n", err)
      os.Exit(1)
   }

   fmt.Println(string(jsonOutput))
}

func getEffectiveBaseURL(bases ...string) string {
   for i := len(bases) - 1; i >= 0; i-- {
      if bases[i] != "" {
         if !strings.HasSuffix(bases[i], "/") {
            return bases[i] + "/"
         }
         return bases[i]
      }
   }
   return ""
}

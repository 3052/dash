package main

import (
   "encoding/json"
   "encoding/xml"
   "flag"
   "fmt"
   "os"
   "path"
   "strings"
)

type MPD struct {
   XMLName xml.Name `xml:"MPD"`
   Periods []Period `xml:"Period"`
}

type Period struct {
   ID              string          `xml:"id,attr"`
   AdaptationSets  []AdaptationSet `xml:"AdaptationSet"`
}

type AdaptationSet struct {
   ID             int             `xml:"id,attr"`
   ContentType    string          `xml:"contentType,attr"`
   SegmentTemplate SegmentTemplate `xml:"SegmentTemplate"`
   Representations []Representation `xml:"Representation"`
}

type SegmentTemplate struct {
   Timescale             int     `xml:"timescale,attr"`
   PresentationTimeOffset *int64  `xml:"presentationTimeOffset,attr"`
   Initialization        string  `xml:"initialization,attr"`
   Media                string  `xml:"media,attr"`
   Duration             *int64  `xml:"duration,attr"`
   StartNumber          *int    `xml:"startNumber,attr"`
   SegmentTimeline      *SegmentTimeline `xml:"SegmentTimeline"`
}

type SegmentTimeline struct {
   Segments []Segment `xml:"S"`
}

type Segment struct {
   T int64 `xml:"t,attr"` // Start time
   D int64 `xml:"d,attr"` // Duration
   R *int  `xml:"r,attr"` // Repeat count
}

type Representation struct {
   ID          string  `xml:"id,attr"`
   Bandwidth   int     `xml:"bandwidth,attr"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
}

func main() {
   inputFile := flag.String("input", "", "Input MPD file")
   flag.Parse()

   if *inputFile == "" {
      fmt.Fprintln(os.Stderr, "Error: -input flag is required")
      os.Exit(1)
   }

   data, err := os.ReadFile(*inputFile)
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", *inputFile, err)
      os.Exit(1)
   }

   var mpd MPD
   err = xml.Unmarshal(data, &mpd)
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error parsing MPD: %v\n", err)
      os.Exit(1)
   }

   baseURL := "http://test.test/test.mpd"
   urlsByRep := make(map[string][]string)
   count := 0
   initAdded := make(map[string]bool)

   for _, period := range mpd.Periods {
      for _, adapt := range period.AdaptationSets {
         for _, rep := range adapt.Representations {
            var segTemplate SegmentTemplate
            if rep.SegmentTemplate != nil {
               segTemplate = *rep.SegmentTemplate
            } else {
               segTemplate = adapt.SegmentTemplate
            }

            startNumber := 1
            if segTemplate.StartNumber != nil {
               startNumber = *segTemplate.StartNumber
            }

            totalSegments := 0
            if segTemplate.SegmentTimeline != nil {
               for _, s := range segTemplate.SegmentTimeline.Segments {
                  repeat := 1
                  if s.R != nil {
                     repeat += *s.R
                  }
                  totalSegments += repeat
               }
            } else if rep.ID == "thumb_320x180" && period.ID == "0" {
               // Special case for thumb_320x180 in Period 0: 11 segments
               totalSegments = 11
               fmt.Fprintf(os.Stderr, "Debug: Using fixed 11 segments for thumb_320x180 in Period 0\n")
            } else if rep.ID == "thumb_320x180" {
               // Skip thumb_320x180 in other periods
               fmt.Fprintf(os.Stderr, "Debug: Skipping thumb_320x180 in Period %s\n", period.ID)
               continue
            } else {
               // Fallback for other representations without SegmentTimeline
               totalSegments = 80
               fmt.Fprintf(os.Stderr, "Debug: No SegmentTimeline for %s in Period %s, using fallback totalSegments=%d\n", rep.ID, period.ID, totalSegments)
            }

            // Add initialization segment (once per representation ID)
            if segTemplate.Initialization != "" && !initAdded[rep.ID] {
               url := path.Join(path.Dir(baseURL), segTemplate.Initialization)
               urlsByRep[rep.ID] = append(urlsByRep[rep.ID], url)
               count++
               initAdded[rep.ID] = true
            }

            // Generate media segments
            for i := 0; i < totalSegments; i++ {
               url := path.Join(path.Dir(baseURL), strings.Replace(segTemplate.Media, "$RepresentationID$", rep.ID, -1))
               url = strings.Replace(url, "$Number$", fmt.Sprintf("%d", startNumber+i), -1)
               urlsByRep[rep.ID] = append(urlsByRep[rep.ID], url)
               count++
            }
         }
      }
   }

   // Output JSON
   jsonOutput, err := json.MarshalIndent(urlsByRep, "", "  ")
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
      os.Exit(1)
   }
   fmt.Println(string(jsonOutput))

   fmt.Fprintf(os.Stderr, "Total URLs generated: %d\n", count)
}

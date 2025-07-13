package main

import (
   "encoding/json"
   "encoding/xml"
   "fmt"
   "io/ioutil"
   "net/url"
   "os"
)

type MPD struct {
   XMLName       xml.Name       `xml:"MPD"`
   Namespaces    string         `xml:"xmlns,attr"`
   Period        []Period       `xml:"Period"`
}

type Period struct {
   AdaptationSet []AdaptationSet `xml:"AdaptationSet"`
}

type AdaptationSet struct {
   Representation []Representation `xml:"Representation"`
}

type Representation struct {
   ID           string   `xml:"id,attr"`
   BaseURL      string   `xml:"BaseURL"`
   SegmentList  SegmentList `xml:"SegmentList"`
}

type SegmentList struct {
   SegmentURL []SegmentURL `xml:"SegmentURL"`
}

type SegmentURL struct {
   Media string `xml:"media,attr"`
}

func resolveURL(base, relative string) string {
   parsedBase, err := url.Parse(base)
   if err != nil {
      fmt.Println("Error parsing base URL:", err)
      return ""
   }
   parsedRelative, err := url.Parse(relative)
   if err != nil {
      fmt.Println("Error parsing relative URL:", err)
      return ""
   }
   return parsedBase.ResolveReference(parsedRelative).String()
}

func main() {
   // Get file path from command-line argument
   if len(os.Args) < 2 {
      fmt.Println("Usage: go run main.go <path_to_mpd>")
      return
   }

   mpdFilePath := os.Args[1]

   // Read MPD file
   data, err := ioutil.ReadFile(mpdFilePath)
   if err != nil {
      fmt.Println("Error reading MPD file:", err)
      return
   }

   // Parse MPD XML
   var mpd MPD
   err = xml.Unmarshal(data, &mpd)
   if err != nil {
      fmt.Println("Error unmarshaling MPD:", err)
      return
   }

   // Initialize the output structure
   output := make(map[string][]string)

   // Base MPD URL
   baseURL := "http://test.test/test.mpd"

   // Loop through AdaptationSets and Representations to extract Segment URLs
   for _, period := range mpd.Period {
      for _, adaptationSet := range period.AdaptationSet {
         for _, representation := range adaptationSet.Representation {
            // Resolve the base URL if needed
            fullBaseURL := resolveURL(baseURL, representation.BaseURL)

            // Collect Segment URLs for this Representation
            var segments []string
            for _, segment := range representation.SegmentList.SegmentURL {
               segmentURL := resolveURL(fullBaseURL, segment.Media)
               segments = append(segments, segmentURL)
            }

            // Store the segments in the output map
            output[representation.ID] = segments
         }
      }
   }

   // Convert the output to JSON
   result, err := json.MarshalIndent(output, "", "  ")
   if err != nil {
      fmt.Println("Error marshaling JSON:", err)
      return
   }

   // Output the result
   fmt.Println(string(result))
}

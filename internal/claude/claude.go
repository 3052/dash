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

// MPD structures
type MPD struct {
   XMLName xml.Name `xml:"MPD"`
   Periods []Period `xml:"Period"`
}

type Period struct {
   AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
}

type AdaptationSet struct {
   Representations []Representation `xml:"Representation"`
}

type Representation struct {
   ID        string         `xml:"id,attr"`
   Bandwidth int            `xml:"bandwidth,attr"`
   MimeType  string         `xml:"mimeType,attr"`
   Codecs    string         `xml:"codecs,attr"`
   Width     int            `xml:"width,attr"`
   Height    int            `xml:"height,attr"`
   BaseURL   string         `xml:"BaseURL"`
   Template  SegmentTemplate `xml:"SegmentTemplate"`
   List      SegmentList     `xml:"SegmentList"`
}

type SegmentTemplate struct {
   Media          string `xml:"media,attr"`
   Initialization string `xml:"initialization,attr"`
   StartNumber    int    `xml:"startNumber,attr"`
   Duration       int    `xml:"duration,attr"`
   Timescale      int    `xml:"timescale,attr"`
   Timeline       SegmentTimeline `xml:"SegmentTimeline"`
}

type SegmentTimeline struct {
   S []TimelineSegment `xml:"S"`
}

type TimelineSegment struct {
   T int `xml:"t,attr"`
   D int `xml:"d,attr"`
   R int `xml:"r,attr"`
}

type SegmentList struct {
   Initialization SegmentURL   `xml:"Initialization"`
   SegmentURLs    []SegmentURL `xml:"SegmentURL"`
}

type SegmentURL struct {
   Media string `xml:"media,attr"`
}

// Output structures
type RepresentationOutput struct {
   ID        string   `json:"id"`
   Bandwidth int      `json:"bandwidth"`
   MimeType  string   `json:"mimeType"`
   Codecs    string   `json:"codecs"`
   Width     int      `json:"width"`
   Height    int      `json:"height"`
   URLs      []string `json:"urls"`
}

type Output struct {
   Representations []RepresentationOutput `json:"representations"`
}

func main() {
   if len(os.Args) != 2 {
      fmt.Fprintf(os.Stderr, "Usage: %s <path-to-mpd-file>\n", os.Args[0])
      os.Exit(1)
   }

   mpdPath := os.Args[1]
   baseURL := "http://test.test/test.mpd"

   // Read and parse MPD file
   data, err := ioutil.ReadFile(mpdPath)
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
      os.Exit(1)
   }

   var mpd MPD
   if err := xml.Unmarshal(data, &mpd); err != nil {
      fmt.Fprintf(os.Stderr, "Error parsing XML: %v\n", err)
      os.Exit(1)
   }

   // Extract segment URLs
   output := Output{
      Representations: []RepresentationOutput{},
   }

   for _, period := range mpd.Periods {
      for _, adaptationSet := range period.AdaptationSets {
         for _, rep := range adaptationSet.Representations {
            repOutput := RepresentationOutput{
               ID:        rep.ID,
               Bandwidth: rep.Bandwidth,
               MimeType:  rep.MimeType,
               Codecs:    rep.Codecs,
               Width:     rep.Width,
               Height:    rep.Height,
               URLs:      []string{},
            }

            // Get segment URLs based on the type of segmentation
            urls := extractSegmentURLs(rep, baseURL)
            repOutput.URLs = urls

            output.Representations = append(output.Representations, repOutput)
         }
      }
   }

   // Output JSON
   jsonData, err := json.MarshalIndent(output, "", "  ")
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
      os.Exit(1)
   }

   fmt.Println(string(jsonData))
}

func extractSegmentURLs(rep Representation, baseURL string) []string {
   var urls []string

   // Handle BaseURL - this represents a single segment/resource
   if rep.BaseURL != "" {
      fullURL := resolveURL(rep.BaseURL, baseURL)
      urls = append(urls, fullURL)
      return urls
   }

   // Handle SegmentTemplate
   if rep.Template.Media != "" {
      urls = extractFromTemplate(rep.Template, baseURL)
   } else if len(rep.List.SegmentURLs) > 0 {
      // Handle SegmentList
      urls = extractFromList(rep.List, baseURL)
   }

   return urls
}

func extractFromTemplate(template SegmentTemplate, baseURL string) []string {
   var urls []string

   // Add initialization segment if present
   if template.Initialization != "" {
      initURL := resolveURL(template.Initialization, baseURL)
      urls = append(urls, initURL)
   }

   // Generate media segments
   if template.Media != "" {
      if len(template.Timeline.S) > 0 {
         // Use timeline
         segmentNumber := template.StartNumber
         if segmentNumber == 0 {
            segmentNumber = 1
         }

         for _, s := range template.Timeline.S {
            repeat := 1
            if s.R > 0 {
               repeat = s.R + 1
            }

            for i := 0; i < repeat; i++ {
               mediaURL := strings.Replace(template.Media, "$Number$", strconv.Itoa(segmentNumber), -1)
               mediaURL = strings.Replace(mediaURL, "$Time$", strconv.Itoa(s.T), -1)
               fullURL := resolveURL(mediaURL, baseURL)
               urls = append(urls, fullURL)
               segmentNumber++
            }
         }
      } else {
         // Simple template without timeline - generate a few segments as example
         startNum := template.StartNumber
         if startNum == 0 {
            startNum = 1
         }

         // Generate 10 segments as example (in real scenario, you'd need duration info)
         for i := 0; i < 10; i++ {
            segmentNumber := startNum + i
            mediaURL := strings.Replace(template.Media, "$Number$", strconv.Itoa(segmentNumber), -1)
            fullURL := resolveURL(mediaURL, baseURL)
            urls = append(urls, fullURL)
         }
      }
   }

   return urls
}

func extractFromList(list SegmentList, baseURL string) []string {
   var urls []string

   // Add initialization segment if present
   if list.Initialization.Media != "" {
      initURL := resolveURL(list.Initialization.Media, baseURL)
      urls = append(urls, initURL)
   }

   // Add all segment URLs
   for _, segURL := range list.SegmentURLs {
      fullURL := resolveURL(segURL.Media, baseURL)
      urls = append(urls, fullURL)
   }

   return urls
}

func resolveURL(segmentURL, baseURL string) string {
   // If segmentURL is already absolute, return as is
   if strings.HasPrefix(segmentURL, "http://") || strings.HasPrefix(segmentURL, "https://") {
      return segmentURL
   }

   // Parse base URL
   base, err := url.Parse(baseURL)
   if err != nil {
      return segmentURL
   }

   // If segmentURL starts with /, it's absolute path
   if strings.HasPrefix(segmentURL, "/") {
      return base.Scheme + "://" + base.Host + segmentURL
   }

   // Otherwise, resolve relative to base URL directory
   // Use URL path manipulation to avoid backslashes
   basePath := base.Path
   lastSlash := strings.LastIndex(basePath, "/")
   if lastSlash >= 0 {
      basePath = basePath[:lastSlash]
   } else {
      basePath = ""
   }

   // Ensure proper path joining with forward slashes
   if basePath == "" || basePath == "/" {
      return base.Scheme + "://" + base.Host + "/" + segmentURL
   }
   
   return base.Scheme + "://" + base.Host + basePath + "/" + segmentURL
}

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

// MPD structures
type MPD struct {
   XMLName xml.Name `xml:"MPD"`
   Periods []Period `xml:"Period"`
}

type Period struct {
   AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
   BaseURL        string          `xml:"BaseURL"`
}

type AdaptationSet struct {
   Representations []Representation `xml:"Representation"`
   Template        SegmentTemplate   `xml:"SegmentTemplate"`
   BaseURL         string            `xml:"BaseURL"`
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
   EndNumber      int    `xml:"endNumber,attr"`
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
   data, err := os.ReadFile(mpdPath)
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

   // Use a map to combine URLs from same representation IDs across multiple periods
   repMap := make(map[string]*RepresentationOutput)

   for _, period := range mpd.Periods {
      for _, adaptationSet := range period.AdaptationSets {
         for _, rep := range adaptationSet.Representations {
            // Get segment URLs based on the type of segmentation
            // Pass adaptation set and period for inheritance
            urls := extractSegmentURLsWithInheritance(rep, adaptationSet, period, baseURL)

            // Check if we already have this representation ID
            if existing, exists := repMap[rep.ID]; exists {
               // Append URLs to existing representation, deduplicating
               existing.URLs = appendUniqueURLs(existing.URLs, urls)
            } else {
               // Create new representation
               repOutput := RepresentationOutput{
                  ID:        rep.ID,
                  Bandwidth: rep.Bandwidth,
                  MimeType:  rep.MimeType,
                  Codecs:    rep.Codecs,
                  Width:     rep.Width,
                  Height:    rep.Height,
                  URLs:      urls,
               }
               repMap[rep.ID] = &repOutput
            }
         }
      }
   }

   // Convert map to slice
   for _, rep := range repMap {
      output.Representations = append(output.Representations, *rep)
   }

   // Output JSON
   jsonData, err := json.MarshalIndent(output, "", "  ")
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
      os.Exit(1)
   }

   fmt.Println(string(jsonData))
}

func appendUniqueURLs(existing []string, newURLs []string) []string {
   // Create a map for fast lookup of existing URLs
   urlSet := make(map[string]bool)
   for _, url := range existing {
      urlSet[url] = true
   }

   // Add new URLs that don't already exist
   for _, url := range newURLs {
      if !urlSet[url] {
         existing = append(existing, url)
         urlSet[url] = true
      }
   }

   return existing
}

func extractSegmentURLsWithInheritance(rep Representation, adaptationSet AdaptationSet, period Period, baseURL string) []string {
   var urls []string

   // Build the effective base URL by combining Period BaseURL with the base URL
   effectiveBaseURL := baseURL
   if period.BaseURL != "" {
      effectiveBaseURL = resolveURL(period.BaseURL, baseURL)
   }

   // Handle BaseURL - this represents a single segment/resource
   if rep.BaseURL != "" {
      fullURL := resolveURL(rep.BaseURL, effectiveBaseURL)
      urls = append(urls, fullURL)
      return urls
   }

   // Check for SegmentTemplate at representation level first, then adaptation set level
   var template SegmentTemplate
   if rep.Template.Media != "" {
      template = rep.Template
   } else if adaptationSet.Template.Media != "" {
      template = adaptationSet.Template
   }

   if template.Media != "" {
      urls = extractFromTemplateWithRepID(template, effectiveBaseURL, rep.ID)
   } else if len(rep.List.SegmentURLs) > 0 {
      // Handle SegmentList
      urls = extractFromList(rep.List, effectiveBaseURL)
   }

   return urls
}


func extractFromTemplateWithRepID(template SegmentTemplate, baseURL string, repID string) []string {
   var urls []string

   // Add initialization segment if present
   if template.Initialization != "" {
      initURL := strings.Replace(template.Initialization, "$RepresentationID$", repID, -1)
      initURL = resolveURL(initURL, baseURL)
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
               mediaURL = strings.Replace(mediaURL, "$RepresentationID$", repID, -1)
               fullURL := resolveURL(mediaURL, baseURL)
               urls = append(urls, fullURL)
               segmentNumber++
            }
         }
      } else {
         // Use endNumber if available, otherwise generate a few segments
         startNum := template.StartNumber
         if startNum == 0 {
            startNum = 1
         }

         endNum := template.EndNumber
         if endNum == 0 {
            endNum = startNum + 9 // Generate 10 segments as fallback
         }

         for i := startNum; i <= endNum; i++ {
            mediaURL := strings.Replace(template.Media, "$Number$", strconv.Itoa(i), -1)
            mediaURL = strings.Replace(mediaURL, "$RepresentationID$", repID, -1)
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

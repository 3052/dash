package dash

import (
   "encoding/xml"
   "errors"
   "fmt"
   "math"
   "strconv"
   "strings"
)

// Representation represents the Representation element.
type Representation struct {
   XMLName           xml.Name         `xml:"Representation"`
   ID                string           `xml:"id,attr"`
   MimeType          string           `xml:"mimeType,attr,omitempty"`
   Bandwidth         int              `xml:"bandwidth,attr"`
   Codecs            string           `xml:"codecs,attr"`
   Width             int              `xml:"width,attr,omitempty"`
   Height            int              `xml:"height,attr,omitempty"`
   FrameRate         string           `xml:"frameRate,attr,omitempty"`
   AudioSamplingRate string           `xml:"audioSamplingRate,attr,omitempty"`
   BaseURL           string           `xml:"BaseURL,omitempty"`
   SegmentTemplate   *SegmentTemplate `xml:"SegmentTemplate,omitempty"`
   SegmentList       *SegmentList     `xml:"SegmentList"`
}

// ResolveURL takes a template string and replaces identifiers like
// $RepresentationID$ and $Bandwidth$ with values from the Representation.
func (r *Representation) ResolveURL(template string) string {
   if r == nil {
      return template
   }
   replacer := strings.NewReplacer(
      "$RepresentationID$", r.ID,
      "$Bandwidth$", strconv.Itoa(r.Bandwidth),
   )
   return replacer.Replace(template)
}

// ListMediaSegmentURLs generates a list of media segment URLs. The parent Period is
// required for manifests that calculate segment count from the Period duration.
func (r *Representation) ListMediaSegmentURLs(p *Period, asTpl *SegmentTemplate) ([]string, error) {
   tpl := r.SegmentTemplate
   if tpl == nil {
      tpl = asTpl
   }
   if tpl == nil {
      return nil, errors.New("no SegmentTemplate available for the representation")
   }

   start := tpl.StartNumber
   if start == 0 {
      start = 1
   }
   mediaURL := r.ResolveURL(tpl.Media)

   // --- Timeline-based logic (highest priority) ---
   if tpl.SegmentTimeline != nil {
      timeline := tpl.SegmentTimeline.GetSegments()
      urls := make([]string, 0, len(timeline))
      for i, segment := range timeline {
         num := start + uint(i)
         replacer := strings.NewReplacer(
            "$Number$", strconv.FormatUint(uint64(num), 10),
            "$Time$", strconv.FormatUint(segment.StartTime, 10),
         )
         segmentURL := replacer.Replace(mediaURL)
         urls = append(urls, segmentURL)
      }
      return urls, nil
   }

   // --- Number-based logic ---
   if tpl.EndNumber > 0 {
      urls := make([]string, 0, tpl.EndNumber-start+1)
      for num := start; num <= tpl.EndNumber; num++ {
         segmentURL := strings.ReplaceAll(mediaURL, "$Number$", strconv.FormatUint(uint64(num), 10))
         urls = append(urls, segmentURL)
      }
      return urls, nil
   }

   // --- Duration-based logic (last resort) ---
   if tpl.Duration > 0 && p != nil && p.Duration != "" {
      periodDurSec, err := p.AsSeconds()
      if err != nil {
         return nil, fmt.Errorf("failed to parse period duration: %w", err)
      }
      if tpl.Timescale <= 0 {
         return nil, errors.New("SegmentTemplate timescale must be positive")
      }
      segmentDurSec := float64(tpl.Duration) / float64(tpl.Timescale)
      if segmentDurSec <= 0 {
         return nil, errors.New("segment duration must be positive")
      }
      count := uint(math.Ceil(periodDurSec / segmentDurSec))
      urls := make([]string, 0, count)
      for i := uint(0); i < count; i++ {
         num := start + i
         segmentURL := strings.ReplaceAll(mediaURL, "$Number$", strconv.FormatUint(uint64(num), 10))
         urls = append(urls, segmentURL)
      }
      return urls, nil
   }

   return nil, errors.New("SegmentTemplate lacks sufficient information to generate segment list")
}

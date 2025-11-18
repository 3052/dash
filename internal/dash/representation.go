package dash

import (
   "encoding/xml"
   "errors"
   "strconv"
   "strings"
)

// Representation represents the Representation element.
type Representation struct {
   XMLName           xml.Name         `xml:"Representation"`
   ID                string           `xml:"id,attr"`
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

// ListMediaSegmentURLs generates a list of media segment URLs.
// It prioritizes SegmentTimeline if present, otherwise falls back to the
// Start/EndNumber logic.
// It uses the Representation's own SegmentTemplate first, then the one passed in
// (e.g., from the parent AdaptationSet).
func (r *Representation) ListMediaSegmentURLs(asTpl *SegmentTemplate) ([]string, error) {
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

   // Pre-resolve the parts of the URL that don't change per segment.
   mediaURL := r.ResolveURL(tpl.Media)

   // --- Timeline-based logic ---
   if tpl.SegmentTimeline != nil {
      timeline := tpl.SegmentTimeline.GetSegments()
      urls := make([]string, 0, len(timeline))
      for i, segment := range timeline { // CORRECTED THIS LINE
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

   return nil, errors.New("SegmentTemplate contains neither a SegmentTimeline nor an endNumber")
}

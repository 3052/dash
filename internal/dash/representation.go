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

// ResolveURL takes a template string (from a SegmentTemplate's initialization
// or media attribute) and replaces the DASH-defined identifiers with values
// from the Representation. It currently supports $RepresentationID$ and $Bandwidth$.
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

// ListMediaSegmentURLs generates a list of media segment URLs based on a SegmentTemplate.
// It uses the Representation's own SegmentTemplate if it exists, otherwise it falls back
// to the one passed in (e.g., from the parent AdaptationSet).
// The method returns an error if no valid SegmentTemplate with an EndNumber is found.
func (r *Representation) ListMediaSegmentURLs(asTpl *SegmentTemplate) ([]string, error) {
   tpl := r.SegmentTemplate
   if tpl == nil {
      tpl = asTpl
   }
   if tpl == nil {
      return nil, errors.New("no SegmentTemplate available for the representation")
   }
   if tpl.EndNumber == 0 {
      return nil, errors.New("SegmentTemplate does not contain an endNumber, cannot generate a finite list")
   }

   start := tpl.StartNumber
   if start == 0 {
      start = 1
   }

   // Pre-resolve the parts of the URL that don't change per segment.
   mediaURL := r.ResolveURL(tpl.Media)
   urls := make([]string, 0, tpl.EndNumber-start+1)

   for num := start; num <= tpl.EndNumber; num++ {
      // Note: This only supports the simple $Number$ identifier. A full implementation
      // would need to parse printf-style format specifiers (e.g., %05d).
      segmentURL := strings.ReplaceAll(mediaURL, "$Number$", strconv.FormatUint(uint64(num), 10))
      urls = append(urls, segmentURL)
   }

   return urls, nil
}

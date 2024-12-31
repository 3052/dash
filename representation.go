package dash

import (
   "iter"
   "strconv"
)

type Representation struct {
   Bandwidth       int64   `xml:"bandwidth,attr"`
   Codecs          *string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height          *int64  `xml:"height,attr"`
   Id              string  `xml:"id,attr"`
   MimeType        *string `xml:"mimeType,attr"`
   SegmentBase *struct {
      Initialization struct {
         Range Range `xml:"range,attr"`
      }
      IndexRange Range `xml:"indexRange,attr"`
   }
   SegmentTemplate *SegmentTemplate
   Width           *int64 `xml:"width,attr"`
   adaptation_set  *AdaptationSet
}

func (r *Representation) set() {
   if r.Codecs == nil {
      r.Codecs = r.adaptation_set.Codecs
   }
   if len(r.ContentProtection) == 0 {
      r.ContentProtection = r.adaptation_set.ContentProtection
   }
   if r.Height == nil {
      r.Height = r.adaptation_set.Height
   }
   if r.MimeType == nil {
      r.MimeType = r.adaptation_set.MimeType
   }
   if r.SegmentTemplate == nil {
      r.SegmentTemplate = r.adaptation_set.SegmentTemplate
   }
   // dashif.org/Guidelines-TimingModel#addressing-simple
   if r.SegmentTemplate != nil {
      if r.SegmentTemplate.StartNumber == nil {
         value := 1
         r.SegmentTemplate.StartNumber = &value
      }
   }
   // dashif.org/Guidelines-TimingModel#timing-sampletimeline
   if r.SegmentTemplate != nil {
      if r.SegmentTemplate.Timescale == nil {
         value := 1
         r.SegmentTemplate.Timescale = &value
      }
   }
   if r.Width == nil {
      r.Width = r.adaptation_set.Width
   }
}

func (r *Representation) String() string {
   var b []byte
   if r.Width != nil {
      b = append(b, "width = "...)
      b = strconv.AppendInt(b, *r.Width, 10)
   }
   if r.Height != nil {
      if b != nil {
         b = append(b, '\n')
      }
      b = append(b, "height = "...)
      b = strconv.AppendInt(b, *r.Height, 10)
   }
   if b != nil {
      b = append(b, '\n')
   }
   b = append(b, "bandwidth = "...)
   b = strconv.AppendInt(b, r.Bandwidth, 10)
   if r.Codecs != nil {
      b = append(b, "\ncodecs = "...)
      b = append(b, *r.Codecs...)
   }
   b = append(b, "\nmimeType = "...)
   b = append(b, *r.MimeType...)
   if role := r.adaptation_set.Role; role != nil {
      b = append(b, "\nrole = "...)
      b = append(b, role.Value...)
   }
   if lang := r.adaptation_set.Lang; lang != "" {
      b = append(b, "\nlang = "...)
      b = append(b, lang...)
   }
   if id := r.adaptation_set.period.Id; id != "" {
      b = append(b, "\nperiod = "...)
      b = append(b, id...)
   }
   b = append(b, "\nid = "...)
   b = append(b, r.Id...)
   return string(b)
}

func (r *Representation) seq() iter.Seq[Representation] {
   return func(yield func(Representation) bool) {
      for r2 := range r.adaptation_set.period.mpd.representation() {
         if r2.Id == r.Id {
            if !yield(r2) {
               return
            }
         }
      }
   }
}

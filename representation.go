package dash

import (
   "iter"
   "net/url"
   "strconv"
)

func (r *Representation) set(adapt *AdaptationSet) {
   r.adaptation_set = adapt
   if v := r.adaptation_set.period.BaseUrl; v != nil {
      if r.BaseUrl == nil {
         r.BaseUrl = &Url{&url.URL{}}
      }
      r.BaseUrl.Url = v.Url.ResolveReference(r.BaseUrl.Url)
   }
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
   if r.SegmentTemplate != nil {
      r.SegmentTemplate.set()
   }
   if r.Width == nil {
      r.Width = r.adaptation_set.Width
   }
}

type Representation struct {
   Bandwidth         int64   `xml:"bandwidth,attr"`
   BaseUrl           *Url    `xml:"BaseURL"`
   Codecs            *string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height            *int64  `xml:"height,attr"`
   Id                string  `xml:"id,attr"`
   MimeType          *string `xml:"mimeType,attr"`
   SegmentBase       *struct {
      Initialization struct {
         Range Range `xml:"range,attr"`
      }
      IndexRange Range `xml:"indexRange,attr"`
   }
   SegmentTemplate *SegmentTemplate
   Width           *int64 `xml:"width,attr"`
   adaptation_set  *AdaptationSet
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

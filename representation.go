package dash

import (
   "iter"
   "net/url"
   "strconv"
)

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
   b = strconv.AppendInt(b, int64(r.Bandwidth), 10)
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

type Representation struct {
   Bandwidth         int     `xml:"bandwidth,attr"`
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
   SegmentList *struct {
      Initialization struct {
         SourceUrl ListUrl `xml:"sourceURL,attr"`
      }
      SegmentUrl []struct {
         Media ListUrl `xml:"media,attr"`
      } `xml:"SegmentURL"`
   }
   SegmentTemplate *SegmentTemplate
   Width           *int64 `xml:"width,attr"`
   adaptation_set  *AdaptationSet
}

func (r *Representation) Segment() iter.Seq[int] {
   template := r.SegmentTemplate
   var address int
   if template.Media.time() {
      address = template.PresentationTimeOffset
   } else {
      address = *template.StartNumber
   }
   return func(yield func(int) bool) {
      if template.SegmentTimeline != nil {
         for _, segment := range template.SegmentTimeline.S {
            for range 1 + segment.R {
               if !yield(address) {
                  return
               }
               if template.Media.time() {
                  address += segment.D
               } else {
                  address++
               }
            }
         }
      } else {
         segment_count := r.adaptation_set.period.segment_count(template)
         for range int64(segment_count) {
            if !yield(address) {
               return
            }
            address++
         }
      }
   }
}

func (r *Representation) Representation() iter.Seq[Representation] {
   return func(yield func(Representation) bool) {
      for _, period0 := range r.adaptation_set.period.mpd.Period {
         for _, adapt := range period0.AdaptationSet {
            for _, represent := range adapt.Representation {
               if represent.Id == r.Id {
                  if adapt.period == nil {
                     period0.set(r.adaptation_set.period.mpd)
                     adapt.set(&period0)
                  }
                  represent.set(&adapt)
                  if !yield(represent) {
                     return
                  }
               }
            }
         }
      }
   }
}

func (r *Representation) set(adapt *AdaptationSet) {
   r.adaptation_set = adapt
   if base := r.adaptation_set.period.BaseUrl; base != nil {
      if r.BaseUrl == nil {
         r.BaseUrl = &Url{&url.URL{}}
      }
      r.BaseUrl.Url = base.Url.ResolveReference(r.BaseUrl.Url)
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

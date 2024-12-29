package dash

import (
   "iter"
   "strconv"
)

func (r *Representation) seq() iter.Seq[Representation] {
   return func(yield func(Representation) bool) {
      for rb := range r.adaptation_set.period.mpd.representation() {
         if rb.Id == r.Id {
            if !yield(rb) {
               return
            }
         }
      }
   }
}

func (m Mpd) representation() iter.Seq[Representation] {
   return func(yield func(Representation) bool) {
      for _, p := range m.Period {
         p.mpd = &m
         for _, adapt := range p.AdaptationSet {
            adapt.period = &p
            for _, represent := range adapt.Representation {
               if represent.Codecs == nil {
                  represent.Codecs = adapt.Codecs
               }
               if represent.Height == nil {
                  represent.Height = adapt.Height
               }
               if represent.MimeType == nil {
                  represent.MimeType = adapt.MimeType
               }
               if represent.SegmentTemplate == nil {
                  represent.SegmentTemplate = adapt.SegmentTemplate
               }
               if represent.Width == nil {
                  represent.Width = adapt.Width
               }
               represent.adaptation_set = &adapt
               if !yield(represent) {
                  return
               }
            }
         }
      }
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

type Mpd struct {
   Period []Period
}

type AdaptationSet struct {
   Codecs         *string `xml:"codecs,attr"`
   Height         *int64  `xml:"height,attr"`
   Lang           string  `xml:"lang,attr"`
   MimeType       *string `xml:"mimeType,attr"`
   Representation []Representation
   Role           *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
   Width  *int64 `xml:"width,attr"`
   period *Period
}

type Representation struct {
   Bandwidth      int64  `xml:"bandwidth,attr"`
   Codecs         *string `xml:"codecs,attr"`
   Height         *int64  `xml:"height,attr"`
   Id             string `xml:"id,attr"`
   MimeType       *string `xml:"mimeType,attr"`
   SegmentTemplate *SegmentTemplate
   Width          *int64  `xml:"width,attr"`
   adaptation_set *AdaptationSet
}

type Period struct {
   AdaptationSet []AdaptationSet
   Id            string `xml:"id,attr"`
   mpd *Mpd
}

type SegmentTemplate struct {
   SegmentTimeline *struct {
      S []struct {
         R int `xml:"r,attr"` // repeat
      }
   }
}

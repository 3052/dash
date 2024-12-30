package dash

import (
   "iter"
   "net/url"
   "strconv"
   "strings"
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

type Url struct {
   Url url.URL
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
   Width           *int64 `xml:"width,attr"`
   period          *Period
}

func (b *Url) UnmarshalText(data []byte) error {
   return b.Url.UnmarshalBinary(data)
}

type Period struct {
   AdaptationSet []AdaptationSet
   BaseUrl       *Url   `xml:"BaseURL"`
   Id            string `xml:"id,attr"`
   mpd           *Mpd
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

type Representation struct {
   Bandwidth       int64   `xml:"bandwidth,attr"`
   Codecs          *string `xml:"codecs,attr"`
   Height          *int64  `xml:"height,attr"`
   Id              string  `xml:"id,attr"`
   MimeType        *string `xml:"mimeType,attr"`
   SegmentTemplate *SegmentTemplate
   Width           *int64 `xml:"width,attr"`
   adaptation_set  *AdaptationSet
}

type SegmentTemplate struct {
   Initialization *Url `xml:"initialization,attr"`
}

func (r *Representation) initialization() (*Url, bool) {
   if v := r.SegmentTemplate; v != nil {
      if v := v.Initialization; v != nil {
         v.Url.Path = strings.Replace(v.Url.Path, "$RepresentationID$", r.Id, 1)
         return v, true
      }
   }
   return nil, false
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

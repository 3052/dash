package dash

import (
   "iter"
   "strconv"
)

type AdaptationSet struct {
   period         *Period
   Representation []Representation
   Role           *struct {
      Value string `xml:"value,attr"`
   }
   Codecs   *string `xml:"codecs,attr"`
   Height   *int64  `xml:"height,attr"`
   Lang     string  `xml:"lang,attr"`
   MimeType *string `xml:"mimeType,attr"`
   Width    *int64  `xml:"width,attr"`
}

func (m Mpd) representation() iter.Seq[Representation] {
   id := map[string]struct{}{}
   return func(yield func(Representation) bool) {
      for _, p := range m.Period {
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
               if represent.Width == nil {
                  represent.Width = adapt.Width
               }
               represent.adaptation_set = &adapt
               _, ok := id[represent.Id]
               if !ok {
                  if !yield(represent) {
                     return
                  }
               }
               id[represent.Id] = struct{}{}
            }
         }
      }
   }
}

type Mpd struct {
   Period []Period
}

type Period struct {
   AdaptationSet []AdaptationSet
   Id            string `xml:"id,attr"`
}

type Representation struct {
   Bandwidth      int64  `xml:"bandwidth,attr"`
   Id             string `xml:"id,attr"`
   adaptation_set *AdaptationSet
   Width          *int64  `xml:"width,attr"`
   Codecs         *string `xml:"codecs,attr"`
   Height         *int64  `xml:"height,attr"`
   MimeType       *string `xml:"mimeType,attr"`
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

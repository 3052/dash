package dash

import "iter"

func (r *Representation) String() string {
   var b []byte
   if v := r.get_width(); v >= 1 {
      b = append(b, "width = "...)
      b = strconv.AppendUint(b, v, 10)
   }
   if v := r.get_height(); v >= 1 {
      if b != nil {
         b = append(b, '\n')
      }
      b = append(b, "height = "...)
      b = strconv.AppendUint(b, v, 10)
   }
   if b != nil {
      b = append(b, '\n')
   }
   b = append(b, "bandwidth = "...)
   b = strconv.AppendUint(b, r.Bandwidth, 10)
   if v := r.get_codecs(); v != "" {
      b = append(b, "\ncodecs = "...)
      b = append(b, v...)
   }
   b = append(b, "\nmimeType = "...)
   b = append(b, r.GetMimeType()...)
   if v := r.adaptation_set.Role; v != nil {
      b = append(b, "\nrole = "...)
      b = append(b, v.Value...)
   }
   if v := r.adaptation_set.Lang; v != "" {
      b = append(b, "\nlang = "...)
      b = append(b, v...)
   }
   if v := r.adaptation_set.period.Id; v != "" {
      b = append(b, "\nperiod = "...)
      b = append(b, v...)
   }
   b = append(b, "\nid = "...)
   b = append(b, r.Id...)
   return string(b)
}

type Mpd struct {
   Period []Period
}

type Period struct {
   AdaptationSet []AdaptationSet
   Id            string `xml:"id,attr"`
}

type AdaptationSet struct {
   Codecs         string `xml:"codecs,attr"`
   Height         int  `xml:"height,attr"`
   Lang           string `xml:"lang,attr"`
   MimeType       string `xml:"mimeType,attr"`
   Representation []Representation
   Role           *struct {
      Value string `xml:"value,attr"`
   }
   Width  int `xml:"width,attr"`
   period *Period
}

type Representation struct {
   Width *int `xml:"width,attr"`
   Height         *int  `xml:"height,attr"`
   Bandwidth      int  `xml:"bandwidth,attr"`
   Codecs         *string `xml:"codecs,attr"`
   MimeType       *string `xml:"mimeType,attr"`
   
   Id             string `xml:"id,attr"`
   adaptation_set *AdaptationSet
}

func (m Mpd) representation() iter.Seq[Representation] {
   id := map[string]struct{}{}
   return func(yield func(Representation) bool) {
      for _, p := range m.Period {
         for _, adapt := range p.AdaptationSet {
            adapt.period = &p
            for _, represent := range adapt.Representation {
               if represent.Codecs == nil {
                  represent.Codecs = &adapt.Codecs
               }
               if represent.Height == nil {
                  represent.Height = &adapt.Height
               }
               if represent.MimeType == nil {
                  represent.MimeType = &adapt.MimeType
               }
               if represent.Width == nil {
                  represent.Width = &adapt.Width
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

package dash

import "iter"

func (p *Period) representation() iter.Seq[Representation] {
   return func(yield func(Representation) bool) {
      for _, adapt := range p.AdaptationSet {
         adapt.period = p
         for _, represent := range adapt.Representation {
            represent.adaptation_set = &adapt
            if !yield(represent) {
               return
            }
         }
      }
   }
}

type AdaptationSet struct {
   Codecs         string `xml:"codecs,attr"`
   Height         int64  `xml:"height,attr"`
   Lang           string `xml:"lang,attr"`
   MimeType       string `xml:"mimeType,attr"`
   Representation []Representation
   Role           *struct {
      Value string `xml:"value,attr"`
   }
   Width  int64 `xml:"width,attr"`
   period *Period
}

type Representation struct {
   Bandwidth      int64  `xml:"bandwidth,attr"`
   adaptation_set *AdaptationSet
   Codecs         string `xml:"codecs,attr"`
   Height         int64  `xml:"height,attr"`
   Id             string `xml:"id,attr"`
   MimeType       string `xml:"mimeType,attr"`
   Width          int64  `xml:"width,attr"`
}

type Period struct {
   AdaptationSet []AdaptationSet
   Id            string `xml:"id,attr"`
}

type Mpd struct {
   Period []Period
}

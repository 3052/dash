package dash

type Mpd struct {
   BaseUrl                   string  `xml:"BaseURL"`
   Period                    []Period
}

type Period struct {
   mpd           *Mpd
   BaseUrl       string  `xml:"BaseURL"`
   Id            string    `xml:"id,attr"`
   AdaptationSet []AdaptationSet
}

type AdaptationSet struct {
   Codecs            string `xml:"codecs,attr"`
   Height            uint64 `xml:"height,attr"`
   Lang              string `xml:"lang,attr"`
   MaxHeight         int    `xml:"maxHeight,attr"`
   MaxWidth          int    `xml:"maxWidth,attr"`
   MimeType          string `xml:"mimeType,attr"`
   Role              *struct {
      Value string `xml:"value,attr"`
   }
   Width           uint64 `xml:"width,attr"`
   Representation    []Representation
}

type Representation struct {
   Bandwidth         uint64   `xml:"bandwidth,attr"`
   BaseUrl           string `xml:"BaseURL"`
   Codecs            string   `xml:"codecs,attr"`
   Height            uint64 `xml:"height,attr"`
   Id                string `xml:"id,attr"`
   MimeType          string `xml:"mimeType,attr"`
   Width             uint64 `xml:"width,attr"`
   SegmentBase       *struct {
      Initialization struct {
         Range string `xml:"range,attr"`
      }
      IndexRange string `xml:"indexRange,attr"`
   }
   adaptation_set    *AdaptationSet
}

package dash

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

package dash

func (m mpd) representation() []representation {
   var represent []representation
   for _, period := range m.Period {
      for _, adapt := range period.AdaptationSet {
         represent = append(represent, adapt.Representation...)
      }
   }
   return represent
}

type mpd struct {
   Period []struct {
      Id            string `xml:"id,attr"`
      AdaptationSet []struct {
         Width uint64 `xml:"width,attr"`
         Height         uint64 `xml:"height,attr"`
         Codecs         string `xml:"codecs,attr"`
         Lang           string `xml:"lang,attr"`
         MimeType       string `xml:"mimeType,attr"`
         Role           *struct {
            Value string `xml:"value,attr"`
         }
         Representation []representation
      }
   }
}

type representation struct {
   Bandwidth uint64 `xml:"bandwidth,attr"`
   Width     uint64 `xml:"width,attr"`
   Height    uint64 `xml:"height,attr"`
   Codecs    string `xml:"codecs,attr"`
   Id        string `xml:"id,attr"`
   MimeType  string `xml:"mimeType,attr"`
}

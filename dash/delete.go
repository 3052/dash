package dash

type Representation struct {
   adaptation_set *adaptation_set
   Bandwidth int64 `xml:"bandwidth,attr"`
   BaseURL string
   Codecs string `xml:"codecs,attr"`
   Height int64 `xml:"height,attr"`
   ID string `xml:"id,attr"`
   MimeType string `xml:"mimeType,attr"`
   SegmentBase *struct {
      Initialization struct {
         Range Range `xml:"range,attr"`
      }
      IndexRange Range `xml:"indexRange,attr"`
   }
   SegmentTemplate *SegmentTemplate
   Width int64 `xml:"width,attr"`
}

func (r Representation) Ext() (string, bool) {
   switch r.GetMimeType() {
   case "audio/mp4":
      return ".m4a", true
   case "video/mp4":
      return ".m4v", true
   }
   return "", false
}

type adaptation_set struct {
   period *period
   Codecs string `xml:"codecs,attr"`
   Lang string `xml:"lang,attr"`
   MimeType string `xml:"mimeType,attr"`
   Representation []Representation
   Role *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
}

type period struct {
   mpd *mpd
   AdaptationSet []adaptation_set
}

package dash

// stream represents a simplified view of a media stream's characteristics,
// combining information typically found across Period, AdaptationSet, and
// Representation types in a DASH MPD.
type stream struct {
   Bandwidth int
}

type Representation struct {
   Bandwidth       int     `xml:"bandwidth,attr"`
   Codecs          *string `xml:"codecs,attr"`
   Id              string  `xml:"id,attr"`
   MimeType        *string `xml:"mimeType,attr"`
   Width           *int    `xml:"width,attr"`
   Height          *int    `xml:"height,attr"`
   SegmentTemplate *SegmentTemplate
   SegmentBase     *struct {
      Initialization struct {
         Range string `xml:"range,attr"`
      }
      IndexRange string `xml:"indexRange,attr"`
   }
   BaseUrl string `xml:"BaseURL"`
}

type Mpd struct {
   BaseUrl                   string `xml:"BaseURL"`
   MediaPresentationDuration string `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

type Period struct {
   AdaptationSet []AdaptationSet
   BaseUrl       string `xml:"BaseURL"`
   Duration      string `xml:"duration,attr"`
   Id            string `xml:"id,attr"`
}

type AdaptationSet struct {
   Codecs         *string `xml:"codecs,attr"`
   Height         *int    `xml:"height,attr"`
   Lang           string  `xml:"lang,attr"`
   MimeType       string  `xml:"mimeType,attr"`
   Representation []Representation
   Role           *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
   Width           *int `xml:"width,attr"`
}

type SegmentTemplate struct {
   Duration               int    `xml:"duration,attr"`
   EndNumber              int    `xml:"endNumber,attr"`
   Initialization         string `xml:"initialization,attr"`
   Media                  string `xml:"media,attr"`
   PresentationTimeOffset int    `xml:"presentationTimeOffset,attr"`
   SegmentTimeline        *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R int `xml:"r,attr"` // repeat
      }
   }
   StartNumber *int `xml:"startNumber,attr"`
   // This can be any frequency but typically is the media clock frequency of
   // one of the media streams (or a positive integer multiple thereof).
   Timescale *int `xml:"timescale,attr"`
}

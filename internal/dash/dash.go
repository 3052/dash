package dash

import (
   "net/url"
   "time"
)

type Mpd struct {
   BaseUrl                   Url      `xml:"BaseURL"`
   MediaPresentationDuration Duration `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

type Url [1]*url.URL

type Duration [1]time.Duration

type Period struct {
   BaseUrl       Url       `xml:"BaseURL"`
   Id            string    `xml:"id,attr"`
   Duration      *Duration `xml:"duration,attr"`
   
   AdaptationSet []AdaptationSet
}

type AdaptationSet struct {
   Lang              string `xml:"lang,attr"`
   MimeType          string `xml:"mimeType,attr"`
   Role              *struct {
      Value string `xml:"value,attr"`
   }
   Codecs *string `xml:"codecs,attr"`
   Height *int    `xml:"height,attr"`
   Width  *int    `xml:"width,attr"`
   ContentProtection []ContentProtection
   Representation    []Representation
   SegmentTemplate *SegmentTemplate
}

type ContentProtection struct {
   Pssh        string `xml:"pssh"`
   SchemeIdUri string `xml:"schemeIdUri,attr"`
}

type Initialization string

type Media string

type SegmentTemplate struct {
   EndNumber              int            `xml:"endNumber,attr"`
   Initialization         Initialization `xml:"initialization,attr"`
   Media                  Media          `xml:"media,attr"`
   PresentationTimeOffset int            `xml:"presentationTimeOffset,attr"`
   SegmentTimeline        *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R int `xml:"r,attr"` // repeat
      }
   }
   StartNumber *int `xml:"startNumber,attr"`
   Duration    int  `xml:"duration,attr"`
   // This can be any frequency but typically is the media clock frequency of
   // one of the media streams (or a positive integer multiple thereof).
   Timescale *int `xml:"timescale,attr"`
}

type SegmentList struct {
   Initialization struct {
      SourceUrl Url `xml:"sourceURL,attr"`
   }
   SegmentUrl []*struct {
      Media Url `xml:"media,attr"`
   } `xml:"SegmentURL"`
}

type Representation struct {
   Bandwidth         int     `xml:"bandwidth,attr"`
   Codecs            *string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Id                string  `xml:"id,attr"`
   MimeType          *string `xml:"mimeType,attr"`
   Width             *int    `xml:"width,attr"`
   Height            *int    `xml:"height,attr"`
   SegmentTemplate   *SegmentTemplate
   SegmentBase       *struct {
      Initialization struct {
         Range string `xml:"range,attr"`
      }
      IndexRange string `xml:"indexRange,attr"`
   }
   BaseUrl           Url     `xml:"BaseURL"`
   SegmentList       *SegmentList
}

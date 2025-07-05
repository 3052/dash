package dash

import (
   "fmt"
   "iter"
   "math"
   "net/url"
   "strings"
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

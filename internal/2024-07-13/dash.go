package dash

import (
   "net/url"
   "time"
)

type Representation struct {
   Bandwidth         uint64 `xml:"bandwidth,attr"`
   BaseUrl           *BaseUrl   `xml:"BaseURL"`
   Codecs            string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height            uint64 `xml:"height,attr"`
   Id                string `xml:"id,attr"`
   MimeType          string `xml:"mimeType,attr"`
   SegmentBase       *SegmentBase
   SegmentTemplate   *SegmentTemplate
   Width             uint64 `xml:"width,attr"`
}

type AdaptationSet struct {
   Codecs            string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height            uint64 `xml:"height,attr"`
   Lang              string `xml:"lang,attr"`
   MimeType          string `xml:"mimeType,attr"`
   Representation    []Representation
   Role              *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
   Width           uint64 `xml:"width,attr"`
}

type BaseUrl struct {
   Url *url.URL
}

type ContentProtection struct {
   Pssh        Pssh   `xml:"pssh"`
   SchemeIdUri string `xml:"schemeIdUri,attr"`
}

type Duration struct {
   Duration time.Duration
}

type Mpd struct {
   BaseUrl *BaseUrl `xml:"BaseURL"`
   MediaPresentationDuration *Duration `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

type Period struct {
   AdaptationSet []AdaptationSet
   BaseUrl *BaseUrl `xml:"BaseURL"`
   Duration      *Duration `xml:"duration,attr"`
   Id            string    `xml:"id,attr"`
}

type Pssh []byte

// SegmentIndexBox uses:
// unsigned int(32) subsegment_duration;
// but range values can exceed 32 bits
type Range struct {
   Start uint64
   End   uint64
}

type SegmentBase struct {
   Initialization struct {
      Range Range `xml:"range,attr"`
   }
   IndexRange Range `xml:"indexRange,attr"`
}

type SegmentTemplate struct {
   StartNumber            uint   `xml:"startNumber,attr"`
   Duration               uint64 `xml:"duration,attr"`
   Initialization         string `xml:"initialization,attr"`
   Media                  string `xml:"media,attr"`
   PresentationTimeOffset uint   `xml:"presentationTimeOffset,attr"`
   Timescale              uint64 `xml:"timescale,attr"`
   SegmentTimeline        *struct {
      S []struct {
         D uint `xml:"d,attr"` // duration
         R uint `xml:"r,attr"` // repeat
      }
   }
}

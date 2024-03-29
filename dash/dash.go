package dash

import (
   "fmt"
   "strconv"
   "strings"
   "time"
)

type RawRange string

func (r RawRange) Scan() (*Range, error) {
   var v Range
   _, err := fmt.Sscanf(string(r), "%v-%v", &v.Start, &v.End)
   if err != nil {
      return nil, err
   }
   return &v, nil
}

type Range struct {
   Start uint64
   End uint64
}

type SegmentTemplate struct {
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   SegmentTimeline *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R int `xml:"r,attr"` // repeat
      }
   }
   StartNumber *int `xml:"startNumber,attr"`
}

func (s SegmentTemplate) replace(old string, number int) string {
   return strings.Replace(s.Media, old, strconv.Itoa(number), 1)
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

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#addressing-simple-to-explicit
type mpd struct {
   MediaPresentationDuration string `xml:"mediaPresentationDuration,attr"`
   Period []period
}

func (m mpd) seconds() (float64, error) {
   s := strings.TrimPrefix(m.MediaPresentationDuration, "PT")
   duration, err := time.ParseDuration(strings.ToLower(s))
   if err != nil {
      return 0, err
   }
   return duration.Seconds(), nil
}

type period struct {
   mpd *mpd
   AdaptationSet []adaptation_set
}

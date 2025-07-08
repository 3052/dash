package dash

import (
   "math"
   "strings"
   "time"
)

type Representation struct {
   Bandwidth       int     `xml:"bandwidth,attr"`
   Codecs          *string `xml:"codecs,attr"`
   Id              string  `xml:"id,attr"`
   MimeType        *string `xml:"mimeType,attr"`
   Width           *int    `xml:"width,attr"`
   Height          *int    `xml:"height,attr"`
   BaseUrl         string  `xml:"BaseURL"`
   SegmentTemplate *SegmentTemplate
   SegmentList     *SegmentList
   SegmentBase     *SegmentBase
}

type SegmentBase struct {
   Initialization struct {
      Range string `xml:"range,attr"`
   }
   IndexRange string `xml:"indexRange,attr"`
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

type Duration [1]time.Duration

type Mpd struct {
   BaseUrl                   string `xml:"BaseURL"`
   MediaPresentationDuration string `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

type Period struct {
   AdaptationSet []AdaptationSet
   BaseUrl       string   `xml:"BaseURL"`
   Duration      Duration `xml:"duration,attr"`
   Id            string   `xml:"id,attr"`
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

// with current data this always uses number addressing
func (s *SegmentTemplate) byEndNumber() []int {
   var numbers []int
   number := *s.StartNumber
   for number <= s.EndNumber {
      numbers = append(numbers, number)
      number++
   }
   return numbers
}

func (s *SegmentTemplate) byTimelineNumber() []int {
   var numbers []int
   number := *s.StartNumber
   for _, segment := range s.SegmentTimeline.S {
      for range 1 + segment.R {
         numbers = append(numbers, number)
         number++
      }
   }
   return numbers
}

// with current data this always uses number addressing
func (s *SegmentTemplate) byPeriod(periodVar *Period) []int {
   var numbers []int
   // dashif.org/Guidelines-TimingModel#addressing-simple-to-explicit
   // SegmentCount = Ceil(
   //    AsSeconds(Period@duration) /
   //    (SegmentTemplate@duration / SegmentTemplate@timescale)
   // )
   segmentCount := int64(math.Ceil(
      periodVar.Duration[0].Seconds() /
         (float64(s.Duration) / float64(*s.Timescale)),
   ))
   number := *s.StartNumber
   for range segmentCount {
      numbers = append(numbers, number)
      number++
   }
   return numbers
}

func (s *SegmentTemplate) byTimelineTime() []int {
   var numbers []int
   number := s.PresentationTimeOffset
   for _, segment := range s.SegmentTimeline.S {
      for range 1 + segment.R {
         numbers = append(numbers, number)
         number += segment.D
      }
   }
   return numbers
}

type SegmentList struct {
   Initialization struct {
      SourceUrl string `xml:"sourceURL,attr"`
   }
   SegmentUrl []*struct {
      Media string `xml:"media,attr"`
   } `xml:"SegmentURL"`
}

func (s *SegmentTemplate) Numbers(periodVar *Period) []int {
   if s.EndNumber >= 1 {
      return s.byEndNumber()
   }
   if s.SegmentTimeline != nil {
      if strings.Contains(s.Media, "$Time$") {
         return s.byTimelineTime()
      }
      return s.byTimelineNumber()
   }
   return s.byPeriod(periodVar)
}

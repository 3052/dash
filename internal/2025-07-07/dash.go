package dash

import (
   "log"
   "math"
   "strings"
   "time"
)

type SegmentBase struct {
   Initialization struct {
      Range string `xml:"range,attr"`
   }
   IndexRange string `xml:"indexRange,attr"`
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

type Representation struct {
   Codecs          string `xml:"codecs,attr"`
   Id              string `xml:"id,attr"`
   MimeType        string `xml:"mimeType,attr"`
   BaseUrl         string `xml:"BaseURL"`
   Bandwidth       uint   `xml:"bandwidth,attr"`
   Width           uint   `xml:"width,attr"`
   Height          uint   `xml:"height,attr"`
   SegmentTemplate *SegmentTemplate
   SegmentList     *SegmentList
   SegmentBase     *SegmentBase
}

type SegmentList struct {
   Initialization struct {
      SourceUrl string `xml:"sourceURL,attr"`
   }
   SegmentUrl []struct {
      Media string `xml:"media,attr"`
   } `xml:"SegmentURL"`
}

type SegmentTemplate struct {
   Duration               uint   `xml:"duration,attr"`
   EndNumber              uint   `xml:"endNumber,attr"`
   Initialization         string `xml:"initialization,attr"`
   Media                  string `xml:"media,attr"`
   PresentationTimeOffset uint   `xml:"presentationTimeOffset,attr"`
   SegmentTimeline        *struct {
      S []struct {
         D uint `xml:"d,attr"` // duration
         R uint `xml:"r,attr"` // repeat
      }
   }
   StartNumber *uint `xml:"startNumber,attr"`
   Timescale   *uint `xml:"timescale,attr"`
}

type AdaptationSet struct {
   Codecs         string `xml:"codecs,attr"`
   Height         uint   `xml:"height,attr"`
   Lang           string `xml:"lang,attr"`
   MimeType       string `xml:"mimeType,attr"`
   Representation []Representation
   Role           *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
   Width           uint `xml:"width,attr"`
}

func (*SegmentBase) segments() []string {
   return nil
}

func (s *SegmentList) segments() []string {
   var segments []string
   for _, segment := range s.SegmentUrl {
      segments = append(segments, segment.Media)
   }
   return segments
}

func (r *Representation) Segments(
   adapt *AdaptationSet, periodVar *Period,
) []string {
   if r.SegmentBase != nil {
      return r.SegmentBase.segments()
   }
   if r.SegmentList != nil {
      return r.SegmentList.segments()
   }
   if r.SegmentTemplate != nil {
      return r.SegmentTemplate.segments(periodVar)
   }
   return adapt.SegmentTemplate.segments(periodVar)
}

type segmentParameter struct {
   number           uint
   RepresentationId string
   time             uint
}

///

func (s *SegmentTemplate) byEndNumber() []segmentParameter {
   var segment []segmentParameter
   number := *s.StartNumber
   for number <= s.EndNumber {
      segment = append(segment, segmentParameter{number: number})
      number++
   }
   return segment
}

func (s *SegmentTemplate) byTimelineNumber() []segmentParameter {
   var segments []segmentParameter
   number := *s.StartNumber
   for _, segment := range s.SegmentTimeline.S {
      for range 1 + segment.R {
         segments = append(segments, segmentParameter{number: number})
         number++
      }
   }
   return segments
}

func (s *SegmentTemplate) byTimelineTime() []segmentParameter {
   var segments []segmentParameter
   number := s.PresentationTimeOffset
   for _, segment := range s.SegmentTimeline.S {
      for range 1 + segment.R {
         segments = append(segments, segmentParameter{time: number})
         number += segment.D
      }
   }
   return segments
}

func (s *SegmentTemplate) byPeriod(periodVar *Period) []segmentParameter {
   var segment []segmentParameter
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
      segment = append(segment, segmentParameter{number: number})
      number++
   }
   return segment
}

func (s *SegmentTemplate) segmentParameter(periodVar *Period) []segmentParameter {
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

func (s *SegmentTemplate) segments(periodVar *Period) []string {
   var segments []string
   for _, segment := range s.segmentParameter(periodVar) {
      log.Print(segment)
   }
   return segments
}

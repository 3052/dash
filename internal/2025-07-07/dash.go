package dash

import (
   "log"
   "math"
   "strings"
   "time"
)

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

///

// 1. iterate Mpd.Period
// 2. iterate Period.AdaptationSet
// 3. iterate AdaptationSet.Representation
// 4. if Representation.SegmentBase
// 5. if Representation.SegmentList
// 6. if Representation.SegmentTemplate
// 7. if AdaptationSet.SegmentTemplate

func (s *SegmentTemplate) byEndNumber() []uint {
   var segment []uint
   number := *s.StartNumber
   for number <= s.EndNumber {
      segment = append(segment, number)
      number++
   }
   return segment
}

func (s *SegmentTemplate) byTimelineNumber() []uint {
   var segments []uint
   number := *s.StartNumber
   for _, segment := range s.SegmentTimeline.S {
      for range 1 + segment.R {
         segments = append(segments, number)
         number++
      }
   }
   return segments
}

func (s *SegmentTemplate) byTimelineTime() []uint {
   var segments []uint
   number := s.PresentationTimeOffset
   for _, segment := range s.SegmentTimeline.S {
      for range 1 + segment.R {
         segments = append(segments, number)
         number += segment.D
      }
   }
   return segments
}

func (s *SegmentTemplate) byPeriod(periodVar *Period) []uint {
   var segment []uint
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
      segment = append(segment, number)
      number++
   }
   return segment
}

func (s *SegmentTemplate) numberTime(periodVar *Period) []uint {
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
   for _, segment := range s.numberTime(periodVar) {
      log.Print(segment)
   }
   return segments
}

package dash

import (
   "log"
   "math"
   "strings"
   "time"
)

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
   return s.byPeriod(periodVar) // SegmentTemplate.duration
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

type Duration [1]time.Duration

type Mpd struct {
   MediaPresentationDuration Duration `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

type Period struct {
   AdaptationSet []AdaptationSet
   Duration      Duration `xml:"duration,attr"`
}

type AdaptationSet struct {
   Representation  []Representation
   SegmentTemplate *SegmentTemplate
}

type Representation struct {
   Id              string `xml:"id,attr"`
   SegmentBase     *SegmentBase
   SegmentList     *SegmentList
   SegmentTemplate *SegmentTemplate
}

type SegmentBase struct {
   IndexRange string `xml:"indexRange,attr"`
}

type SegmentList struct {
   SegmentUrl []struct {
      Media string `xml:"media,attr"`
   } `xml:"SegmentURL"`
}

type SegmentTemplate struct {
   Duration               uint   `xml:"duration,attr"`
   EndNumber              uint   `xml:"endNumber,attr"`
   Media                  string `xml:"media,attr"`
   PresentationTimeOffset uint   `xml:"presentationTimeOffset,attr"`
   SegmentTimeline        *SegmentTimeline
   StartNumber            *uint `xml:"startNumber,attr"`
   Timescale              *uint `xml:"timescale,attr"`
}

type SegmentTimeline struct {
   S []struct {
      D uint `xml:"d,attr"` // duration
      R uint `xml:"r,attr"` // repeat
   }
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

func (s *SegmentTemplate) segments(periodVar *Period) []string {
   var segments []string
   for _, segment := range s.numberTime(periodVar) {
      log.Print(segment)
   }
   return segments
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

func (s *SegmentTemplate) byEndNumber() []uint {
   var segment []uint
   number := *s.StartNumber
   for number <= s.EndNumber {
      segment = append(segment, number)
      number++
   }
   return segment
}

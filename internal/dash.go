package dash

import (
   "math"
   "strings"
   "time"
)

// Period represents a DASH Period element.
type Period struct {
   Duration Duration `xml:"duration,attr"`
}

// SegmentTemplate represents a DASH SegmentTemplate element.
type SegmentTemplate struct {
   EndNumber        int `xml:"endNumber,attr"`
   Media            Media `xml:"media,attr"`
   PresentationTimeOffset int `xml:"presentationTimeOffset,attr"`
   SegmentTimeline  *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R int `xml:"r,attr"` // repeat
      }
   } `xml:"SegmentTimeline"`
   StartNumber *int `xml:"startNumber,attr"`
   Duration    int `xml:"duration,attr"`
   Timescale   *int `xml:"timescale,attr"`
}

// Duration represents a duration in time.Duration.
type Duration [1]time.Duration

// Media represents the media attribute string.
type Media string

// time_address checks if the media string contains "$Time$".
func (m Media) time_address() bool {
   return strings.Contains(string(m), "$Time$")
}

func (s *SegmentTemplate) byEndNumber(address int) []int {
   var segments []int
   for address <= s.EndNumber {
      segments = append(segments, address)
      address++
   }
   return segments
}

func (s *SegmentTemplate) byTimeline(address int) []int {
   var segments []int
   for _, segment := range s.SegmentTimeline.S {
      for range 1 + segment.R {
         segments = append(segments, address)
         if s.Media.time_address() {
            address += segment.D
         } else {
            address++
         }
      }
   }
   return segments
}

// dashif.org/Guidelines-TimingModel#addressing-simple-to-explicit
// SegmentCount = Ceil(
//    AsSeconds(Period@duration) /
//    (SegmentTemplate@duration / SegmentTemplate@timescale)
// )
func (s *SegmentTemplate) byPeriodDuration(address int, periodVar *Period) []int {
   segmentCount := int64(math.Ceil(
      periodVar.Duration[0].Seconds() /
      (float64(s.Duration) / float64(*s.Timescale)),
   ))
   var segments []int
   for range segmentCount {
      segments = append(segments, address)
      address++
   }
   return segments
}

func (s *SegmentTemplate) Segment(periodVar *Period) []int {
   var address int
   if s.Media.time_address() {
      address = s.PresentationTimeOffset
   } else {
      address = *s.StartNumber
   }
   var segments []int
   if s.EndNumber >= 1 {
      segments = s.byEndNumber(address)
   } else if s.SegmentTimeline != nil {
      segments = s.byTimeline(address)
   } else {
      segments = s.byPeriodDuration(address, periodVar)
   }
   return segments
}

package dash

import (
   "iter"
   "math"
   "strings"
   "time"
)

type Duration [1]time.Duration

type Media string

func (m Media) time_address() bool {
   return strings.Contains(string(m), "$Time$")
}

func (p *Period) segment_count(template *SegmentTemplate) int64 {
   durationVar := float64(template.Duration) / float64(*template.Timescale)
   return int64(math.Ceil(p.Duration[0].Seconds() / durationVar))
}

type Period struct {
   Id            string    `xml:"id,attr"`
   BaseUrl       string       `xml:"BaseURL"`
   Duration      Duration `xml:"duration,attr"`
}

type SegmentTemplate struct {
   EndNumber              int            `xml:"endNumber,attr"`
   Initialization         string `xml:"initialization,attr"`
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
   Timescale *int `xml:"timescale,attr"`
}

func (s *SegmentTemplate) Segment(periodVar *Period) iter.Seq[int] {
   var address int
   if s.Media.time_address() {
      address = s.PresentationTimeOffset
   } else {
      address = *s.StartNumber
   }
   return func(yield func(int) bool) {
      if s.EndNumber >= 1 {
         for address <= s.EndNumber {
            if !yield(address) {
               return
            }
            address++
         }
      } else if s.SegmentTimeline != nil {
         for _, segment := range s.SegmentTimeline.S {
            for range 1 + segment.R {
               if !yield(address) {
                  return
               }
               if s.Media.time_address() {
                  address += segment.D
               } else {
                  address++
               }
            }
         }
      } else {
         for range periodVar.segment_count(s) {
            if !yield(address) {
               return
            }
            address++
         }
      }
   }
}

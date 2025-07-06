package dash

import "iter"

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

func (s *SegmentTemplate) set() {
   // dashif.org/Guidelines-TimingModel#addressing-simple
   if s.StartNumber == nil {
      start := 1
      s.StartNumber = &start
   }
   // dashif.org/Guidelines-TimingModel#timing-sampletimeline
   if s.Timescale == nil {
      scale := 1
      s.Timescale = &scale
   }
}

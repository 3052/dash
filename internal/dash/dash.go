package dash

import (
   "iter"
   "math"
   "strings"
   "time"
)

func (m *Mpd) Stream() []Stream {
   return nil
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

type Representation struct {
   Bandwidth       int     `xml:"bandwidth,attr"`
   Codecs          *string `xml:"codecs,attr"`
   Id              string  `xml:"id,attr"`
   MimeType        *string `xml:"mimeType,attr"`
   Width           *int    `xml:"width,attr"`
   Height          *int    `xml:"height,attr"`
   SegmentTemplate *SegmentTemplate
   SegmentBase     *struct {
      Initialization struct {
         Range string `xml:"range,attr"`
      }
      IndexRange string `xml:"indexRange,attr"`
   }
   BaseUrl string `xml:"BaseURL"`
}

// with current data this always uses number addressing
func (s *SegmentTemplate) byEndNumber() iter.Seq[int] {
   return func(yield func(int) bool) {
      number := *s.StartNumber
      for number <= s.EndNumber {
         if !yield(number) {
            return
         }
         number++
      }
   }
}

func (s *SegmentTemplate) byTimelineNumber() iter.Seq[int] {
   return func(yield func(int) bool) {
      number := *s.StartNumber
      for _, segment := range s.SegmentTimeline.S {
         for range 1 + segment.R {
            if !yield(number) {
               return
            }
            number++
         }
      }
   }
}

func (s *SegmentTemplate) byTimelineTime() iter.Seq[int] {
   return func(yield func(int) bool) {
      number := s.PresentationTimeOffset
      for _, segment := range s.SegmentTimeline.S {
         for range 1 + segment.R {
            if !yield(number) {
               return
            }
            number += segment.D
         }
      }
   }
}

// with current data this always uses number addressing
func (s *SegmentTemplate) byPeriod(periodVar *Period) iter.Seq[int] {
   // dashif.org/Guidelines-TimingModel#addressing-simple-to-explicit
   // SegmentCount = Ceil(
   //    AsSeconds(Period@duration) /
   //    (SegmentTemplate@duration / SegmentTemplate@timescale)
   // )
   segmentCount := int64(math.Ceil(
      periodVar.Duration[0].Seconds() /
         (float64(s.Duration) / float64(*s.Timescale)),
   ))
   return func(yield func(int) bool) {
      number := *s.StartNumber
      for range segmentCount {
         if !yield(number) {
            return
         }
         number++
      }
   }
}

type Duration [1]time.Duration

type SegmentTemplate struct {
   Duration               int    `xml:"duration,attr"`
   EndNumber              int    `xml:"endNumber,attr"`
   Initialization         string `xml:"initialization,attr"`
   Media                  string  `xml:"media,attr"`
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

// stream represents a simplified view of a media stream's characteristics,
// combining information typically found across Period, AdaptationSet, and
// Representation types in a DASH MPD.
type Stream struct {
   Bandwidth int
   Segment   []string
}

func (s *SegmentTemplate) Segments(periodVar *Period) iter.Seq[int] {
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


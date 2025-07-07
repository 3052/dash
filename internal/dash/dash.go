package dash

import (
   "math"
   "strings"
   "time"
)

// Media represents the media attribute string.
func (m Media) time_address() bool {
   return strings.Contains(string(m), "$Time$")
}

// Media represents the media attribute string.
type Media string

type Duration [1]time.Duration

func (s *SegmentTemplate) Segment(periodVar *Period) []int {
   if s.EndNumber >= 1 {
      return s.byEndNumber()
   }
   if s.SegmentTimeline != nil {
      if s.Media.time_address() {
         return s.byTimelineTimeAddress()
      }
      return s.byTimelineNumberAddress()
   }
   return s.byPeriodDuration(periodVar)
}

// byTimelineTimeAddress generates segments based on time addressing.
func (s *SegmentTemplate) byTimelineTimeAddress() []int {
   var segments []int
   number := s.PresentationTimeOffset
   for _, segment := range s.SegmentTimeline.S {
      for range 1 + segment.R {
         segments = append(segments, number)
         number += segment.D
      }
   }
   return segments
}

// byTimelineNumberAddress generates segments based on number addressing.
func (s *SegmentTemplate) byTimelineNumberAddress() []int {
   var segments []int
   number := *s.StartNumber
   for _, segment := range s.SegmentTimeline.S {
      for range 1 + segment.R {
         segments = append(segments, number)
         number++
      }
   }
   return segments
}

// with current data this always uses number addressing
func (s *SegmentTemplate) byEndNumber() []int {
   var segments []int
   number := *s.StartNumber
   for number <= s.EndNumber {
      segments = append(segments, number)
      number++
   }
   return segments
}

// with current data this always uses number addressing
func (s *SegmentTemplate) byPeriodDuration(periodVar *Period) []int {
   var segments []int
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
      segments = append(segments, number)
      number++
   }
   return segments
}

func (*Mpd) Stream() []Stream {
   return nil
}

// stream represents a simplified view of a media stream's characteristics,
// combining information typically found across Period, AdaptationSet, and
// Representation types in a DASH MPD.
type Stream struct {
   Bandwidth int
   Segment   []string
}

type Mpd struct {
   BaseUrl                   string `xml:"BaseURL"`
   MediaPresentationDuration string `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
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

type SegmentTemplate struct {
   Duration               int    `xml:"duration,attr"`
   EndNumber              int    `xml:"endNumber,attr"`
   Initialization         string `xml:"initialization,attr"`
   Media                  Media  `xml:"media,attr"`
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

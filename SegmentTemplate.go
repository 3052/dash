package dash

import (
   "fmt"
   "math"
   "strings"
)

func (s SegmentTemplate) GetInitialization(r *Representation) (string, bool) {
   if v := s.Initialization; v != "" {
      return r.id(v), true
   }
   return "", false
}

type SegmentTemplate struct {
   Duration float64 `xml:"duration,attr"`
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   StartNumber int `xml:"startNumber,attr"`
   PresentationTimeOffset int `xml:"presentationTimeOffset,attr"`
   Timescale float64 `xml:"timescale,attr"`
   SegmentTimeline *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R int `xml:"r,attr"` // repeat
      }
   }
}

func (s SegmentTemplate) start() int {
   if v := s.PresentationTimeOffset; v >= 1 {
      return v
   }
   return s.StartNumber
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#timing-sampletimeline
func (s SegmentTemplate) get_timescale() float64 {
   if v := s.Timescale; v >= 1 {
      return v
   }
   return 1
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#addressing-simple-to-explicit
func (s SegmentTemplate) segment_count(seconds float64) float64 {
   seconds /= s.Duration / s.get_timescale()
   return math.Ceil(seconds)
}

func (s SegmentTemplate) number(value int) string {
   f := strings.Replace(s.Media, "$Number$", "%d", 1)
   f = strings.Replace(f, "$Number%02d$", "%02d", 1)
   f = strings.Replace(f, "$Number%03d$", "%03d", 1)
   f = strings.Replace(f, "$Number%04d$", "%04d", 1)
   f = strings.Replace(f, "$Number%05d$", "%05d", 1)
   f = strings.Replace(f, "$Number%06d$", "%06d", 1)
   f = strings.Replace(f, "$Number%07d$", "%07d", 1)
   f = strings.Replace(f, "$Number%08d$", "%08d", 1)
   f = strings.Replace(f, "$Number%09d$", "%09d", 1)
   return fmt.Sprintf(f, value)
}

func (s SegmentTemplate) time(value int) string {
   f := strings.Replace(s.Media, "$Time$", "%d", 1)
   return fmt.Sprintf(f, value)
}

func (s SegmentTemplate) GetMedia(r *Representation) ([]string, error) {
   s.Media = r.id(s.Media)
   var media []string
   number := s.start()
   if s.SegmentTimeline != nil {
      for _, segment := range s.SegmentTimeline.S {
         var repeat int
         if segment.R >= 1 {
            repeat = segment.R
         }
         for range 1 + repeat {
            var medium string
            if strings.Contains(s.Media, "$Time$") {
               medium = s.time(number)
               number += segment.D
            } else {
               medium = s.number(number)
               number++
            }
            media = append(media, medium)
         }
      }
   } else {
      seconds, err := r.adaptation_set.period.Seconds()
      if err != nil {
         return nil, err
      }
      for range int(s.segment_count(seconds)) {
         media = append(media, s.number(number))
         number++
      }
   }
   return media, nil
}

func (r Representation) id(value string) string {
   return strings.Replace(value, "$RepresentationID$", r.Id, 1)
}

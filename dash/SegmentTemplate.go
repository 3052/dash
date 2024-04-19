package dash

import (
   "fmt"
   "math"
   "strings"
)

type SegmentTemplate struct {
   Duration *float64 `xml:"duration,attr"`
   Initialization *string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   SegmentTimeline *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R *int `xml:"r,attr"` // repeat
      }
   }
   StartNumber *int `xml:"startNumber,attr"`
   Timescale *float64 `xml:"timescale,attr"`
}

func (s SegmentTemplate) GetInitialization(r *Representation) (string, bool) {
   if v := s.Initialization; v != nil {
      return strings.Replace(*v, "$RepresentationID$", r.ID, 1), true
   }
   return "", false
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#timing-sampletimeline
func (s SegmentTemplate) get_timescale() float64 {
   if v := s.Timescale; v != nil {
      return *v
   }
   return 1
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#addressing-simple-to-explicit
func (s SegmentTemplate) segment_count(seconds float64) float64 {
   seconds /= *s.Duration / s.get_timescale()
   return math.Ceil(seconds)
}

func (s SegmentTemplate) GetMedia(r *Representation) ([]string, error) {
   s.Media = strings.Replace(s.Media, "$RepresentationID$", r.ID, 1)
   var (
      media []string
      number int
   )
   if s.StartNumber != nil {
      number = *s.StartNumber
   }
   if s.SegmentTimeline != nil {
      for _, segment := range s.SegmentTimeline.S {
         var repeat int
         if segment.R != nil {
            repeat = *segment.R
         }
         for range 1 + repeat {
            var medium string
            if s.StartNumber != nil {
               medium = sprintf(s.Media, number)
               number++
            } else {
               medium = strings.Replace(s.Media, "$Time$", fmt.Sprint(number), 1)
               number += segment.D
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
         medium := strings.Replace(s.Media, "$Number$", fmt.Sprint(number), 1)
         media = append(media, medium)
         number++
      }
   }
   return media, nil
}

// github.com/Dash-Industry-Forum/DASH-IF-Conformance/blob/development/Utils/impl/MPDHandler/computeUrls.php
func sprintf(format string, number int) string {
   format = strings.Replace(format, "$Number$", fmt.Sprint(number), 1)
   if strings.Contains(format, "$Number") {
      format = fmt.Sprintf(format, number)
      format = strings.Replace(format, "$Number", "", 1)
      format = strings.Replace(format, "$", "", 1)
   }
   return format
}

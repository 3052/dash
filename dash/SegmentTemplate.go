package dash

import (
   "fmt"
   "math"
   "strings"
)

// github.com/Dash-Industry-Forum/DASH-IF-Conformance/blob/development/Utils/impl/MPDHandler/computeUrls.php
func replace(s, old string, number int) string {
   s = strings.Replace(s, old, fmt.Sprint(number), 1)
   if old = strings.TrimSuffix(old, "$"); strings.Contains(s, old) {
      s = fmt.Sprintf(s, number)
      s = strings.Replace(s, old, "", 1)
      s = strings.Replace(s, "$", "", 1)
   }
   return s
}

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
               medium = replace(s.Media, "$Number$", number)
               number++
            } else {
               medium = replace(s.Media, "$Time$", number)
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
         medium := replace(s.Media, "$Number$", number)
         media = append(media, medium)
         number++
      }
   }
   return media, nil
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

package dash

import (
   "fmt"
   "math"
   "strings"
)

func (s SegmentTemplate) GetMedia(r *Representation) ([]string, error) {
   s.Media = strings.Replace(s.Media, "$RepresentationID$", r.Id, 1)
   var media []string
   number := s.start()
   if s.SegmentTimeline != nil {
      for _, segment := range s.SegmentTimeline.S {
         var repeat int
         if segment.R != nil {
            repeat = *segment.R
         }
         for range 1 + repeat {
            var medium string
            switch {
            case media_attr(s.Media).number():
               medium = replace(s.Media, "$Number", "$", number)
               number++
            case media_attr(s.Media).time():
               medium = replace(s.Media, "$Time", "$", number)
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
         medium := replace(s.Media, "$Number", "$", number)
         media = append(media, medium)
         number++
      }
   }
   return media, nil
}

func (m media_attr) time() bool {
   return strings.Contains(string(m), "$Time$")
}

func (m media_attr) number() bool {
   numbers := []string{
      "$Number$",
      "$Number%02d$",
      "$Number%03d$",
      "$Number%04d$",
      "$Number%05d$",
      "$Number%06d$",
      "$Number%07d$",
      "$Number%08d$",
      "$Number%09d$",
   }
   for _, number := range numbers {
      if strings.Contains(string(m), number) {
         return true
      }
   }
   return false
}

// github.com/Dash-Industry-Forum/DASH-IF-Conformance/blob/development/Utils/impl/MPDHandler/computeUrls.php
func replace(s, before, after string, number int) string {
   s = strings.Replace(s, before+after, fmt.Sprint(number), 1)
   if strings.Contains(s, before) {
      s = fmt.Sprintf(s, number)
      s = strings.Replace(s, before, "", 1)
      s = strings.Replace(s, after, "", 1)
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
   Timescale *float64 `xml:"timescale,attr"`
   StartNumber *int `xml:"startNumber,attr"`
   PresentationTimeOffset int `xml:"presentationTimeOffset,attr"`
}

type media_attr string

func (s SegmentTemplate) start() int {
   if v := s.PresentationTimeOffset; v >= 1 {
      return v
   }
   if v := s.StartNumber; v != nil {
      return *v
   }
   return 0
}

func (s SegmentTemplate) GetInitialization(r *Representation) (string, bool) {
   if v := s.Initialization; v != nil {
      return strings.Replace(*v, "$RepresentationID$", r.Id, 1), true
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

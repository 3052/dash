package dash

import (
   "strconv"
   "strings"
)

func (s SegmentTemplate) GetMedia(id string) []string {
   timeline := s.SegmentTimeline
   if timeline == nil {
      return nil
   }
   s.Media = strings.Replace(s.Media, "$RepresentationID$", id, 1)
   var number int
   if s.StartNumber != nil {
      number = *s.StartNumber
   }
   var media []string
   for _, segment := range timeline.S {
      var repeat int
      if segment.R != nil {
         repeat = *segment.R
      }
      for repeat >= 0 {
         var medium string
         replace := strconv.Itoa(number)
         if s.StartNumber != nil {
            medium = strings.Replace(s.Media, "$Number$", replace, 1)
            number++
         } else {
            medium = strings.Replace(s.Media, "$Time$", replace, 1)
            number += segment.D
         }
         media = append(media, medium)
         repeat--
      }
   }
   return media
}

type SegmentTemplate struct {
   Initialization *string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   SegmentTimeline *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R *int `xml:"r,attr"` // repeat
      }
   }
   StartNumber *int `xml:"startNumber,attr"`
}

func (s SegmentTemplate) GetInitialization(id string) (string, bool) {
   if v := s.Initialization; v != nil {
      return strings.Replace(*v, "$RepresentationID$", id, 1), true
   }
   return "", false
}

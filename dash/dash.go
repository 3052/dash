package dash

import (
   "fmt"
   "strconv"
   "strings"
   "time"
)

type Range string

// range-start and range-end can both exceed 32 bits, so we must use 64 bit
func (r Range) Scan() (uint64, uint64, error) {
   var start, end uint64
   _, err := fmt.Sscanf(string(r), "%v-%v", &start, &end)
   if err != nil {
      return 0, 0, err
   }
   return start, end, nil
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#addressing-simple-to-explicit
type mpd struct {
   MediaPresentationDuration string `xml:"mediaPresentationDuration,attr"`
   Period []period
}

func (m mpd) Seconds() (float64, error) {
   s := strings.TrimPrefix(m.MediaPresentationDuration, "PT")
   duration, err := time.ParseDuration(strings.ToLower(s))
   if err != nil {
      return 0, err
   }
   return duration.Seconds(), nil
}

type period struct {
   AdaptationSet []adaptation_set
   mpd *mpd
}

type adaptation_set struct {
   Codecs *string `xml:"codecs,attr"`
   Lang *string `xml:"lang,attr"`
   MimeType *string `xml:"mimeType,attr"`
   Representation []Representation
   Role *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
   period *period
}

func (s SegmentTemplate) GetInitialization(id string) (string, bool) {
   if v := s.Initialization; v != nil {
      return strings.Replace(*v, "$RepresentationID$", id, 1), true
   }
   return "", false
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

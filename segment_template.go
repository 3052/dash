package dash

import (
   "errors"
   "fmt"
   "math"
   "net/url"
   "strings"
)

// SegmentTemplate defines specific rules for generating segment URLs.
type SegmentTemplate struct {
   Duration               uint             `xml:"duration,attr"`
   EndNumber              uint             `xml:"endNumber,attr"`
   Initialization         string           `xml:"initialization,attr"`
   Media                  string           `xml:"media,attr"`
   PresentationTimeOffset uint             `xml:"presentationTimeOffset,attr"`
   StartNumber            *uint            `xml:"startNumber,attr"`
   Timescale              *uint            `xml:"timescale,attr"`
   SegmentTimeline        *SegmentTimeline `xml:"SegmentTimeline"`

   // Navigation
   ParentAdaptationSet  *AdaptationSet  `xml:"-"`
   ParentRepresentation *Representation `xml:"-"`
}

func (st *SegmentTemplate) GetStartNumber() uint {
   if st.StartNumber != nil {
      return *st.StartNumber
   }
   return 1
}

func (st *SegmentTemplate) GetTimescale() uint {
   if st.Timescale != nil {
      return *st.Timescale
   }
   return 1
}

func (st *SegmentTemplate) GetNumberRange() []uint {
   start := st.GetStartNumber()
   end := st.EndNumber
   if end < start {
      return nil
   }
   size := end - start + 1
   nums := make([]uint, size)
   for idx := uint(0); idx < size; idx++ {
      nums[idx] = start + idx
   }
   return nums
}

func (st *SegmentTemplate) GetTimelineNumbers() []uint {
   if st.SegmentTimeline == nil {
      return nil
   }
   var numbers []uint
   current := st.GetStartNumber()
   for _, segment := range st.SegmentTimeline.S {
      count := 1
      if segment.R > 0 {
         count += segment.R
      }
      for idx := 0; idx < count; idx++ {
         numbers = append(numbers, current)
         current++
      }
   }
   return numbers
}

func (st *SegmentTemplate) GetTimelineTimes() []uint {
   if st.SegmentTimeline == nil {
      return nil
   }
   var times []uint
   currentTime := st.PresentationTimeOffset
   for _, segment := range st.SegmentTimeline.S {
      count := 1
      if segment.R > 0 {
         count += segment.R
      }
      for idx := 0; idx < count; idx++ {
         times = append(times, currentTime)
         currentTime += segment.D
      }
   }
   return times
}

func (st *SegmentTemplate) GetDurationBasedNumbers() ([]uint, error) {
   var manifestPeriod *Period
   if st.ParentRepresentation != nil && st.ParentRepresentation.Parent != nil {
      manifestPeriod = st.ParentRepresentation.Parent.Parent
   } else if st.ParentAdaptationSet != nil {
      manifestPeriod = st.ParentAdaptationSet.Parent
   }

   if manifestPeriod == nil {
      return nil, errors.New("SegmentTemplate is not properly linked to a Period")
   }

   pDuration, err := manifestPeriod.GetDuration()
   if err != nil {
      return nil, err
   }
   if pDuration == 0 {
      return nil, errors.New("period duration is zero or missing")
   }

   if st.Duration == 0 {
      return nil, errors.New("SegmentTemplate duration is zero")
   }

   segDurSec := float64(st.Duration) / float64(st.GetTimescale())
   count := uint(math.Ceil(pDuration.Seconds() / segDurSec))

   start := st.GetStartNumber()
   numbers := make([]uint, count)
   for idx := uint(0); idx < count; idx++ {
      numbers[idx] = start + idx
   }
   return numbers, nil
}

func (st *SegmentTemplate) GetSegmentUrls(rep *Representation) ([]*url.URL, error) {
   if st.Media == "" {
      return nil, nil
   }

   if strings.Contains(st.Media, "$Time$") {
      times := st.GetTimelineTimes()
      if len(times) == 0 {
         return nil, errors.New("media template requires $Time$ but no SegmentTimeline found")
      }
      var urls []*url.URL
      for _, timeVal := range times {
         parsedUrl, err := st.ResolveMediaTime(rep, int(timeVal))
         if err != nil {
            return nil, err
         }
         urls = append(urls, parsedUrl)
      }
      return urls, nil
   }

   var numbers []uint
   if st.SegmentTimeline != nil {
      numbers = st.GetTimelineNumbers()
   }
   if len(numbers) == 0 && st.EndNumber > 0 {
      numbers = st.GetNumberRange()
   }
   if len(numbers) == 0 && st.Duration > 0 {
      var err error
      numbers, err = st.GetDurationBasedNumbers()
      if err != nil {
         return nil, err
      }
   }

   var urls []*url.URL
   for _, num := range numbers {
      parsedUrl, err := st.ResolveMedia(rep, int(num))
      if err != nil {
         return nil, err
      }
      urls = append(urls, parsedUrl)
   }
   return urls, nil
}

func (st *SegmentTemplate) ResolveInitialization(rep *Representation) (*url.URL, error) {
   base, initStr, err := st.prepareTemplateString(rep, st.Initialization)
   if err != nil {
      return nil, err
   }
   return resolveRef(base, initStr)
}

func (st *SegmentTemplate) ResolveMedia(rep *Representation, number int) (*url.URL, error) {
   base, mediaStr, err := st.prepareTemplateString(rep, st.Media)
   if err != nil {
      return nil, err
   }

   formats := []string{"%02d", "%03d", "%04d", "%05d", "%06d", "%07d", "%08d", "%09d"}
   for _, format := range formats {
      token := fmt.Sprintf("$Number%s$", format)
      if strings.Contains(mediaStr, token) {
         mediaStr = strings.ReplaceAll(mediaStr, token, fmt.Sprintf(format, number))
      }
   }

   mediaStr = strings.ReplaceAll(mediaStr, "$Number$", fmt.Sprintf("%d", number))
   return resolveRef(base, mediaStr)
}

func (st *SegmentTemplate) ResolveMediaTime(rep *Representation, timeVal int) (*url.URL, error) {
   base, mediaStr, err := st.prepareTemplateString(rep, st.Media)
   if err != nil {
      return nil, err
   }
   mediaStr = strings.ReplaceAll(mediaStr, "$Time$", fmt.Sprintf("%d", timeVal))
   return resolveRef(base, mediaStr)
}

func (st *SegmentTemplate) prepareTemplateString(rep *Representation, templateStr string) (*url.URL, string, error) {
   base, err := st.getParentBaseUrl()
   if err != nil {
      return nil, "", err
   }

   var repId string
   if rep != nil {
      repId = rep.Id
   } else if st.ParentRepresentation != nil {
      repId = st.ParentRepresentation.Id
   }

   if repId != "" {
      templateStr = strings.ReplaceAll(templateStr, "$RepresentationID$", repId)
   }

   return base, templateStr, nil
}

func (st *SegmentTemplate) getParentBaseUrl() (*url.URL, error) {
   if st.ParentRepresentation != nil {
      return st.ParentRepresentation.ResolveBaseUrl()
   }
   if st.ParentAdaptationSet != nil {
      return st.ParentAdaptationSet.getAbsoluteBaseUrl()
   }
   return nil, errors.New("SegmentTemplate has no parent linked")
}

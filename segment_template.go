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
   for i := uint(0); i < size; i++ {
      nums[i] = start + i
   }
   return nums
}

func (st *SegmentTemplate) GetTimelineNumbers() []uint {
   if st.SegmentTimeline == nil {
      return nil
   }
   var numbers []uint
   current := st.GetStartNumber()
   for _, s := range st.SegmentTimeline.S {
      count := 1
      if s.R > 0 {
         count += s.R
      }
      for i := 0; i < count; i++ {
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
   for _, s := range st.SegmentTimeline.S {
      count := 1
      if s.R > 0 {
         count += s.R
      }
      for i := 0; i < count; i++ {
         times = append(times, currentTime)
         currentTime += s.D
      }
   }
   return times
}

func (st *SegmentTemplate) GetDurationBasedNumbers() ([]uint, error) {
   var period *Period
   if st.ParentRepresentation != nil && st.ParentRepresentation.Parent != nil {
      period = st.ParentRepresentation.Parent.Parent
   } else if st.ParentAdaptationSet != nil {
      period = st.ParentAdaptationSet.Parent
   }

   if period == nil {
      return nil, errors.New("SegmentTemplate is not properly linked to a Period")
   }

   pDur, err := period.GetDuration()
   if err != nil {
      return nil, err
   }
   if pDur == 0 {
      return nil, errors.New("period duration is zero or missing")
   }

   if st.Duration == 0 {
      return nil, errors.New("SegmentTemplate duration is zero")
   }

   segDurSec := float64(st.Duration) / float64(st.GetTimescale())
   count := uint(math.Ceil(pDur.Seconds() / segDurSec))

   start := st.GetStartNumber()
   numbers := make([]uint, count)
   for i := uint(0); i < count; i++ {
      numbers[i] = start + i
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
      for _, t := range times {
         u, err := st.ResolveMediaTime(rep, int(t))
         if err != nil {
            return nil, err
         }
         urls = append(urls, u)
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
   for _, n := range numbers {
      u, err := st.ResolveMedia(rep, int(n))
      if err != nil {
         return nil, err
      }
      urls = append(urls, u)
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
   for _, f := range formats {
      token := fmt.Sprintf("$Number%s$", f)
      if strings.Contains(mediaStr, token) {
         mediaStr = strings.ReplaceAll(mediaStr, token, fmt.Sprintf(f, number))
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

func (st *SegmentTemplate) prepareTemplateString(rep *Representation, tmpl string) (*url.URL, string, error) {
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
      tmpl = strings.ReplaceAll(tmpl, "$RepresentationID$", repId)
   }

   return base, tmpl, nil
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

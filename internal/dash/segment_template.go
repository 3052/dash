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
   Duration               uint   `xml:"duration,attr,omitempty"`
   EndNumber              uint   `xml:"endNumber,attr,omitempty"`
   Initialization         string `xml:"initialization,attr,omitempty"`
   Media                  string `xml:"media,attr,omitempty"`
   PresentationTimeOffset uint   `xml:"presentationTimeOffset,attr,omitempty"`
   // StartNumber is a pointer to distinguish between missing (nil) and explicitly 0.
   StartNumber *uint `xml:"startNumber,attr,omitempty"`
   // Timescale is a pointer to distinguish between missing (nil) and explicitly 0.
   Timescale       *uint            `xml:"timescale,attr,omitempty"`
   SegmentTimeline *SegmentTimeline `xml:"SegmentTimeline"`

   // Navigation
   ParentAdaptationSet  *AdaptationSet  `xml:"-"`
   ParentRepresentation *Representation `xml:"-"`
}

// GetStartNumber returns the StartNumber if present, otherwise returns default 1.
func (st *SegmentTemplate) GetStartNumber() uint {
   if st.StartNumber != nil {
      return *st.StartNumber
   }
   return 1
}

// GetTimescale returns the Timescale if present, otherwise returns default 1.
func (st *SegmentTemplate) GetTimescale() uint {
   if st.Timescale != nil {
      return *st.Timescale
   }
   return 1
}

// GetNumberRange returns a slice of numbers from StartNumber to EndNumber (inclusive).
// If EndNumber is less than StartNumber, it returns nil.
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

// GetTimelineNumbers returns a slice of segment numbers derived from the SegmentTimeline.
// It starts at StartNumber and iterates through the S elements, accounting for repetitions (@r).
func (st *SegmentTemplate) GetTimelineNumbers() []uint {
   if st.SegmentTimeline == nil {
      return nil
   }

   var numbers []uint
   current := st.GetStartNumber()

   for _, s := range st.SegmentTimeline.S {
      // Each 'S' element represents a segment duration, repeated 'r' times.
      // Total segments described by this S element = 1 + r.
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

// GetTimelineTimes returns a slice of start times for segments derived from the SegmentTimeline.
// It starts at PresentationTimeOffset and accumulates durations ('d').
func (st *SegmentTemplate) GetTimelineTimes() []uint {
   if st.SegmentTimeline == nil {
      return nil
   }

   var times []uint
   // Req 39: Initialize with PresentationTimeOffset
   var currentTime uint = st.PresentationTimeOffset

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

// GetDurationBasedNumbers calculates segment numbers based on Period duration and fixed Segment duration.
// Formula: Ceil( AsSeconds(Period@duration) / (SegmentTemplate@duration / SegmentTemplate@timescale) )
func (st *SegmentTemplate) GetDurationBasedNumbers() ([]uint, error) {
   // 1. Find the parent Period
   var period *Period
   if st.ParentRepresentation != nil && st.ParentRepresentation.Parent != nil {
      period = st.ParentRepresentation.Parent.Parent
   } else if st.ParentAdaptationSet != nil {
      period = st.ParentAdaptationSet.Parent
   }

   if period == nil {
      return nil, errors.New("SegmentTemplate is not properly linked to a Period")
   }

   // 2. Get Period Duration (Seconds)
   pDur, err := period.GetDuration()
   if err != nil {
      return nil, err
   }
   if pDur == 0 {
      return nil, errors.New("period duration is zero or missing")
   }
   pDurSec := pDur.Seconds()

   // 3. Get Segment Duration parameters
   if st.Duration == 0 {
      return nil, errors.New("SegmentTemplate duration is zero")
   }

   timescale := float64(st.GetTimescale())
   segDurSec := float64(st.Duration) / timescale

   // 4. Calculate count
   countFloat := math.Ceil(pDurSec / segDurSec)
   count := uint(countFloat)

   // 5. Generate numbers
   start := st.GetStartNumber()
   numbers := make([]uint, count)
   for i := uint(0); i < count; i++ {
      numbers[i] = start + i
   }

   return numbers, nil
}

// GetSegmentURLs returns all segment URLs defined by this template.
// It determines the generation strategy (Time-based vs Number-based)
// and the source of the segment list (Timeline, Range, or Duration).
func (st *SegmentTemplate) GetSegmentURLs(rep *Representation) ([]*url.URL, error) {
   if st.Media == "" {
      return nil, nil
   }

   // Strategy 1: Time-based addressing ($Time$ present in template)
   if strings.Contains(st.Media, "$Time$") {
      times := st.GetTimelineTimes()
      if len(times) == 0 {
         // If $Time$ is required but no timeline exists, we cannot generate URLs reliably.
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

   // Strategy 2: Number-based addressing
   var numbers []uint

   // Source A: SegmentTimeline
   if st.SegmentTimeline != nil {
      numbers = st.GetTimelineNumbers()
   }

   // Source B: explicit EndNumber range
   if len(numbers) == 0 && st.EndNumber > 0 {
      numbers = st.GetNumberRange()
   }

   // Source C: Duration-based calculation
   if len(numbers) == 0 && st.Duration > 0 {
      var err error
      numbers, err = st.GetDurationBasedNumbers()
      if err != nil {
         // If calculation fails (e.g. unlinked parent), return error
         return nil, err
      }
   }

   if len(numbers) == 0 {
      return nil, nil
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

// ResolveInitialization resolves the @initialization attribute against the parent BaseURL.
// It replaces the literal "$RepresentationID$" with the ID of the provided Representation.
func (st *SegmentTemplate) ResolveInitialization(rep *Representation) (*url.URL, error) {
   base, err := st.getParentBaseURL()
   if err != nil {
      return nil, err
   }

   initStr := st.Initialization

   // Determine the ID to use for replacement
   var repID string
   if rep != nil {
      repID = rep.ID
   } else if st.ParentRepresentation != nil {
      repID = st.ParentRepresentation.ID
   }

   // Perform replacement if an ID was found
   if repID != "" {
      initStr = strings.ReplaceAll(initStr, "$RepresentationID$", repID)
   }

   return resolveRef(base, initStr)
}

// ResolveMedia resolves the @media attribute for number-based addressing.
// It performs substitutions for $RepresentationID$ and $Number$ (including format variants).
func (st *SegmentTemplate) ResolveMedia(rep *Representation, number int) (*url.URL, error) {
   base, err := st.getParentBaseURL()
   if err != nil {
      return nil, err
   }

   mediaStr := st.Media

   // 1. Replace $RepresentationID$
   var repID string
   if rep != nil {
      repID = rep.ID
   } else if st.ParentRepresentation != nil {
      repID = st.ParentRepresentation.ID
   }
   if repID != "" {
      mediaStr = strings.ReplaceAll(mediaStr, "$RepresentationID$", repID)
   }

   // 2. Replace $Number%0xd$ variants
   formats := []string{
      "%02d", "%03d", "%04d", "%05d",
      "%06d", "%07d", "%08d", "%09d",
   }

   for _, f := range formats {
      token := fmt.Sprintf("$Number%s$", f)
      if strings.Contains(mediaStr, token) {
         replacement := fmt.Sprintf(f, number)
         mediaStr = strings.ReplaceAll(mediaStr, token, replacement)
      }
   }

   // 3. Replace bare $Number$
   mediaStr = strings.ReplaceAll(mediaStr, "$Number$", fmt.Sprintf("%d", number))

   return resolveRef(base, mediaStr)
}

// ResolveMediaTime resolves the @media attribute for time-based addressing.
// It performs substitutions for $RepresentationID$ and $Time$.
func (st *SegmentTemplate) ResolveMediaTime(rep *Representation, timeVal int) (*url.URL, error) {
   base, err := st.getParentBaseURL()
   if err != nil {
      return nil, err
   }

   mediaStr := st.Media

   // 1. Replace $RepresentationID$
   var repID string
   if rep != nil {
      repID = rep.ID
   } else if st.ParentRepresentation != nil {
      repID = st.ParentRepresentation.ID
   }
   if repID != "" {
      mediaStr = strings.ReplaceAll(mediaStr, "$RepresentationID$", repID)
   }

   // 2. Replace $Time$
   mediaStr = strings.ReplaceAll(mediaStr, "$Time$", fmt.Sprintf("%d", timeVal))

   return resolveRef(base, mediaStr)
}

func (st *SegmentTemplate) getParentBaseURL() (*url.URL, error) {
   if st.ParentRepresentation != nil {
      return st.ParentRepresentation.ResolveBaseURL()
   }
   if st.ParentAdaptationSet != nil {
      return st.ParentAdaptationSet.getAbsoluteBaseURL()
   }
   return nil, errors.New("SegmentTemplate has no parent linked")
}

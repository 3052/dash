package dash

import (
   "errors"
   "fmt"
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
   StartNumber     *uint            `xml:"startNumber,attr,omitempty"`
   Timescale       uint             `xml:"timescale,attr,omitempty"`
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

// ResolveMedia resolves the @media attribute against the parent BaseURL.
// It performs substitutions for $RepresentationID$, $Time$, and $Number$ (including format variants).
func (st *SegmentTemplate) ResolveMedia(rep *Representation, number, timeVal int) (*url.URL, error) {
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

   // 3. Replace $Number%0xd$ variants
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

   // 4. Replace bare $Number$
   mediaStr = strings.ReplaceAll(mediaStr, "$Number$", fmt.Sprintf("%d", number))

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

package dash

import (
   "errors"
   "fmt"
   "net/url"
   "strings"
)

// SegmentTemplate defines specific rules for generating segment URLs.
type SegmentTemplate struct {
   Duration               uint             `xml:"duration,attr,omitempty"`
   EndNumber              uint             `xml:"endNumber,attr,omitempty"`
   Initialization         string           `xml:"initialization,attr,omitempty"`
   Media                  string           `xml:"media,attr,omitempty"`
   PresentationTimeOffset uint             `xml:"presentationTimeOffset,attr,omitempty"`
   StartNumber            uint             `xml:"startNumber,attr,omitempty"`
   Timescale              uint             `xml:"timescale,attr,omitempty"`
   SegmentTimeline        *SegmentTimeline `xml:"SegmentTimeline"`

   // Navigation
   ParentAdaptationSet  *AdaptationSet  `xml:"-"`
   ParentRepresentation *Representation `xml:"-"`
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
   // We iterate through the specific formats requested
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

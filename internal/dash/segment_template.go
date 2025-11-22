package dash

import (
   "errors"
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
// If rep is nil, it attempts to use the SegmentTemplate's direct parent Representation.
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
func (st *SegmentTemplate) ResolveMedia() (*url.URL, error) {
   base, err := st.getParentBaseURL()
   if err != nil {
      return nil, err
   }
   return resolveRef(base, st.Media)
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

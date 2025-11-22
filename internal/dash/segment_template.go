package dash

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
   // Note: SegmentTemplate can belong to AdaptationSet OR Representation
   ParentAdaptationSet  *AdaptationSet  `xml:"-"`
   ParentRepresentation *Representation `xml:"-"`
}

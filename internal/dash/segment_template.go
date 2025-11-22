package dash

// SegmentTemplate specifies a mechanism to construct segment URLs.
type SegmentTemplate struct {
   Duration               uint64           `xml:"duration,attr,omitempty"`
   EndNumber              uint64           `xml:"endNumber,attr,omitempty"`
   Initialization         string           `xml:"initialization,attr,omitempty"`
   Media                  string           `xml:"media,attr,omitempty"`
   PresentationTimeOffset uint64           `xml:"presentationTimeOffset,attr,omitempty"`
   StartNumber            uint64           `xml:"startNumber,attr,omitempty"`
   Timescale              uint64           `xml:"timescale,attr,omitempty"`
   SegmentTimeline        *SegmentTimeline `xml:"SegmentTimeline"`
}

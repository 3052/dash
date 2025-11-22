package dash

// SegmentTemplate represents the SegmentTemplate element.
type SegmentTemplate struct {
   Timescale              uint64           `xml:"timescale,attr,omitempty"`
   Duration               uint64           `xml:"duration,attr,omitempty"`
   StartNumber            uint64           `xml:"startNumber,attr,omitempty"`
   EndNumber              uint64           `xml:"endNumber,attr,omitempty"`
   PresentationTimeOffset uint64           `xml:"presentationTimeOffset,attr,omitempty"`
   Initialization         string           `xml:"initialization,attr,omitempty"`
   Media                  string           `xml:"media,attr,omitempty"`
   SegmentTimeline        *SegmentTimeline `xml:"SegmentTimeline,omitempty"`
}

// SegmentTimeline represents the SegmentTimeline element.
type SegmentTimeline struct {
   S []SegmentTimelineS `xml:"S"`
}

// SegmentTimelineS represents the S element within a timeline.
type SegmentTimelineS struct {
   D uint64 `xml:"d,attr"`           // Duration
   R int    `xml:"r,attr,omitempty"` // Repeat count
}

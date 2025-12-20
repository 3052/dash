package dash

// SegmentTimeline defines specific timing for segments.
type SegmentTimeline struct {
   S []*S `xml:"S"`
}

// S represents a segment within the timeline.
type S struct {
   Duration uint `xml:"d,attr"` // Duration
   Repeat   int  `xml:"r,attr"` // Repeat count
}

package dash

// SegmentTimeline defines specific timing for segments.
type SegmentTimeline struct {
   S []*S `xml:"S"`
}

// S represents a segment within the timeline.
type S struct {
   D uint `xml:"d,attr"`           // Duration
   R int  `xml:"r,attr,omitempty"` // Repeat count
}

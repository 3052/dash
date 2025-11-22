package dash

// S represents a segment within the SegmentTimeline.
type S struct {
   D uint64 `xml:"d,attr"`
   R int64  `xml:"r,attr,omitempty"`
}

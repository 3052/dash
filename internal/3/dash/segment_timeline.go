package dash

import "encoding/xml"

// SegmentTimeline contains a list of Segment timings (S elements).
type SegmentTimeline struct {
   XMLName xml.Name `xml:"SegmentTimeline"`
   S       []S      `xml:"S"`
}

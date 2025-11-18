package dash

import "encoding/xml"

// S represents the S (segment) element within a SegmentTimeline.
type S struct {
   XMLName xml.Name `xml:"S"`
   // T is the presentation time of the first segment in the series.
   T uint64 `xml:"t,attr,omitempty"`
   // D is the duration of the segment in timescale units.
   D uint64 `xml:"d,attr"`
   // R is the repeat count. A value of 2 means this segment is repeated 2 more times.
   R int `xml:"r,attr,omitempty"`
}

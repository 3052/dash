package dash

import "encoding/xml"

// S represents a Segment in the timeline.
type S struct {
   XMLName xml.Name `xml:"S"`
   T       uint64   `xml:"t,attr,omitempty"`
   D       uint64   `xml:"d,attr"`
   R       int      `xml:"r,attr,omitempty"`
}

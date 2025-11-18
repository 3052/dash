package dash

import "encoding/xml"

// SegmentList contains a list of SegmentURLs.
type SegmentList struct {
   XMLName        xml.Name        `xml:"SegmentList"`
   Timescale      uint32          `xml:"timescale,attr,omitempty"`
   Duration       uint32          `xml:"duration,attr,omitempty"`
   Initialization *Initialization `xml:"Initialization,omitempty"`
   SegmentURLs    []SegmentURL    `xml:"SegmentURL,omitempty"`
}

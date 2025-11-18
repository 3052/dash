package dash

import "encoding/xml"

// SegmentList represents the SegmentList element, providing an explicit list of media segments.
type SegmentList struct {
   XMLName        xml.Name        `xml:"SegmentList"`
   Timescale      int             `xml:"timescale,attr,omitempty"`
   Duration       int             `xml:"duration,attr,omitempty"`
   Initialization *Initialization `xml:"Initialization,omitempty"`
   SegmentURLs    []*SegmentURL   `xml:"SegmentURL"`
}

package dash

import "encoding/xml"

// SegmentBase represents a SegmentBase element in the MPD.
type SegmentBase struct {
   XMLName        xml.Name        `xml:"SegmentBase"`
   IndexRange     *string         `xml:"indexRange,attr,omitempty"`
   Initialization *Initialization `xml:"Initialization,omitempty"`
}

// Initialization represents an Initialization element in the MPD.
type Initialization struct {
   XMLName xml.Name `xml:"Initialization"`
   Range   *string  `xml:"range,attr,omitempty"`
}

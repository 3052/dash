package dash

import "encoding/xml"

// SegmentURL explicitly lists a segment media URL/range.
type SegmentURL struct {
   XMLName    xml.Name `xml:"SegmentURL"`
   Media      string   `xml:"media,attr,omitempty"`
   MediaRange string   `xml:"mediaRange,attr,omitempty"`
}

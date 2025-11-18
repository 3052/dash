package dash

import "encoding/xml"

// SegmentURL represents the SegmentURL element, which provides the URL for a media segment.
type SegmentURL struct {
   XMLName xml.Name `xml:"SegmentURL"`
   Media   string   `xml:"media,attr"`
}

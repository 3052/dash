package dash

import "encoding/xml"

// SegmentTemplate represents the SegmentTemplate element.
type SegmentTemplate struct {
   XMLName         xml.Name         `xml:"SegmentTemplate"`
   Timescale       int              `xml:"timescale,attr"`
   Duration        uint             `xml:"duration,attr,omitempty"`
   Media           string           `xml:"media,attr"`
   Initialization  string           `xml:"initialization,attr"`
   StartNumber     uint             `xml:"startNumber,attr,omitempty"`
   EndNumber       uint             `xml:"endNumber,attr,omitempty"`
   SegmentTimeline *SegmentTimeline `xml:"SegmentTimeline,omitempty"`
}

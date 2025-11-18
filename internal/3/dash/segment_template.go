package dash

import "encoding/xml"

// SegmentTemplate allows generating segment URLs using a template.
type SegmentTemplate struct {
   XMLName                xml.Name         `xml:"SegmentTemplate"`
   Timescale              uint32           `xml:"timescale,attr,omitempty"`
   Media                  string           `xml:"media,attr,omitempty"`
   Initialization         string           `xml:"initialization,attr,omitempty"`
   StartNumber            uint32           `xml:"startNumber,attr,omitempty"`
   PresentationTimeOffset uint64           `xml:"presentationTimeOffset,attr,omitempty"`
   SegmentTimeline        *SegmentTimeline `xml:"SegmentTimeline,omitempty"`
}

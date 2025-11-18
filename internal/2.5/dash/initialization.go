package dash

import "encoding/xml"

// Initialization represents the Initialization element, typically found
// within a SegmentList, which points to the initialization segment.
type Initialization struct {
   XMLName   xml.Name `xml:"Initialization"`
   SourceURL string   `xml:"sourceURL,attr"`
}

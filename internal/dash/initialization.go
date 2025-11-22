package dash

import "encoding/xml"

// Initialization URL info found in SegmentList.
type Initialization struct {
   XMLName   xml.Name `xml:"Initialization"`
   SourceURL string   `xml:"sourceURL,attr,omitempty"`
   Range     string   `xml:"range,attr,omitempty"`
}

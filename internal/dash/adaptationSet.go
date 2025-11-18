package dash

import "encoding/xml"

// AdaptationSet represents the AdaptationSet element.
type AdaptationSet struct {
   XMLName            xml.Name             `xml:"AdaptationSet"`
   ContentType        string               `xml:"contentType,attr,omitempty"`
   Lang               string               `xml:"lang,attr,omitempty"`
   MimeType           string               `xml:"mimeType,attr,omitempty"`
   SegmentAlignment   bool                 `xml:"segmentAlignment,attr,omitempty"`
   StartWithSAP       int                  `xml:"startWithSAP,attr,omitempty"`
   Representations    []*Representation    `xml:"Representation"`
   SegmentTemplate    *SegmentTemplate     `xml:"SegmentTemplate"`
   ContentProtections []*ContentProtection `xml:"ContentProtection"`
}

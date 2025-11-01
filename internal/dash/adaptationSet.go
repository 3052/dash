package dash

import "encoding/xml"

// AdaptationSet represents an AdaptationSet element in the MPD.
type AdaptationSet struct {
   XMLName            xml.Name             `xml:"AdaptationSet"`
   MimeType           *string              `xml:"mimeType,attr,omitempty"`
   Lang               *string              `xml:"lang,attr,omitempty"`
   MaxWidth           *uint64              `xml:"maxWidth,attr,omitempty"`
   MaxHeight          *uint64              `xml:"maxHeight,attr,omitempty"`
   MinWidth           *uint64              `xml:"minWidth,attr,omitempty"`
   MinHeight          *uint64              `xml:"minHeight,attr,omitempty"`
   PAR                *string              `xml:"par,attr,omitempty"`
   SAR                *string              `xml:"sar,attr,omitempty"`
   SegmentAlignment   *bool                `xml:"segmentAlignment,attr,omitempty"`
   StartWithSAP       *uint64              `xml:"startWithSAP,attr,omitempty"`
   ContentProtections []*ContentProtection `xml:"ContentProtection,omitempty"`
   Representations    []*Representation    `xml:"Representation,omitempty"`
}

package dash

// AdaptationSet represents the AdaptationSet element.
type AdaptationSet struct {
   MimeType           string              `xml:"mimeType,attr,omitempty"`
   Codecs             string              `xml:"codecs,attr,omitempty"`
   Lang               string              `xml:"lang,attr,omitempty"`
   Width              int                 `xml:"width,attr,omitempty"`
   Height             int                 `xml:"height,attr,omitempty"`
   ContentProtections []ContentProtection `xml:"ContentProtection,omitempty"`
   Roles              []Role              `xml:"Role,omitempty"`
   SegmentTemplate    *SegmentTemplate    `xml:"SegmentTemplate,omitempty"`
   Representations    []Representation    `xml:"Representation,omitempty"`
}

package dash

// Representation represents the Representation element.
type Representation struct {
   ID                 string              `xml:"id,attr,omitempty"`
   Bandwidth          int                 `xml:"bandwidth,attr,omitempty"`
   MimeType           string              `xml:"mimeType,attr,omitempty"`
   Codecs             string              `xml:"codecs,attr,omitempty"`
   Width              int                 `xml:"width,attr,omitempty"`
   Height             int                 `xml:"height,attr,omitempty"`
   BaseURL            string              `xml:"BaseURL,omitempty"`
   SegmentTemplate    *SegmentTemplate    `xml:"SegmentTemplate,omitempty"`
   ContentProtections []ContentProtection `xml:"ContentProtection,omitempty"`
   SegmentBase        *SegmentBase        `xml:"SegmentBase,omitempty"`
   SegmentList        *SegmentList        `xml:"SegmentList,omitempty"`
}

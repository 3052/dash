package dash

import "encoding/xml"

// Representation describes a version of the content (e.g., a specific bitrate).
type Representation struct {
   XMLName            xml.Name            `xml:"Representation"`
   ID                 string              `xml:"id,attr"`
   Bandwidth          int                 `xml:"bandwidth,attr"`
   Width              int                 `xml:"width,attr,omitempty"`
   Height             int                 `xml:"height,attr,omitempty"`
   FrameRate          string              `xml:"frameRate,attr,omitempty"`
   Codecs             string              `xml:"codecs,attr,omitempty"`
   MimeType           string              `xml:"mimeType,attr,omitempty"`
   BaseURL            string              `xml:"BaseURL,omitempty"`
   ContentProtections []ContentProtection `xml:"ContentProtection,omitempty"`
   SegmentTemplate    *SegmentTemplate    `xml:"SegmentTemplate,omitempty"`
   SegmentList        *SegmentList        `xml:"SegmentList,omitempty"`
}

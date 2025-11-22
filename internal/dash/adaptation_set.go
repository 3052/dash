package dash

import "encoding/xml"

// AdaptationSet contains a set of Representations.
type AdaptationSet struct {
   XMLName             xml.Name            `xml:"AdaptationSet"`
   ID                  string              `xml:"id,attr,omitempty"`
   MimeType            string              `xml:"mimeType,attr,omitempty"`
   Codecs              string              `xml:"codecs,attr,omitempty"`
   Width               int                 `xml:"width,attr,omitempty"`
   Height              int                 `xml:"height,attr,omitempty"`
   FrameRate           string              `xml:"frameRate,attr,omitempty"`
   Lang                string              `xml:"lang,attr,omitempty"`
   SegmentAlignment    string              `xml:"segmentAlignment,attr,omitempty"`
   SubsegmentAlignment string              `xml:"subsegmentAlignment,attr,omitempty"`
   BaseURL             string              `xml:"BaseURL,omitempty"`
   ContentProtections  []ContentProtection `xml:"ContentProtection,omitempty"`
   Representations     []Representation    `xml:"Representation,omitempty"`
   SegmentTemplate     *SegmentTemplate    `xml:"SegmentTemplate,omitempty"`
}

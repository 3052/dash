package dash

// AdaptationSet represents a set of interchangeable encoded versions of one or several media content components.
type AdaptationSet struct {
   Codecs            string               `xml:"codecs,attr,omitempty"`
   Height            int                  `xml:"height,attr,omitempty"`
   Lang              string               `xml:"lang,attr,omitempty"`
   MimeType          string               `xml:"mimeType,attr,omitempty"`
   Width             int                  `xml:"width,attr,omitempty"`
   ContentProtection []*ContentProtection `xml:"ContentProtection"`
   Role              []*Role              `xml:"Role"`
   SegmentTemplate   *SegmentTemplate     `xml:"SegmentTemplate"`
   Representation    []*Representation    `xml:"Representation"`
}

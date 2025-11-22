package dash

// SegmentBase represents the SegmentBase element.
type SegmentBase struct {
   IndexRange     string                 `xml:"indexRange,attr,omitempty"`
   Initialization *SegmentInitialization `xml:"Initialization,omitempty"`
}

// SegmentInitialization represents the Initialization element.
type SegmentInitialization struct {
   Range     string `xml:"range,attr,omitempty"`
   SourceURL string `xml:"sourceURL,attr,omitempty"` // Used by SegmentList
}

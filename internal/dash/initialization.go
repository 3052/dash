package dash

// Initialization contains URL and byte range information for initialization segments.
type Initialization struct {
   // Used in SegmentBase
   Range string `xml:"range,attr,omitempty"`
   // Used in SegmentList
   SourceURL string `xml:"sourceURL,attr,omitempty"`

   // Navigation
   Parent *SegmentList `xml:"-"`
}

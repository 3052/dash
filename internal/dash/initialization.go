package dash

// Initialization specifies the initialization segment info.
// Note: This type merges attributes used in both SegmentBase (range) and SegmentList (sourceURL).
type Initialization struct {
   Range     string `xml:"range,attr,omitempty"`
   SourceURL string `xml:"sourceURL,attr,omitempty"`
}

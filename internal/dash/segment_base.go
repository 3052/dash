package dash

// SegmentBase specifies the base information for a segment.
type SegmentBase struct {
   IndexRange     string          `xml:"indexRange,attr,omitempty"`
   Initialization *Initialization `xml:"Initialization"`
}

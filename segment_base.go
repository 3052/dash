package dash

// SegmentBase defines base information for segments.
type SegmentBase struct {
   IndexRange     string          `xml:"indexRange,attr"`
   Initialization *Initialization `xml:"Initialization"`
}

func (sb *SegmentBase) link() {
}

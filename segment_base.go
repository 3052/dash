package dash

// SegmentBase defines base information for segments.
type SegmentBase struct {
   IndexRange     string          `xml:"indexRange,attr"`
   Initialization *Initialization `xml:"Initialization"`
}

func (sb *SegmentBase) link() {
   // Note: Req 10.2 specifies Initialization -> SegmentList
   // Logic for Initialization inside SegmentBase is not strictly required by Req 10,
   // but structure is provided here for completeness of parsing.
}

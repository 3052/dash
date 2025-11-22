package dash

// SegmentList defines a list of segments.
type SegmentList struct {
   Initialization *Initialization `xml:"Initialization"`
   SegmentURL     []*SegmentURL   `xml:"SegmentURL"`
}

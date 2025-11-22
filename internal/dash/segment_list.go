package dash

// SegmentList represents the SegmentList element.
type SegmentList struct {
   Initialization *SegmentInitialization `xml:"Initialization,omitempty"`
   SegmentURLs    []SegmentURL           `xml:"SegmentURL,omitempty"`
}

// SegmentURL represents the SegmentURL element.
type SegmentURL struct {
   Media string `xml:"media,attr,omitempty"`
}

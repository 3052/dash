package dash

// SegmentList contains a list of SegmentURLs.
type SegmentList struct {
   Initialization *Initialization `xml:"Initialization"`
   SegmentURLs    []*SegmentURL   `xml:"SegmentURL"`

   // Navigation
   Parent *Representation `xml:"-"`
}

func (sl *SegmentList) link() {
   if sl.Initialization != nil {
      // Req 10.2: Initialization to SegmentList
      sl.Initialization.Parent = sl
   }
   for _, u := range sl.SegmentURLs {
      // Req 10.8: SegmentURL to SegmentList
      u.Parent = sl
   }
}

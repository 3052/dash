package dash

import "net/url"

// SegmentList contains a list of SegmentURLs.
type SegmentList struct {
   Duration       uint            `xml:"duration,attr"`
   Timescale      *uint           `xml:"timescale,attr"`
   Initialization *Initialization `xml:"Initialization"`
   SegmentURLs    []*SegmentURL   `xml:"SegmentURL"`
   // Navigation
   Parent *Representation `xml:"-"`
}

// GetTimescale returns the Timescale if present, otherwise returns default 1.
func (sl *SegmentList) GetTimescale() uint {
   if sl.Timescale != nil {
      return *sl.Timescale
   }
   return 1
}

// getParentBaseURL retrieves the resolved BaseURL from the parent Representation.
func (sl *SegmentList) getParentBaseURL() (*url.URL, error) {
   return sl.Parent.ResolveBaseURL()
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

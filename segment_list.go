package dash

import "net/url"

// SegmentList contains a list of SegmentUrls.
type SegmentList struct {
   Duration       uint            `xml:"duration,attr"`
   Timescale      *uint           `xml:"timescale,attr"`
   Initialization *Initialization `xml:"Initialization"`
   SegmentUrls    []*SegmentUrl   `xml:"SegmentURL"`
   // Navigation
   Parent *Representation `xml:"-"`
}

func (sl *SegmentList) GetTimescale() uint {
   if sl.Timescale != nil {
      return *sl.Timescale
   }
   return 1
}

func (sl *SegmentList) getParentBaseUrl() (*url.URL, error) {
   return sl.Parent.ResolveBaseUrl()
}

func (sl *SegmentList) link() {
   if sl.Initialization != nil {
      sl.Initialization.Parent = sl
   }
   for _, mediaUrl := range sl.SegmentUrls {
      mediaUrl.Parent = sl
   }
}

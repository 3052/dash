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

// SegmentUrl defines a specific media segment source.
type SegmentUrl struct {
   Media string `xml:"media,attr"`
   // Navigation
   Parent *SegmentList `xml:"-"`
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

// ResolveMedia resolves the @media attribute against the parent SegmentList's context.
func (su *SegmentUrl) ResolveMedia() (*url.URL, error) {
   base, err := su.Parent.getParentBaseUrl()
   if err != nil {
      return nil, err
   }
   return resolveRef(base, su.Media)
}

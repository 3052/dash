package dash

import "net/url"

// SegmentUrl defines a specific media segment source.
type SegmentUrl struct {
   Media string `xml:"media,attr"`
   // Navigation
   Parent *SegmentList `xml:"-"`
}

// ResolveMedia resolves the @media attribute against the parent SegmentList's context.
func (su *SegmentUrl) ResolveMedia() (*url.URL, error) {
   base, err := su.Parent.getParentBaseUrl()
   if err != nil {
      return nil, err
   }
   return resolveRef(base, su.Media)
}

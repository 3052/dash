package dash

import "net/url"

// SegmentURL defines a specific media segment source.
type SegmentURL struct {
   Media string `xml:"media,attr,omitempty"`

   // Navigation
   Parent *SegmentList `xml:"-"`
}

// ResolveMedia resolves the @media attribute against the parent SegmentList's context.
func (su *SegmentURL) ResolveMedia() (*url.URL, error) {
   base, err := su.Parent.getParentBaseURL()
   if err != nil {
      return nil, err
   }
   return resolveRef(base, su.Media)
}

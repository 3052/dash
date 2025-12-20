package dash

import "net/url"

// Initialization contains URL and byte range information for initialization segments.
type Initialization struct {
   // Used in SegmentBase
   Range string `xml:"range,attr"`
   // Used in SegmentList
   SourceUrl string `xml:"sourceURL,attr"`
   // Navigation
   Parent *SegmentList `xml:"-"`
}

// ResolveSourceUrl resolves the @sourceURL attribute against the parent SegmentList's context.
func (i *Initialization) ResolveSourceUrl() (*url.URL, error) {
   if i.Parent != nil {
      base, err := i.Parent.getParentBaseUrl()
      if err != nil {
         return nil, err
      }
      return resolveRef(base, i.SourceUrl)
   }
   return url.Parse(i.SourceUrl)
}

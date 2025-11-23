package dash

import (
   "net/url"
)

// Initialization contains URL and byte range information for initialization segments.
type Initialization struct {
   // Used in SegmentBase
   Range string `xml:"range,attr"`
   // Used in SegmentList
   SourceURL string `xml:"sourceURL,attr"`

   // Navigation
   Parent *SegmentList `xml:"-"`
}

// ResolveSourceURL resolves the @sourceURL attribute against the parent SegmentList's context.
func (i *Initialization) ResolveSourceURL() (*url.URL, error) {
   // Initialization inside SegmentList resolves against Representation BaseURL
   if i.Parent != nil {
      base, err := i.Parent.getParentBaseURL()
      if err != nil {
         return nil, err
      }
      return resolveRef(base, i.SourceURL)
   }

   // If not parented (e.g., inside SegmentBase without full linking),
   // attempt to parse the source URL directly.
   return url.Parse(i.SourceURL)
}

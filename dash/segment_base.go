package dash

import (
   "errors"
   "net/url"
   "strconv"
   "strings"
)

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

// SegmentBase defines base information for segments.
type SegmentBase struct {
   IndexRange     string          `xml:"indexRange,attr"`
   Initialization *Initialization `xml:"Initialization"`
}

func ParseRange(input string) (uint64, uint64, error) {
   startStr, endStr, found := strings.Cut(input, "-")
   if !found {
      return 0, 0, errors.New("invalid range format")
   }
   start, err := strconv.ParseUint(startStr, 10, 64)
   if err != nil {
      return 0, 0, err
   }
   end, err := strconv.ParseUint(endStr, 10, 64)
   if err != nil {
      return 0, 0, err
   }
   return start, end, nil
}

func FormatRange(start, end uint64) string {
   var sb strings.Builder
   sb.WriteString(strconv.FormatUint(start, 10))
   sb.WriteByte('-')
   sb.WriteString(strconv.FormatUint(end, 10))
   return sb.String()
}

package dash

import (
   "errors"
   "strconv"
   "strings"
)

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

// SegmentBase defines base information for segments.
type SegmentBase struct {
   IndexRange     string          `xml:"indexRange,attr"`
   Initialization *Initialization `xml:"Initialization"`
}

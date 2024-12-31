package dash

import (
   "strconv"
   "strings"
)

// the current testdata always has `-`, so lets assume for now
func (r *Range) UnmarshalText(data []byte) error {
   before, after, _ := strings.Cut(string(data), "-")
   var err error
   (*r)[0], err = strconv.ParseUint(before, 10, 64)
   if err != nil {
      return err
   }
   (*r)[1], err = strconv.ParseUint(after, 10, 64)
   if err != nil {
      return err
   }
   return nil
}

func (r Range) MarshalText() ([]byte, error) {
   b := strconv.AppendUint(nil, r[0], 10)
   b = append(b, '-')
   return strconv.AppendUint(b, r[1], 10), nil
}

// SegmentIndexBox uses:
// unsigned int(32) subsegment_duration;
// but range values can exceed 32 bits
type Range [2]uint64

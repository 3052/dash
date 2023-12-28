package dash

import (
   "errors"
   "fmt"
)

func (r Representation) Sidx_Moof() (uint32, uint32, error) {
   if r.SegmentBase == nil {
      return 0, 0, errors.New("SegmentBase")
   }
   var start uint32
   var end uint32
   _, err := fmt.Sscanf(r.SegmentBase.IndexRange, "%v-%v", &start, &end)
   if err != nil {
      return 0, 0, err
   }
   return start, end+1, nil
}

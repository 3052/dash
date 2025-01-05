package dash

import (
   "iter"
   "strconv"
   "strings"
)

type Mpd struct {
   BaseUrl                   *Url      `xml:"BaseURL"`
   MediaPresentationDuration *Duration `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

func (m *Mpd) representation() iter.Seq[Representation] {
   return func(yield func(Representation) bool) {
      for _, p := range m.Period {
         p.set(m)
         for _, adapt := range p.AdaptationSet {
            adapt.set(&p)
            for _, represent := range adapt.Representation {
               represent.set(&adapt)
               if !yield(represent) {
                  return
               }
            }
         }
      }
   }
}
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

package dash

import (
   "fmt"
   "iter"
   "net/url"
   "strconv"
   "strings"
)

func (m Media) time() bool {
   return strings.Contains(m.S, "$Time$")
}

func replace(s *string, from, to string) {
   *s = strings.Replace(*s, from, to, 1)
}

type Media struct {
   S string
}

func (m *Media) UnmarshalText(data []byte) error {
   m.S = string(data)
   return nil
}

func (m Media) Url(r *Representation, i int) (*url.URL, error) {
   replace(&m.S, "$RepresentationID$", r.Id)
   if m.time() {
      replace(&m.S, "$Time$", fmt.Sprint(i))
   } else {
      replace(&m.S, "$Number$", fmt.Sprint(i))
      replace(&m.S, "$Number%02d$", fmt.Sprintf("%02d", i))
      replace(&m.S, "$Number%03d$", fmt.Sprintf("%03d", i))
      replace(&m.S, "$Number%04d$", fmt.Sprintf("%04d", i))
      replace(&m.S, "$Number%05d$", fmt.Sprintf("%05d", i))
      replace(&m.S, "$Number%06d$", fmt.Sprintf("%06d", i))
      replace(&m.S, "$Number%07d$", fmt.Sprintf("%07d", i))
      replace(&m.S, "$Number%08d$", fmt.Sprintf("%08d", i))
      replace(&m.S, "$Number%09d$", fmt.Sprintf("%09d", i))
   }
   u, err := url.Parse(m.S)
   if err != nil {
      return nil, err
   }
   if r.BaseUrl != nil {
      u = r.BaseUrl.Url.ResolveReference(u)
   }
   return u, nil
}
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

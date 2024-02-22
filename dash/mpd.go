package dash

import "fmt"

// media presentation description
// wikipedia.org/wiki/Dynamic_Adaptive_Streaming_over_HTTP
type MPD struct {
   Period []Period
}

func (m MPD) String() string {
   var b []byte
   for _, p := range m.Period {
      for _, a := range p.AdaptationSet {
         for _, r := range a.Representation {
            if b != nil {
               b = append(b, "\n\n"...)
            }
            var c []byte
            if r.Width >= 1 {
               c = fmt.Append(c, "width = ", r.Width)
            }
            if r.Height >= 1 {
               if c != nil {
                  c = append(c, '\n')
               }
               c = fmt.Append(c, "height = ", r.Height)
            }
            if c != nil {
               c = append(c, '\n')
            }
            c = fmt.Append(c, "bandwidth = ", r.Bandwidth)
            if r.Codecs != "" {
               c = fmt.Append(c, "\ncodecs = ", r.Codecs)
            }
            if a.Codecs != "" {
               c = fmt.Append(c, "\ncodecs = ", a.Codecs)
            }
            if r.MimeType != "" {
               c = fmt.Append(c, "\ntype = ", r.MimeType)
            } else {
               c = fmt.Append(c, "\ntype = ", a.MimeType)
            }
            if a.Role != nil {
               c = fmt.Append(c, "\nrole = ", a.Role.Value)
            }
            if a.Lang != "" {
               c = fmt.Append(c, "\nlang = ", a.Lang)
            }
            c = fmt.Append(c, "\nid = ", r.ID)
            if p.ID != "" {
               c = fmt.Append(c, "\nperiod = ", p.ID)
            }
            b = append(b, c...)
         }
      }
   }
   return string(b)
}

// godocs.io/flag#Visit
func (m MPD) Visit(f func(Pointer)) {
   m.Contains(func(p Pointer) bool {
      f(p)
      return false
   })
}

func (m MPD) Contains(f func(Pointer) bool) bool {
   for _, period := range m.Period {
      for _, adapt := range period.AdaptationSet {
         for _, represent := range adapt.Representation {
            var p Pointer
            p.AdaptationSet = &adapt
            p.Period = &period
            p.Representation = &represent
            if f(p) {
               return true
            }
         }
      }
   }
   return false
}

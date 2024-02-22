package dash

import "fmt"

// media presentation description
// wikipedia.org/wiki/Dynamic_Adaptive_Streaming_over_HTTP
type MPD struct {
   Period []Period
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

func (m MPD) Visit(f func(Pointer)) {
   m.Contains(func(p Pointer) bool {
      f(p)
      return false
   })
}

func (m MPD) String() string {
   var b []byte
   m.Visit(func(p Pointer) {
      if b != nil {
         b = append(b, "\n\n"...)
      }
      var c []byte
      if v := p.Representation.Width; v >= 1 {
         c = fmt.Append(c, "width = ", v)
      }
      if v := p.Representation.Height; v >= 1 {
         if c != nil {
            c = append(c, '\n')
         }
         c = fmt.Append(c, "height = ", v)
      }
      if c != nil {
         c = append(c, '\n')
      }
      c = fmt.Append(c, "bandwidth = ", p.Representation.Bandwidth)
      if v, ok := p.codecs(); ok {
         c = fmt.Append(c, "\ncodecs = ", v)
      }
      c = fmt.Append(c, "\ntype = ", p.mime_type())
      if v := p.AdaptationSet.Role; v != nil {
         c = fmt.Append(c, "\nrole = ", v.Value)
      }
      if v := p.AdaptationSet.Lang; v != "" {
         c = fmt.Append(c, "\nlang = ", v)
      }
      c = fmt.Append(c, "\nid = ", p.Representation.ID)
      if v := p.Period.ID; v != "" {
         c = fmt.Append(c, "\nperiod = ", v)
      }
      b = append(b, c...)
   })
   return string(b)
}

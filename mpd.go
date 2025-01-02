package dash

import "iter"

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

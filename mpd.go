package dash

import "iter"

type Mpd struct {
   MediaPresentationDuration *Duration `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

func (m *Mpd) representation() iter.Seq[Representation] {
   return func(yield func(Representation) bool) {
      for _, p := range m.Period {
         p.mpd = m
         for _, adapt := range p.AdaptationSet {
            adapt.period = &p
            for _, represent := range adapt.Representation {
               represent.adaptation_set = &adapt
               represent.set()
               if !yield(represent) {
                  return
               }
            }
         }
      }
   }
}

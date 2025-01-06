package dash

import "iter"

func (m *Mpd) Representation() iter.Seq[Representation] {
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

func (r *Representation) Representation() iter.Seq[Representation] {
   return func(yield func(Representation) bool) {
      for r2 := range r.adaptation_set.period.mpd.Representation() {
         if r2.Id == r.Id {
            if !yield(r2) {
               return
            }
         }
      }
   }
}

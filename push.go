package dash

import "iter"

func (p *Period) hd_push() (*Representation, bool) {
   for value := range p.push() {
      if value.Height > 576 {
         return &value, true
      }
   }
   return nil, false
}

func (p *Period) push() iter.Seq[Representation] {
   return func(yield func(Representation) bool) {
      for _, adapt := range p.AdaptationSet {
         for _, represent := range adapt.Representation {
            if !yield(represent) {
               return
            }
         }
      }
   }
}

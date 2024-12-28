package dash

func (p *Period) hd_pull() (*Representation, bool) {
   values := p.pull()
   for {
      value, ok := values()
      if !ok {
         return nil, false
      }
      if value.Height > 576 {
         return value, true
      }
   }
}

func (p *Period) pull() func() (*Representation, bool) {
   var a, b int
   return func() (*Representation, bool) {
      for a < len(p.AdaptationSet) {
         adapt := p.AdaptationSet[a]
         for b < len(adapt.Representation) {
            represent := adapt.Representation[b]
            b++
            return &represent, true
         }
         a++
         b = 0
      }
      return nil, false
   }
}

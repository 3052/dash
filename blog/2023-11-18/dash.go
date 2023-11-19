package dash

type representation struct {
   adaptationSet *adaptationSet
   contentProtection []int
}

type adaptationSet struct {
   contentProtection []int
   representation []*representation
}

func (a adaptationSet) setContentProtection() {
   for _, r := range a.representation {
      r.contentProtection = a.contentProtection
   }
}

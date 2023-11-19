package dash

import (
   "fmt"
   "testing"
)

func Test_DASH(t *testing.T) {
   a := adaptationSet{
      contentProtection: []int{1,2},
      representation: []*representation{
         {},
         {},
      },
   }
   a.setContentProtection()
   for _, r := range a.representation {
      fmt.Printf("%+v\n", r)
   }
}

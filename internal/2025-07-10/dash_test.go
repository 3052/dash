package dash

import (
   "fmt"
   "slices"
   "testing"
)

func Test(t *testing.T) {
   for _, state := range states {
      if slices.Contains(state.example, "criterion.mpd") {
         fmt.Printf("%#q,\n", state.state)
      }
   }
}

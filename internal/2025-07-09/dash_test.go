package dash

import (
   "fmt"
   "slices"
   "testing"
)

func Test(t *testing.T) {
   for _, state := range states {
      if slices.Contains(state.example, "canal.mpd") {
         fmt.Println(state.state)
      }
   }
}

package dash

import (
   "fmt"
   "testing"
)

func TestDuration(t *testing.T) {
   for _, name := range tests {
      reps, err := reader(name)
      if err != nil {
         t.Fatal(err)
      }
      duration, err := reps[0].adaptation_set.period.mpd.seconds()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(name, duration)
   }
}

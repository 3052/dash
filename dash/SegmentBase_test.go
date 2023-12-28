package dash

import (
   "fmt"
   "testing"
)

func Test_SegmentBase(t *testing.T) {
   reps, err := read_file("mpd/hulu.mpd")
   if err != nil {
      t.Fatal(err)
   }
   for _, rep := range reps {
      fmt.Println(rep.Sidx_Moof())
   }
}

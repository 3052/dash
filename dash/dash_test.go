package dash

import (
   "fmt"
   "os"
   "testing"
)

func TestRange(t *testing.T) {
   reps, err := reader("mpd/hulu.mpd")
   if err != nil {
      t.Fatal(err)
   }
   for _, rep := range reps {
      a, b, ok := rep.SegmentBase.Initialization.Range.Cut()
      fmt.Print(a, " ", b, " ", ok, " ")
      a, b, ok = rep.SegmentBase.IndexRange.Cut()
      fmt.Print(a, " ", b, " ", ok, "\n")
   }
}

func reader(name string) ([]Representation, error) {
   text, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   return Unmarshal(text)
}

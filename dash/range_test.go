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
      start, end, err := rep.SegmentBase.Initialization.Range.Scan()
      fmt.Print(start, " ", end, " ", err, " ")
      start, end, err = rep.SegmentBase.IndexRange.Scan()
      fmt.Print(start, " ", end, " ", err, "\n")
   }
}

func reader(name string) ([]Representation, error) {
   text, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   return Unmarshal(text)
}

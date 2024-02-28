package dash

import (
   "fmt"
   "os"
   "testing"
)

func reader(name string) ([]Representation, error) {
   text, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   return Unmarshal(text)
}

func TestRange(t *testing.T) {
   reps, err := reader("mpd/hulu.mpd")
   if err != nil {
      t.Fatal(err)
   }
   for _, rep := range reps {
      start, end, err := rep.SegmentBase.Initialization.Range.Scan()
      fmt.Printf("%v %v %v ", start, end, err)
      start, end, err = rep.SegmentBase.IndexRange.Scan()
      fmt.Printf("%v %v %v\n", start, end, err)
   }
}

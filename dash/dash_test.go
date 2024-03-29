package dash

import (
   "fmt"
   "os"
   "slices"
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

var tests = []string{
   "mpd/amc.mpd",
   "mpd/hulu.mpd",
   "mpd/mubi.mpd",
   "mpd/nbc.mpd",
   "mpd/paramount.mpd",
   "mpd/roku.mpd",
   "mpd/stan.mpd",
}

func reader(name string) ([]Representation, error) {
   text, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   return Unmarshal(text)
}

func TestDelete(t *testing.T) {
   for i, name := range tests {
      if i >= 1 {
         fmt.Println()
      }
      reps, err := reader(name)
      if err != nil {
         t.Fatal(err)
      }
      reps = slices.DeleteFunc(reps, func(r Representation) bool {
         if _, ok := r.Ext(); !ok {
            return true
         }
         return false
      })
      for i, rep := range reps {
         if i >= 1 {
            fmt.Println()
         }
         fmt.Println(rep)
      }
   }
}

func TestRange(t *testing.T) {
   reps, err := reader("mpd/hulu.mpd")
   if err != nil {
      t.Fatal(err)
   }
   for _, rep := range reps {
      r, err := rep.SegmentBase.Initialization.Range.Scan()
      fmt.Print(r.Start, " ", r.End, " ", err, " ")
      r, err = rep.SegmentBase.IndexRange.Scan()
      fmt.Print(r.Start, " ", r.End, " ", err, "\n")
   }
}

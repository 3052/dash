package dash

import (
   "fmt"
   "testing"
)

func TestProtection(t *testing.T) {
   for _, test := range tests {
      reps, err := reader(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, rep := range reps {
         _, pssh := rep.PSSH()
         _, kid := rep.Default_KID()
         fmt.Printf(
            "mpd:%v period:%q type:%v pssh:%v kid:%v\n",
            test, rep.adaptation_set.period.ID, rep.mime_type(), pssh, kid,
         )
      }
   }
}

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

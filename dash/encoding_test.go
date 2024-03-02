package dash

import (
   "fmt"
   "testing"
)

func TestRange(t *testing.T) {
   reps, err := reader("mpd/hulu.mpd")
   if err != nil {
      t.Fatal(err)
   }
   for _, rep := range reps {
      a, b, err := rep.SegmentBase.Initialization.Range.Scan()
      fmt.Print(a, " ", b, " ", err, " ")
      a, b, err = rep.SegmentBase.IndexRange.Scan()
      fmt.Print(a, " ", b, " ", err, "\n")
   }
}

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

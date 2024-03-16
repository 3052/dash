package dash

import (
   "fmt"
   "testing"
)

func TestPsshKid(t *testing.T) {
   for _, test := range tests {
      reps, err := reader(test)
      if err != nil {
         t.Fatal(err)
      }
      for i, rep := range reps {
         if i >= 1 {
            fmt.Println()
         }
         protect := rep.content_protection()
         fmt.Println("mpd =", test)
         fmt.Println("protect == nil", protect == nil)
         fmt.Println("type =", rep.mime_type())
         if protect != nil {
            _, pssh := rep.PSSH()
            _, kid := rep.Default_KID()
            fmt.Println("kid =", kid)
            fmt.Println("pssh =", pssh)
         }
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

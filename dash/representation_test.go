package dash

import (
   "fmt"
   "os"
   "testing"
)

var tests = []string{
   "mpd/amc.mpd",
   "mpd/hulu.mpd",
   "mpd/nbc.mpd",
   "mpd/paramount.mpd",
   "mpd/roku.mpd",
}

func read_file(name string) ([]*Representation, error) {
   file, err := os.Open(name)
   if err != nil {
      return nil, err
   }
   defer file.Close()
   var m Media
   if err := m.Decode(file); err != nil {
      return nil, err
   }
   return m.Representation(tests[name])
}

func Test_Info(t *testing.T) {
   for name := range tests {
      reps, err := read_file(name)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      for i, rep := range reps {
         if i >= 1 {
            fmt.Println()
         }
         fmt.Println(rep)
      }
      fmt.Println()
   }
}

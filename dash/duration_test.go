package dash

import (
   "fmt"
   "os"
   "testing"
   "time"
)

func TestDuration(t *testing.T) {
   for _, name := range tests {
      reps, err := reader(name)
      if err != nil {
         t.Fatal(err)
      }
      raw := reps[0].adaptation_set.period.mpd.MediaPresentationDuration
      fmt.Printf("%v %q\n", name, raw)
      duration, err := time.ParseDuration(raw)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(duration)
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

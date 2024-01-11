package dash

import (
   "fmt"
   "net/http"
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

func Test_Initialization(t *testing.T) {
   reps, err := read_file("mpd/amc.mpd")
   if err != nil {
      t.Fatal(err)
   }
   for _, rep := range reps {
      v, ok := rep.Initialization()
      fmt.Printf("%v %q %v\n\n", rep.ID, v, ok)
   }
}

func Test_Media(t *testing.T) {
   reps, err := read_file("mpd/roku.mpd")
   if err != nil {
      t.Fatal(err)
   }
   base, err := http.NewRequest("", "http://example.com", nil)
   if err != nil {
      t.Fatal(err)
   }
   media, ok := reps[0].Media()
   if !ok {
      t.Fatal("Media")
   }
   for _, medium := range media {
      req, err := http.NewRequest("", medium, nil)
      if err != nil {
         t.Fatal(err)
      }
      req.URL = base.URL.ResolveReference(req.URL)
      fmt.Println(req.URL)
   }
}

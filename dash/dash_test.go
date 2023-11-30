package dash

import (
   "fmt"
   "net/http"
   "os"
   "testing"
)

func read_file(name string) ([]*Representation, error) {
   file, err := os.Open(name)
   if err != nil {
      return nil, err
   }
   defer file.Close()
   return Representations(file)
}

var tests = []string{
   "mpd/amc.mpd",
   "mpd/hulu.mpd",
   "mpd/nbc.mpd",
   "mpd/paramount.mpd",
   "mpd/roku.mpd",
}

func Test_Ext(t *testing.T) {
   for _, name := range tests {
      reps, err := read_file(name)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      for _, rep := range reps {
         v, ok := rep.Ext()
         fmt.Printf("%q %v\n", v, ok)
      }
      fmt.Println()
   }
}

func Test_Info(t *testing.T) {
   for _, name := range tests {
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

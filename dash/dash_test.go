package dash

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

func reader(name string) (*Media, error) {
   file, err := os.Open(name)
   if err != nil {
      return nil, err
   }
   defer file.Close()
   var m Media
   if err := m.Decode(file); err != nil {
      return nil, err
   }
   return &m, nil
}

func Test_SegmentBase(t *testing.T) {
   m, err := reader("mpd/hulu.mpd")
   if err != nil {
      t.Fatal(err)
   }
   looper(m, func(r Representation) bool {
      fmt.Println(r.Sidx_Moof())
      return true
   })
}

func Test_Initialization(t *testing.T) {
   m, err := reader("mpd/amc.mpd")
   if err != nil {
      t.Fatal(err)
   }
   looper(m, func(r Representation) bool {
      v, ok := r.Initialization()
      fmt.Printf("%v %q %v\n\n", r.ID, v, ok)
      return true
   })
}

func looper(m *Media, fn func(Representation) bool) {
   for _, period := range m.Period {
      for _, adaptation := range period.AdaptationSet {
         for _, representation := range adaptation.Representation {
            if !fn(representation) {
               return
            }
         }
      }
   }
}

func Test_Media(t *testing.T) {
   media_struct, err := reader("mpd/roku.mpd")
   if err != nil {
      t.Fatal(err)
   }
   base, err := url.Parse("http://example.com")
   if err != nil {
      t.Fatal(err)
   }
   looper(media_struct, func(r Representation) bool {
      media_string, ok := r.Media()
      if !ok {
         t.Fatal("Representation.Media")
      }
      for _, medium := range media_string {
         ref, err := base.Parse(medium)
         if err != nil {
            t.Fatal(err)
         }
         fmt.Println(ref)
      }
      return false
   })
}

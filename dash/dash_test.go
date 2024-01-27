package dash

import (
   "encoding/xml"
   "fmt"
   "net/url"
   "os"
   "testing"
)

func Test_SegmentBase(t *testing.T) {
   media, err := reader("mpd/hulu.mpd")
   if err != nil {
      t.Fatal(err)
   }
   media.Every(func(r Representation) {
      fmt.Printf("%+v\n", r.SegmentBase)
   })
}

func reader(name string) (*MPD, error) {
   text, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   media := new(MPD)
   if err := xml.Unmarshal(text, &media); err != nil {
      return nil, err
   }
   return media, nil
}

func Test_Initialization(t *testing.T) {
   media, err := reader("mpd/amc.mpd")
   if err != nil {
      t.Fatal(err)
   }
   media.Every(func(r Representation) {
      v, ok := r.Initialization()
      fmt.Printf("%v %q %v\n\n", r.ID, v, ok)
   })
}

func Test_Media(t *testing.T) {
   roku, err := reader("mpd/roku.mpd")
   if err != nil {
      t.Fatal(err)
   }
   base, err := url.Parse("http://example.com")
   if err != nil {
      t.Fatal(err)
   }
   roku.Some(func(r Representation) bool {
      media, ok := r.Media()
      if !ok {
         t.Fatal("Representation.Media")
      }
      for _, medium := range media {
         ref, err := base.Parse(medium)
         if err != nil {
            t.Fatal(err)
         }
         fmt.Println(ref)
      }
      return true
   })
}

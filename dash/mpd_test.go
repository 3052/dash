package dash

import (
   "fmt"
   "net/url"
   "testing"
)

func TestInitialization(t *testing.T) {
   media, err := reader("mpd/amc.mpd")
   if err != nil {
      t.Fatal(err)
   }
   media.Visit(func(p Pointer) {
      v, ok := p.Initialization()
      fmt.Printf("%v %q %v\n\n", p.Representation.ID, v, ok)
   })
}

func TestMedia(t *testing.T) {
   roku, err := reader("mpd/roku.mpd")
   if err != nil {
      t.Fatal(err)
   }
   base, err := url.Parse("http://example.com")
   if err != nil {
      t.Fatal(err)
   }
   roku.Contains(func(p Pointer) bool {
      for _, raw := range p.Media() {
         medium, err := base.Parse(raw)
         if err != nil {
            t.Fatal(err)
         }
         fmt.Println(medium)
      }
      return true
   })
}

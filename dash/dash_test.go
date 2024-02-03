package dash

import (
   "encoding/xml"
   "fmt"
   "net/url"
   "os"
   "testing"
)

func Test_Range(t *testing.T) {
   media, err := reader("mpd/hulu.mpd")
   if err != nil {
      t.Fatal(err)
   }
   media.Every(func(p Pointer) {
      sb := p.Representation.SegmentBase
      start, end, err := sb.Initialization.Range.Scan()
      fmt.Printf("%v %v %v ", start, end, err)
      start, end, err = sb.IndexRange.Scan()
      fmt.Printf("%v %v %v\n", start, end, err)
   })
}

var tests = []string{
   "mpd/amc.mpd",
   "mpd/hulu.mpd",
   "mpd/nbc.mpd",
   "mpd/paramount.mpd",
   "mpd/roku.mpd",
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

func Test_Media(t *testing.T) {
   roku, err := reader("mpd/roku.mpd")
   if err != nil {
      t.Fatal(err)
   }
   base, err := url.Parse("http://example.com")
   if err != nil {
      t.Fatal(err)
   }
   roku.Some(func(p Pointer) bool {
      for _, ref := range p.Media() {
         media, err := base.Parse(ref)
         if err != nil {
            t.Fatal(err)
         }
         fmt.Println(media)
      }
      return false
   })
}

func Test_Initialization(t *testing.T) {
   media, err := reader("mpd/amc.mpd")
   if err != nil {
      t.Fatal(err)
   }
   media.Every(func(p Pointer) {
      v, ok := p.Initialization()
      fmt.Printf("%v %q %v\n\n", p.Representation.ID, v, ok)
   })
}

package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestString(t *testing.T) {
   for i, name := range tests {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(name)
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      var media MPD
      if err := xml.Unmarshal(text, &media); err != nil {
         t.Fatal(err)
      }
      fmt.Println(media)
   }
}

func TestMedia(t *testing.T) {
   tests := [][]string{
      // startNumber == nil
      {"mpd/mubi.mpd", "dallas-edge3.mubicdn.net/stream/7a947d0bd29b2e5c17ec36497b1c67e4/65d5e449/da8710c0/mubi-films/325455/passages_eng_zxx_1800x1080_50000_mezz40828/ae8c88ed4e/drm_playlist.0ff148ef80.ism/default/"},
      // startNumber == 0
      {"mpd/amc.mpd", ""},
      // startNumber == 1
      {"mpd/paramount.mpd", "vod-gcs-cedexis.cbsaavideo.com/intl_vms/2022/02/24/2006197315671/77016_cenc_dash/"},
   }
   for _, test := range tests {
      fmt.Println(test[0] + ":")
      media, err := reader(test[0])
      if err != nil {
         t.Fatal(err)
      }
      media.Contains(func(p Pointer) bool {
         for _, medium := range p.Media() {
            fmt.Println(test[1] + medium)
         }
         return true
      })
   }
}

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

func TestRange(t *testing.T) {
   media, err := reader("mpd/hulu.mpd")
   if err != nil {
      t.Fatal(err)
   }
   media.Visit(func(p Pointer) {
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

func TestProtection(t *testing.T) {
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media MPD
      if err := xml.Unmarshal(text, &media); err != nil {
         t.Fatal(err)
      }
      media.Visit(func(p Pointer) {
         _, pssh := p.PSSH()
         _, kid := p.Default_KID()
         fmt.Printf(
            "mpd:%v period:%q type:%v pssh:%v kid:%v\n",
            test, p.Period.ID, p.mime_type(), pssh, kid,
         )
      })
   }
}

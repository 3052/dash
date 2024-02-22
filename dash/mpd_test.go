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


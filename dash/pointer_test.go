package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func Test_Pointer(t *testing.T) {
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media MPD
      if err := xml.Unmarshal(text, &media); err != nil {
         t.Fatal(err)
      }
      media.Every(func(p Pointer) {
         fmt.Printf("name:%v period:%q ", test, p.Period.ID)
         _, ok := p.Default_KID()
         fmt.Printf("kid:%v ", ok)
         _, ok = p.PSSH()
         fmt.Printf("pssh:%v ", ok)
         fmt.Printf("mimeType:%q\n", p.MimeType)
      })
   }
}

var tests = []string{
   "mpd/amc.mpd",
   "mpd/hulu.mpd",
   "mpd/nbc.mpd",
   "mpd/paramount.mpd",
   "mpd/roku.mpd",
}

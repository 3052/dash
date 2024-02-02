package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func Test_ContentProtection(t *testing.T) {
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media MPD
      if err := xml.Unmarshal(text, &media); err != nil {
         t.Fatal(err)
      }
      media.Every(func(m MPD, i Index) {
         fmt.Printf("%v %q ", test, i.GetPeriod(m).ID)
         _, err := represent.Default_KID()
         fmt.Print(err, " ")
         _, err = represent.PSSH()
         fmt.Print(err, " ")
         fmt.Printf("%q %q\n", adapt.MimeType, represent.MimeType)
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

package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

var tests = []string{
   "mpd/amc.mpd",
   "mpd/hulu.mpd",
   "mpd/nbc.mpd",
   "mpd/paramount.mpd",
   "mpd/roku.mpd",
}

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
      for _, period := range media.Period {
         for _, adapt := range period.AdaptationSet {
            for _, represent := range adapt.Representation {
               fmt.Printf("%v %q ", test, period.ID)
               _, err := represent.Default_KID()
               fmt.Print(err, " ")
               _, err = represent.PSSH()
               fmt.Print(err, " ")
               fmt.Printf("%q %q\n", adapt.MimeType, represent.MimeType)
            }
         }
      }
   }
}

package dash

import (
   "fmt"
   "os"
   "testing"
)

var tests = []string{
   "testdata/amc.mpd",
   "testdata/cine-member.mpd",
   "testdata/criterion.mpd",
   "testdata/ctv.mpd",
   "testdata/draken.mpd",
   "testdata/hulu.mpd",
   "testdata/max.mpd",
   "testdata/mubi.mpd",
   "testdata/nbc.mpd",
   "testdata/paramount.mpd",
   "testdata/plex.mpd",
   "testdata/pluto.mpd",
   "testdata/rakuten.mpd",
   "testdata/roku.mpd",
   "testdata/stan.mpd",
   "testdata/tubi.mpd",
}

func TestMedia(t *testing.T) {
   for _, test := range tests {
      fmt.Println(test)
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media Mpd
      err = media.Unmarshal(text)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            for _, represent := range v.Representation {
               if v, ok := represent.GetSegmentTemplate(); ok {
                  initial, ok := v.GetInitialization(represent)
                  fmt.Printf("%q %v\n", initial, ok)
                  media, err := v.GetMedia(represent)
                  if err != nil {
                     t.Fatal(err)
                  }
                  for i := range min(len(media), 9) {
                     fmt.Println(media[i])
                  }
               }
            }
         }
      }
   }
}

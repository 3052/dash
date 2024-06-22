package dash

import (
   "fmt"
   "os"
   "testing"
)

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

var tests = []string{
   "mpd/amc.mpd",
   "mpd/cine-member.mpd",
   "mpd/hulu.mpd",
   "mpd/mubi.mpd",
   "mpd/nbc.mpd",
   "mpd/paramount.mpd",
   "mpd/plex.mpd",
   "mpd/pluto.mpd",
   "mpd/roku.mpd",
   "mpd/stan.mpd",
   "mpd/tubi.mpd",
}

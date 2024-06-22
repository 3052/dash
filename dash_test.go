package dash

import (
   "os"
   "testing"
)

func TestMpd(t *testing.T) {
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media Mpd
      err = media.Unmarshal(text)
      if err != nil {
         t.Fatal(err)
      }
      if media.MediaPresentationDuration == "" {
         t.Fatal("MediaPresentationDuration", test)
      }
      if len(media.Period) == 0 {
         t.Fatal("Period", test)
      }
      for _, v := range media.Period {
         if v.mpd == nil {
            t.Fatal("mpd")
         }
         for _, v := range v.AdaptationSet {
            if v.period == nil {
               t.Fatal("period")
            }
            for _, v := range v.Representation {
               if v.adaptation_set == nil {
                  t.Fatal("adaptation_set")
               }
            }
         }
      }
   }
}

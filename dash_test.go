package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestChapters(t *testing.T) {
   text, err := os.ReadFile("testdata/max.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = media.Unmarshal(text)
   if err != nil {
      t.Fatal(err)
   }
   var line bool
   for _, v := range media.Period {
      for _, v := range v.AdaptationSet {
         for _, v := range v.Representation {
            if line {
               fmt.Println()
            } else {
               line = true
            }
            fmt.Println(v)
         }
      }
   }
}

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

func TestPeriod(t *testing.T) {
   duration := make(set)
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media Mpd
      err = xml.Unmarshal(text, &media)
      if err != nil {
         t.Fatal(err)
      }
      for _, p := range media.Period {
         if len(p.AdaptationSet) == 0 {
            t.Fatal("AdaptationSet", test)
         }
         if p.Duration != nil {
            duration[1] = struct{}{}
         } else {
            duration[0] = struct{}{}
         }
      }
   }
   fmt.Println(duration)
}

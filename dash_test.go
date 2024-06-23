package dash

import (
   "encoding/xml"
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

func TestDuration(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("MPD %q\n", media.MediaPresentationDuration)
      for _, v := range media.Period {
         fmt.Printf("Period %q\n", v.Duration)
      }
      fmt.Println()
   }
}

func TestMimeType(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            fmt.Printf("AdaptationSet %q\n", v.MimeType)
            for _, v := range v.Representation {
               fmt.Printf("Representation %q\n", v.MimeType)
            }
            fmt.Println()
         }
      }
   }
}

func TestCodecs(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            fmt.Printf("AdaptationSet %q\n", v.Codecs)
            for _, v := range v.Representation {
               fmt.Printf("Representation %q\n", v.Codecs)
            }
         }
         fmt.Println()
      }
   }
}

func TestSegmentBase(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            for _, v := range v.Representation {
               fmt.Printf("%+v\n", v.SegmentBase)
            }
         }
      }
   }
}

func TestSegmentTemplate(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            fmt.Printf("AdaptationSet %+v\n", v.SegmentTemplate)
            for _, v := range v.Representation {
               fmt.Printf("Representation %+v\n", v.SegmentTemplate)
            }
         }
      }
   }
}

func TestRole(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            fmt.Println(v.Role)
         }
      }
   }
}

func new_mpd(name string) (*Mpd, error) {
   text, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   media := new(Mpd)
   err = xml.Unmarshal(text, media)
   if err != nil {
      return nil, err
   }
   return media, nil
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

package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "slices"
   "testing"
)

var tests = []string{
   "mpd/amc.mpd",
   "mpd/hulu.mpd",
   "mpd/mubi.mpd",
   "mpd/nbc.mpd",
   "mpd/paramount.mpd",
   "mpd/peacock.mpd",
   "mpd/roku-clear.mpd",
   "mpd/roku-protected.mpd",
   "mpd/stan.mpd",
}

func TestMpd(t *testing.T) {
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media mpd
      xml.Unmarshal(text, &media)
      if media.MediaPresentationDuration == "" {
         t.Fatal("MediaPresentationDuration", test)
      }
      if len(media.Period) == 0 {
         t.Fatal("Period", test)
      }
   }
}

func TestPeriod(t *testing.T) {
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      reps, err := Unmarshal(text)
      if err != nil {
         t.Fatal(err)
      }
      for _, rep := range reps {
         if rep.adaptation_set.period.mpd == nil {
            t.Fatal("mpd", test)
         }
      }
      var media mpd
      xml.Unmarshal(text, &media)
      for _, p := range media.Period {
         if len(p.AdaptationSet) == 0 {
            t.Fatal("AdaptationSet", test)
         }
      }
   }
}

//////////////////////////////////////////////////////

func TestAdaptation(t *testing.T) {
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      reps, err := Unmarshal(text)
      if err != nil {
         t.Fatal(err)
      }
      for _, rep := range reps {
         if _, ok := rep.GetCodecs(); !ok {
            fmt.Println("GetCodecs needed")
         }
      }
      var media mpd
      xml.Unmarshal(text, &media)
      for _, per := range media.Period {
         for _, ada := range per.AdaptationSet {
            if ada.Lang == "" {
               t.Fatal("Lang")
            }
         }
      }
   }
}

func TestDelete(t *testing.T) {
   for i, name := range tests {
      if i >= 1 {
         fmt.Println()
      }
      reps, err := reader(name)
      if err != nil {
         t.Fatal(err)
      }
      reps = slices.DeleteFunc(reps, func(r Representation) bool {
         if _, ok := r.Ext(); !ok {
            return true
         }
         return false
      })
      for i, rep := range reps {
         if i >= 1 {
            fmt.Println()
         }
         fmt.Println(rep)
      }
   }
}

func reader(name string) ([]Representation, error) {
   text, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   return Unmarshal(text)
}
func TestRange(t *testing.T) {
   reps, err := reader("mpd/hulu.mpd")
   if err != nil {
      t.Fatal(err)
   }
   for _, rep := range reps {
      start, end, err := rep.SegmentBase.Initialization.Range.Scan()
      fmt.Print(start, " ", end, " ", err, " ")
      start, end, err = rep.SegmentBase.IndexRange.Scan()
      fmt.Print(start, " ", end, " ", err, "\n")
   }
}

func TestInitialization(t *testing.T) {
   reps, err := reader("mpd/amc.mpd")
   if err != nil {
      t.Fatal(err)
   }
   for _, rep := range reps {
      v, ok := rep.Initialization()
      fmt.Printf("%v %q %v\n\n", rep.ID, v, ok)
   }
}

func TestDuration(t *testing.T) {
   for _, name := range tests {
      reps, err := reader(name)
      if err != nil {
         t.Fatal(err)
      }
      duration, err := reps[0].adaptation_set.period.mpd.seconds()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(name, duration)
   }
}

var media_tests = []struct{
   name string
   base string
}{
   // startNumber == nil
   {"mpd/mubi.mpd", "new-york-edge2.mubicdn.net/stream/43cac9f0138aaa566a429be4542ff21c/65df1dc5/728eb9fc/mubi-films/325455/passages_eng_zxx_1800x1080_50000_mezz40828/ae8c88ed4e/drm_playlist.0ff148ef80.ism/default/"},
   // startNumber == 0
   {"mpd/amc.mpd", ""},
   // startNumber == 1
   {"mpd/paramount.mpd", "vod-gcs-cedexis.cbsaavideo.com/intl_vms/2022/02/24/2006197315671/77016_cenc_dash/"},
}

func TestMedia(t *testing.T) {
   for _, test := range media_tests {
      fmt.Println(test.name + ":")
      reps, err := reader(test.name)
      if err != nil {
         t.Fatal(err)
      }
      for _, media := range reps[0].Media() {
         fmt.Println(test.base + media)
      }
   }
}

package dash

import (
   "fmt"
   "testing"
)

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

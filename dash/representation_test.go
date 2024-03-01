package dash

import (
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

func TestString(t *testing.T) {
   for i, name := range tests {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(name)
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      reps, err := Unmarshal(text)
      if err != nil {
         t.Fatal(err)
      }
      for i, rep := range reps {
         if i >= 1 {
            fmt.Println()
         }
         fmt.Println(rep)
      }
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

var media_tests = [][]string{
   // startNumber == nil
   {"mpd/mubi.mpd", "new-york-edge2.mubicdn.net/stream/43cac9f0138aaa566a429be4542ff21c/65df1dc5/728eb9fc/mubi-films/325455/passages_eng_zxx_1800x1080_50000_mezz40828/ae8c88ed4e/drm_playlist.0ff148ef80.ism/default/"},
   // startNumber == 0
   {"mpd/amc.mpd", ""},
   // startNumber == 1
   {"mpd/paramount.mpd", "vod-gcs-cedexis.cbsaavideo.com/intl_vms/2022/02/24/2006197315671/77016_cenc_dash/"},
}

func TestMedia(t *testing.T) {
   for _, test := range media_tests {
      fmt.Println(test[0] + ":")
      reps, err := reader(test[0])
      if err != nil {
         t.Fatal(err)
      }
      for _, media := range reps[0].Media() {
         fmt.Println(test[1] + media)
      }
   }
}

func TestProtection(t *testing.T) {
   for _, test := range tests {
      reps, err := reader(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, rep := range reps {
         _, pssh := rep.PSSH()
         _, kid := rep.Default_KID()
         fmt.Printf(
            "mpd:%v period:%q type:%v pssh:%v kid:%v\n",
            test, rep.adaptation_set.period.ID, rep.mime_type(), pssh, kid,
         )
      }
   }
}

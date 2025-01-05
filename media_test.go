package dash

import (
   "encoding/xml"
   "os"
   "testing"
)

func Test2MediaUrl(t *testing.T) {
   _, err := Media{"\n"}.Url(&Representation{}, 0)
   if err == nil {
      t.Fatal("Media.Url")
   }
}

func TestMediaUrl(t *testing.T) {
   data, err := os.ReadFile("testdata/itv.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = xml.Unmarshal(data, &media)
   if err != nil {
      t.Fatal(err)
   }
   var represent Representation
   for represent = range media.Representation() {
      break
   }
   _, err = represent.SegmentTemplate.Media.Url(&represent, 0)
   if err != nil {
      t.Fatal(err)
   }
}

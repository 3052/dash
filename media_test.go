package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestMedia(t *testing.T) {
   data, err := os.ReadFile("testdata/amc.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var present Mpd
   err = xml.Unmarshal(data, &present)
   if err != nil {
      t.Fatal(err)
   }
   represent, _ := present.hd()
   fmt.Printf("%q\n", represent.SegmentTemplate.Media("HELLO", 999))
}

package dash

import (
   "encoding/xml"
   "os"
   "testing"
)

func Test(t *testing.T) {
   data, err := os.ReadFile("../../testdata/pluto.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var newMpd Mpd
   err = xml.Unmarshal(data, &newMpd)
   if err != nil {
      t.Fatal(err)
   }
}

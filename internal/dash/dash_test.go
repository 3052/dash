package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func Test(t *testing.T) {
   data, err := os.ReadFile("../../testdata/pluto.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var mpdVar Mpd
   err = xml.Unmarshal(data, &mpdVar)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", mpdVar)
}

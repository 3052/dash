package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestPull(t *testing.T) {
   data, err := os.ReadFile("testdata/amc.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = xml.Unmarshal(data, &media)
   if err != nil {
      t.Fatal(err)
   }
   for _, p := range media.Period {
      represent, ok := p.hd_pull()
      fmt.Printf("%+v %v\n", represent, ok)
   }
}

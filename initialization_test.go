package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestInitialization(t *testing.T) {
   data, err := os.ReadFile("testdata/amc.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = xml.Unmarshal(data, &media)
   if err != nil {
      t.Fatal(err)
   }
   represent, _ := media.hd()
   initial := represent.SegmentTemplate.Initialization
   fmt.Printf("%q\n", initial("HELLO"))
   fmt.Printf("%q\n", initial(represent.Id))
}

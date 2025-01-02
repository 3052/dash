package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestInitialization(t *testing.T) {
   data, err := os.ReadFile("testdata/cineMember.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var present Mpd
   err = xml.Unmarshal(data, &present)
   if err != nil {
      t.Fatal(err)
   }
   represent, _ := present.hd()
   template := represent.SegmentTemplate.Initialization
   initial, err := template.Url(&Representation{Id: "HELLO"})
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", initial)
   initial, err = template.Url(represent)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", initial)
}

package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestSegment(t *testing.T) {
   data, err := os.ReadFile("testdata/cineMember.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = xml.Unmarshal(data, &media)
   if err != nil {
      t.Fatal(err)
   }
   represent, _ := media.hd()
   for segment := range represent.segment() {
      fmt.Println(segment)
   }
}

func TestRepresentRepresent(t *testing.T) {
   data, err := os.ReadFile("testdata/paramount.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = xml.Unmarshal(data, &media)
   if err != nil {
      t.Fatal(err)
   }
   represent, _ := media.hd()
   for represent := range represent.representation() {
      fmt.Print(&represent, "\n\n")
   }
}

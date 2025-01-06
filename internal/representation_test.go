package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestRepresentation(t *testing.T) {
   data, err := os.ReadFile("../testdata/paramount.mpd")
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
   for represent2 := range represent.Representation() {
      fmt.Print(&represent2, "\n\n")
   }
}

package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestMpdRepresent(t *testing.T) {
   data, err := os.ReadFile("testdata/paramount.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = xml.Unmarshal(data, &media)
   if err != nil {
      t.Fatal(err)
   }
   for represent := range media.representation() {
      fmt.Print(&represent, "\n\n")
   }
}

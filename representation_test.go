package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestPluto(t *testing.T) {
   data, err := os.ReadFile("ignore/pluto.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = xml.Unmarshal(data, &media)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(media)
}

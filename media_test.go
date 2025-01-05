package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestMedia(t *testing.T) {
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
   template := represent.SegmentTemplate.Media
   media_url, err := template.Url(&Representation{Id: "HELLO"}, 999)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", media_url)
   media_url, err = template.Url(represent, 999)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", media_url)
}

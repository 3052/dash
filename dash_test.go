package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestHd(t *testing.T) {
   data, err := os.ReadFile("testdata/amc.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = xml.Unmarshal(data, &media)
   if err != nil {
      t.Fatal(err)
   }
   represent, ok := media.hd()
   fmt.Printf("%+v %v\n", represent, ok)
}

func (m Mpd) hd() (*Representation, bool) {
   for represent := range m.representation() {
      if represent.Height != nil {
         if *represent.Height > 576 {
            return &represent, true
         }
      }
   }
   return nil, false
}

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

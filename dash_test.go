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
   var media Mpd
   err = xml.Unmarshal(data, &media)
   if err != nil {
      t.Fatal(err)
   }
   represent, _ := media.hd()
   initial, ok := represent.initialization()
   fmt.Println(&initial.Url, ok)
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

func TestRepresentation(t *testing.T) {
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

func TestSeq(t *testing.T) {
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
   for represent := range represent.seq() {
      fmt.Print(&represent, "\n\n")
   }
}

package dash

import (
   "encoding/xml"
   "os"
   "testing"
)

func TestRepresentation(t *testing.T) {
   t.Run("itv", func(t *testing.T) {
      data, err := os.ReadFile("testdata/itv.mpd")
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
      for segment := range represent.Segment() {
         if segment >= 1 {
            break
         }
      }
   })
   t.Run("pluto", func(t *testing.T) {
      data, err := os.ReadFile("testdata/pluto.mpd")
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
         data := represent.String()
         if data == "" {
            t.Fatal(represent)
         }
      }
      for range represent.Representation() {
         break
      }
      for segment := range represent.Segment() {
         if segment >= 9 {
            break
         }
      }
   })
}

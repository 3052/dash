package dash

import (
   "encoding/xml"
   "os"
   "testing"
)

func TestRepresentation(t *testing.T) {
   data, err := os.ReadFile("testdata/max.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = xml.Unmarshal(data, &media)
   if err != nil {
      t.Fatal(err)
   }
   t.Run("representation", func(t *testing.T) {
      var represent Representation
      for represent = range media.Representation() {
         break
      }
      for range represent.Representation() {
         break
      }
   })
   t.Run("string", func(t *testing.T) {
      for represent := range media.Representation() {
         data := represent.String()
         if data == "" {
            t.Fatal(represent)
         }
      }
   })
}

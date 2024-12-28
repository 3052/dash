package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func (p *Period) hd() (*Representation, bool) {
   for value := range p.representation() {
      if value.Height > 576 {
         return &value, true
      }
   }
   return nil, false
}

func TestPush(t *testing.T) {
   data, err := os.ReadFile("testdata/amc.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = xml.Unmarshal(data, &media)
   if err != nil {
      t.Fatal(err)
   }
   for _, p := range media.Period {
      represent, ok := p.hd()
      fmt.Printf("%+v %v\n", represent, ok)
   }
}

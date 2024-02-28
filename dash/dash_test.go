package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestRange(t *testing.T) {
   media, err := reader("mpd/hulu.mpd")
   if err != nil {
      t.Fatal(err)
   }
   media.Visit(func(p Pointer) {
      sb := p.Representation.SegmentBase
      start, end, err := sb.Initialization.Range.Scan()
      fmt.Printf("%v %v %v ", start, end, err)
      start, end, err = sb.IndexRange.Scan()
      fmt.Printf("%v %v %v\n", start, end, err)
   })
}

func reader(name string) (*MPD, error) {
   text, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   media := new(MPD)
   if err := xml.Unmarshal(text, &media); err != nil {
      return nil, err
   }
   return media, nil
}

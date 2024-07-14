package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "slices"
   "testing"
)

func TestDownload(t *testing.T) {
   text, err := os.ReadFile("../../testdata/paramount.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = xml.Unmarshal(text, &media)
   if err != nil {
      t.Fatal(err)
   }
   for _, p := range media.Period {
      for _, adapt := range p.AdaptationSet {
         adapt.period = &p
         for _, represent := range adapt.Representation {
            if represent.Id == "thumb_160x90" {
               represent.adaptation_set = &adapt
               fmt.Printf("%+v\n", adapt.SegmentTemplate)
               represent.Media()
               fmt.Printf("%+v\n\n", adapt.SegmentTemplate)
            }
         }
      }
   }
}

func TestPrint(t *testing.T) {
   text, err := os.ReadFile("../../testdata/paramount.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = xml.Unmarshal(text, &media)
   if err != nil {
      t.Fatal(err)
   }
   var represents []Representation
   for _, p := range media.Period {
      if p.Id == "0" {
         for _, adapt := range p.AdaptationSet {
            represents = append(represents, adapt.Representation...)
         }
      }
   }
   slices.SortFunc(represents, func(a, b Representation) int {
      return int(a.Bandwidth - b.Bandwidth)
   })
   for _, represent := range represents {
      fmt.Printf("%+v\n\n", represent)
   }
}

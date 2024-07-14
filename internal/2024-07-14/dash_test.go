package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "slices"
   "testing"
)

func (r Representation) get_media() {
   for _, media := range r.Media() {
      fmt.Printf("%q\n", media)
   }
   fmt.Println()
}

func (r Representation) get_initial() {
   if v, ok := r.Initialization(); ok {
      fmt.Printf("%q\n\n", v)
   }
}

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
   var initial bool
   for _, p := range media.Period {
      for _, adapt := range p.AdaptationSet {
         for _, represent := range adapt.Representation {
            if represent.Id == "0" {
               if !initial {
                  represent.get_initial()
                  initial = true
               }
               represent.get_media()
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

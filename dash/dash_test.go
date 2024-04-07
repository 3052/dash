package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "slices"
   "testing"
)

func TestMpd(t *testing.T) {
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var m mpd
      if err := xml.Unmarshal(text, &m); err != nil {
         t.Fatal(err)
      }
      if m.MediaPresentationDuration == "" {
         t.Fatal("MediaPresentationDuration", test)
      }
      if len(m.Period) == 0 {
         t.Fatal("Period", test)
      }
   }
}

func TestRepresentation(t *testing.T) {
   sets := struct{
      bandwidth set
      base_url set
      height set
      id set
      indexRange set
      initialization set
      segmentBase set
      width set
   }{
      make(set),
      make(set),
      make(set),
      make(set),
      make(set),
      make(set),
      make(set),
      make(set),
   }
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media mpd
      if err := xml.Unmarshal(text, &media); err != nil {
         t.Fatal(err)
      }
      for _, per := range media.Period {
         for _, ada := range per.AdaptationSet {
            for _, rep := range ada.Representation {
               if rep.Bandwidth >= 1 {
                  sets.bandwidth[1] = struct{}{}
               } else {
                  sets.bandwidth[0] = struct{}{}
               }
               if rep.BaseURL != nil {
                  sets.base_url[1] = struct{}{}
               } else {
                  sets.base_url[0] = struct{}{}
               }
               if rep.Height != nil {
                  sets.height[1] = struct{}{}
               } else {
                  sets.height[0] = struct{}{}
               }
               if rep.Width != nil {
                  sets.width[1] = struct{}{}
               } else {
                  sets.width[0] = struct{}{}
               }
               if rep.ID != "" {
                  sets.id[1] = struct{}{}
               } else {
                  sets.id[0] = struct{}{}
               }
               if v := rep.SegmentBase; v != nil {
                  sets.segmentBase[1] = struct{}{}
                  if v.Initialization.Range != "" {
                     sets.initialization[1] = struct{}{}
                  } else {
                     sets.initialization[0] = struct{}{}
                  }
                  if v.IndexRange != "" {
                     sets.indexRange[1] = struct{}{}
                  } else {
                     sets.indexRange[0] = struct{}{}
                  }
               } else {
                  sets.segmentBase[0] = struct{}{}
               }
            }
         }
      }
   }
   fmt.Printf("%+v\n", sets)
}

func TestPeriod(t *testing.T) {
   duration := make(set)
   for i, test := range tests {
      if i >= 1 {
         fmt.Println()
      }
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media mpd
      if err := xml.Unmarshal(text, &media); err != nil {
         t.Fatal(err)
      }
      for _, per := range media.Period {
         if len(per.AdaptationSet) == 0 {
            t.Fatal("AdaptationSet", test)
         }
         if per.Duration != nil {
            duration[1] = struct{}{}
         } else {
            duration[0] = struct{}{}
         }
      }
   }
   fmt.Println(duration)
}

func TestAdaptation(t *testing.T) {
   sets := struct{
      codecs set
      lang set
      mimeType set
      role set
      segmentTemplate set
      value set
   }{
      make(set),
      make(set),
      make(set),
      make(set),
      make(set),
      make(set),
   }
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media mpd
      if err := xml.Unmarshal(text, &media); err != nil {
         t.Fatal(err)
      }
      for _, per := range media.Period {
         for _, ada := range per.AdaptationSet {
            if len(ada.Representation) == 0 {
               t.Fatal("Representation")
            }
            if ada.Lang != nil {
               sets.lang[1] = struct{}{}
            } else {
               sets.lang[0] = struct{}{}
            }
            if ada.Role != nil {
               sets.role[1] = struct{}{}
               if ada.Role.Value != "" {
                  sets.value[1] = struct{}{}
               } else {
                  sets.value[0] = struct{}{}
               }
            } else {
               sets.role[0] = struct{}{}
            }
            for _, rep := range ada.Representation {
               var value byte
               if ada.Codecs != nil {
                  value += 10
               }
               if rep.Codecs != nil {
                  value++
               }
               sets.codecs[value] = struct{}{}
               value = 0
               if ada.MimeType != nil {
                  value += 10
               }
               if rep.MimeType != nil {
                  value++
               }
               sets.mimeType[value] = struct{}{}
               value = 0
               if ada.SegmentTemplate != nil {
                  value += 10
               }
               if rep.SegmentTemplate != nil {
                  value++
               }
               sets.segmentTemplate[value] = struct{}{}
            }
         }
      }
   }
   fmt.Printf("%+v\n", sets)
}

func TestDelete(t *testing.T) {
   for i, test := range tests {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(test)
      fmt.Println("------------------------------------------------------------")
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      reps, err := Unmarshal(text)
      if err != nil {
         t.Fatal(err)
      }
      reps = slices.DeleteFunc(reps, func(r Representation) bool {
         if _, ok := r.Ext(); !ok {
            return true
         }
         if v, _ := r.GetAdaptationSet().GetPeriod().Seconds(); v < 9 {
            return true
         }
         return false
      })
      for i, rep := range reps {
         if i >= 1 {
            fmt.Println()
         }
         fmt.Println(rep)
      }
   }
}

package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

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
      var media MPD
      err = xml.Unmarshal(text, &media)
      if err != nil {
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
      var media MPD
      err = media.Unmarshal(text)
      if err != nil {
         t.Fatal(err)
      }
      var line bool
      for _, v := range media.Period {
         seconds, err := v.Seconds()
         if err != nil {
            t.Fatal(err)
         }
         for _, v := range v.AdaptationSet {
            for _, v := range v.Representation {
               if seconds > 9 {
                  if _, ok := v.Ext(); ok {
                     if line {
                        fmt.Println()
                     } else {
                        line = true
                     }
                     fmt.Println(v)
                  }
               }
            }
         }
      }
   }
}

func TestMpd(t *testing.T) {
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media MPD
      err = media.Unmarshal(text)
      if err != nil {
         t.Fatal(err)
      }
      if media.MediaPresentationDuration == "" {
         t.Fatal("MediaPresentationDuration", test)
      }
      if len(media.Period) == 0 {
         t.Fatal("Period", test)
      }
      for _, v := range media.Period {
         if v.mpd == nil {
            t.Fatal("mpd")
         }
         for _, v := range v.AdaptationSet {
            if v.period == nil {
               t.Fatal("period")
            }
            for _, v := range v.Representation {
               if v.adaptation_set == nil {
                  t.Fatal("adaptation_set")
               }
            }
         }
      }
   }
}

func TestPeriod(t *testing.T) {
   duration := make(set)
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media MPD
      err = xml.Unmarshal(text, &media)
      if err != nil {
         t.Fatal(err)
      }
      for _, p := range media.Period {
         if len(p.AdaptationSet) == 0 {
            t.Fatal("AdaptationSet", test)
         }
         if p.Duration != nil {
            duration[1] = struct{}{}
         } else {
            duration[0] = struct{}{}
         }
      }
   }
   fmt.Println(duration)
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
      var media MPD
      err = xml.Unmarshal(text, &media)
      if err != nil {
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

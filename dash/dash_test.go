package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestMedia(t *testing.T) {
   for i, test := range tests {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(test)
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      reps, err := Unmarshal(text)
      if err != nil {
         t.Fatal(err)
      }
      for _, rep := range reps {
         if v, ok := rep.GetSegmentTemplate(); ok {
            media := v.GetMedia(rep.ID)
            length := len(media)
            if length >= 1 {
               fmt.Println(media[length-1])
            }
         }
      }
   }
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

func TestMpd(t *testing.T) {
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media mpd
      xml.Unmarshal(text, &media)
      if media.MediaPresentationDuration == "" {
         t.Fatal("MediaPresentationDuration", test)
      }
      if len(media.Period) == 0 {
         t.Fatal("Period", test)
      }
   }
}

func TestSegmentTemplate(t *testing.T) {
   sets := struct{
      d set
      r set
      initialization set
      media set
      segmentTimeline set
   }{
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
      reps, err := Unmarshal(text)
      if err != nil {
         t.Fatal(err)
      }
      for _, rep := range reps {
         if v, ok := rep.GetSegmentTemplate(); ok {
            if v.Initialization != nil {
               sets.initialization[1] = struct{}{}
            } else {
               sets.initialization[0] = struct{}{}
            }
            if v.Media != "" {
               sets.media[1] = struct{}{}
            } else {
               sets.media[0] = struct{}{}
            }
            if v := v.SegmentTimeline; v != nil {
               sets.segmentTimeline[1] = struct{}{}
               for _, v := range v.S {
                  if v.D >= 1 {
                     sets.d[1] = struct{}{}
                  } else {
                     sets.d[0] = struct{}{}
                  }
                  if v.R != nil {
                     sets.r[1] = struct{}{}
                  } else {
                     sets.r[0] = struct{}{}
                  }
               }
            } else {
               sets.segmentTimeline[0] = struct{}{}
            }
         }
      }
   }
   fmt.Printf("%+v\n", sets)
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

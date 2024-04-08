package dash

import (
   "encoding/xml"
   "fmt"
   "io"
   "net/http"
   "os"
   "testing"
)

func TestMpd(t *testing.T) {
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media MPD
      if err := media.Unmarshal(text); err != nil {
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

func TestSegmentTemplate(t *testing.T) {
   sets := struct{
      d set
      duration set
      r set
      initialization set
      media set
      segmentTimeline set
      timescale set
   }{
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
      if err := media.Unmarshal(text); err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            for _, v := range v.Representation {
               if v, ok := v.GetSegmentTemplate(); ok {
                  if v.Duration != nil {
                     sets.duration[1] = struct{}{}
                  } else {
                     sets.duration[0] = struct{}{}
                  }
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
                  if v.Timescale != nil {
                     sets.timescale[1] = struct{}{}
                  } else {
                     sets.timescale[0] = struct{}{}
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
      }
   }
   fmt.Printf("%+v\n", sets)
}

func TestMedia(t *testing.T) {
   base := "https://gec.stan.video/09/dash/live/1540676B-1/hd/sdr/"
   res, err := http.Get(base + "high_h264-59fcad98.mpd")
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   fmt.Println(res.Request.URL)
   text, err := io.ReadAll(res.Body)
   if err != nil {
      t.Fatal(err)
   }
   var media MPD
   if err := media.Unmarshal(text); err != nil {
      t.Fatal(err)
   }
   for _, v := range media.Period {
      for _, v := range v.AdaptationSet {
         for _, represent := range v.Representation {
            if v, ok := represent.GetSegmentTemplate(); ok {
               media, err := v.GetMedia(represent)
               if err != nil {
                  t.Fatal(err)
               }
               length := len(media)
               if length >= 1 {
                  fmt.Println(base + media[length-1])
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
      if err := xml.Unmarshal(text, &media); err != nil {
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

package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "slices"
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
      var media mpd
      xml.Unmarshal(text, &media)
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
               if ada.Codecs != "" {
                  value += 10
               }
               if rep.Codecs != "" {
                  value++
               }
               sets.codecs[value] = struct{}{}
               value = 0
               if ada.MimeType != "" {
                  value += 10
               }
               if rep.MimeType != "" {
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
      xml.Unmarshal(text, &media)
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

var media_tests = []struct{
   name string
   base string
}{
   // startNumber == nil
   {"mpd/mubi.mpd", "new-york-edge2.mubicdn.net/stream/43cac9f0138aaa566a429be4542ff21c/65df1dc5/728eb9fc/mubi-films/325455/passages_eng_zxx_1800x1080_50000_mezz40828/ae8c88ed4e/drm_playlist.0ff148ef80.ism/default/"},
   // startNumber == 0
   {"mpd/amc.mpd", ""},
   // startNumber == 1
   {"mpd/paramount.mpd", "vod-gcs-cedexis.cbsaavideo.com/intl_vms/2022/02/24/2006197315671/77016_cenc_dash/"},
}

func TestMedia(t *testing.T) {
   for _, test := range media_tests {
      fmt.Println(test.name + ":")
      text, err := os.ReadFile(test.name)
      if err != nil {
         t.Fatal(err)
      }
      reps, err := Unmarshal(text)
      if err != nil {
         t.Fatal(err)
      }
      for _, media := range reps[0].Media() {
         fmt.Println(test.base + media)
      }
   }
}

func TestDelete(t *testing.T) {
   for i, name := range tests {
      if i >= 1 {
         fmt.Println()
      }
      reps = slices.DeleteFunc(reps, func(r Representation) bool {
         if _, ok := r.Ext(); !ok {
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

func TestSegmentTemplate(t *testing.T) {
   sets := struct{
      initialization set
   }{
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
      for _, rep := range ada.Representation {
         if v, ok := rep.GetSegmentTemplate(); ok {
            if v.Initialization != "" {
               sets.initialization[1] = struct{}{}
            } else {
               sets.initialization[0] = struct{}{}
            }
         }
      }
   }
   fmt.Printf("%+v\n", sets)
}

var tests = []string{
   "mpd/amc.mpd",
   "mpd/hulu.mpd",
   "mpd/mubi.mpd",
   "mpd/nbc.mpd",
   "mpd/paramount.mpd",
   "mpd/peacock.mpd",
   "mpd/roku-clear.mpd",
   "mpd/roku-protected.mpd",
   "mpd/stan.mpd",
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

func TestPeriod(t *testing.T) {
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
         if rep.adaptation_set.period.mpd == nil {
            t.Fatal("mpd", test)
         }
      }
      var media mpd
      xml.Unmarshal(text, &media)
      for _, p := range media.Period {
         if len(p.AdaptationSet) == 0 {
            t.Fatal("AdaptationSet", test)
         }
      }
   }
}

type set map[byte]struct{}


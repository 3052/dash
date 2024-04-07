package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

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
      if err := media.Unmarshal(text); err != nil {
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

type set map[byte]struct{}

var tests = []string{
   "mpd/amc.mpd",
   "mpd/hulu.mpd",
   "mpd/mubi.mpd",
   "mpd/nbc.mpd",
   "mpd/paramount.mpd",
   "mpd/peacock.mpd",
   "mpd/plex.mpd",
   "mpd/roku.mpd",
   "mpd/stan.mpd",
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

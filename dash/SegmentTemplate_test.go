package dash

import (
   "fmt"
   "os"
   "testing"
)

func TestMedia(t *testing.T) {
   for _, test := range tests {
      fmt.Println(test)
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media MPD
      err = media.Unmarshal(text)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            for _, represent := range v.Representation {
               if v, ok := represent.GetSegmentTemplate(); ok {
                  initial, ok := v.GetInitialization(represent)
                  fmt.Printf("%q %v\n", initial, ok)
                  media, err := v.GetMedia(represent)
                  if err != nil {
                     t.Fatal(err)
                  }
                  for i := range min(len(media), 9) {
                     fmt.Println(media[i])
                  }
               }
            }
         }
      }
   }
}

var tests = []string{
   "mpd/amc.mpd",
   "mpd/cine-member.mpd",
   "mpd/hulu.mpd",
   "mpd/mubi.mpd",
   "mpd/nbc.mpd",
   "mpd/paramount.mpd",
   "mpd/plex.mpd",
   "mpd/pluto.mpd",
   "mpd/roku.mpd",
   "mpd/stan.mpd",
   "mpd/tubi.mpd",
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
      err = media.Unmarshal(text)
      if err != nil {
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

type set map[byte]struct{}

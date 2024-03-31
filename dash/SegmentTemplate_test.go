package dash

import (
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
            media, err := v.GetMedia(rep)
            if err != nil {
               t.Fatal(err)
            }
            length := len(media)
            if length >= 1 {
               fmt.Println(media[length-1])
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
      reps, err := Unmarshal(text)
      if err != nil {
         t.Fatal(err)
      }
      for _, rep := range reps {
         if v, ok := rep.GetSegmentTemplate(); ok {
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
   fmt.Printf("%+v\n", sets)
}

type set map[byte]struct{}

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

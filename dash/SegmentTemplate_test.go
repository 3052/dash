package dash

import (
   "fmt"
   "io"
   "net/http"
   "os"
   "testing"
)

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
               fmt.Println(v.GetInitialization(represent))
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
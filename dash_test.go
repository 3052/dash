package dash

import (
   "fmt"
   "os"
   "testing"
)

func TestBandwidth(t *testing.T) {
   for _, test := range tests {
      fmt.Print("\n", test, "\n")
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            for _, v := range v.Representation {
               fmt.Println(v.Bandwidth)
            }
         }
      }
   }
}

func TestCodecs(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            fmt.Printf("AdaptationSet %q\n", v.Codecs)
            for _, v := range v.Representation {
               fmt.Printf("Representation %q\n", v.Codecs)
            }
         }
         fmt.Println()
      }
   }
}

func TestMedia(t *testing.T) {
   text, err := os.ReadFile("testdata/paramount.mpd")
   if err != nil {
      t.Fatal(err)
   }
   represents, err := Unmarshal(text, nil)
   if err != nil {
      t.Fatal(err)
   }
   for _, represent := range represents {
      media := represent.Media()
      fmt.Println(represent.Id)
      fmt.Println(media[0])
      fmt.Println(media[len(media)-1])
   }
}

func TestMimeType(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            fmt.Printf("AdaptationSet %q\n", v.MimeType)
            for _, v := range v.Representation {
               fmt.Printf("Representation %q\n", v.MimeType)
            }
         }
      }
   }
}

func TestRepresentation(t *testing.T) {
   text, err := os.ReadFile("testdata/max.mpd")
   if err != nil {
      t.Fatal(err)
   }
   reps, err := Unmarshal(text, nil)
   if err != nil {
      t.Fatal(err)
   }
   for _, rep := range reps {
      fmt.Print(rep, "\n\n")
   }
}

func TestSegmentBase(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            for _, v := range v.Representation {
               fmt.Printf("%+v\n", v.SegmentBase)
            }
         }
      }
   }
}

func TestSegmentTemplate(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            fmt.Printf("AdaptationSet %+v\n", v.SegmentTemplate)
            for _, v := range v.Representation {
               fmt.Printf("Representation %+v\n", v.SegmentTemplate)
            }
         }
      }
   }
}

func TestBaseUrl(t *testing.T) {
   for _, test := range tests {
      fmt.Println(test)
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      reps, err := Unmarshal(text, nil)
      if err != nil {
         t.Fatal(err)
      }
      for _, rep := range reps {
         if v, ok := rep.GetBaseUrl(); ok {
            fmt.Println(v)
         }
      }
      fmt.Println()
   }
}

func TestDuration(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("MPD %q\n", media.MediaPresentationDuration)
      for _, v := range media.Period {
         fmt.Printf("Period %q\n", v.Duration)
      }
      fmt.Println()
   }
}

func TestPeriodId(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         fmt.Printf("%q\n", v.Id)
      }
   }
}

func TestLang(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            fmt.Printf("%q\n", v.Lang)
         }
      }
   }
}

func TestPssh(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            for _, v := range v.ContentProtection {
               fmt.Printf("%q\n", v.Pssh)
            }
            for _, v := range v.Representation {
               for _, v := range v.ContentProtection {
                  fmt.Printf("%q\n", v.Pssh)
               }
            }
            fmt.Println()
         }
      }
   }
}

func TestRole(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            fmt.Println(v.Role)
         }
      }
   }
}

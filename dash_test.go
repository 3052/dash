package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

var tests = []string{
   "testdata/amc.mpd",
   "testdata/cine-member.mpd",
   "testdata/criterion.mpd",
   "testdata/ctv.mpd",
   "testdata/draken.mpd",
   "testdata/hulu.mpd",
   "testdata/max.mpd",
   "testdata/mubi.mpd",
   "testdata/nbc.mpd",
   "testdata/paramount.mpd",
   "testdata/plex.mpd",
   "testdata/pluto.mpd",
   "testdata/rakuten.mpd",
   "testdata/roku.mpd",
   "testdata/stan.mpd",
   "testdata/tubi.mpd",
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

func TestId(t *testing.T) {
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

func TestInitialization(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            if v := v.SegmentTemplate; v != nil {
               fmt.Printf("%q\n", v.Initialization)
            }
            for _, v := range v.Representation {
               if v := v.SegmentTemplate; v != nil {
                  fmt.Printf("%q\n", v.Initialization)
               }
            }
         }
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

func TestRepresentation(t *testing.T) {
   for _, test := range tests {
      fmt.Print(test, ":\n\n")
      text, err := os.ReadFile(test)
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

func new_mpd(name string) (*Mpd, error) {
   text, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   media := new(Mpd)
   err = xml.Unmarshal(text, media)
   if err != nil {
      return nil, err
   }
   return media, nil
}

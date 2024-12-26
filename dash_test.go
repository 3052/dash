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
   "testdata/tubi.mpd",
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

func new_mpd(name string) (*Mpd, error) {
   text, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   media := &Mpd{}
   err = xml.Unmarshal(text, media)
   if err != nil {
      return nil, err
   }
   return media, nil
}

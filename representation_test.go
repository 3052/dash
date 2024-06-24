package dash

import (
   "fmt"
   "testing"
)

func TestBaseUrl(t *testing.T) {
   for _, test := range tests {
      media, err := new_mpd(test)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(media.BaseUrl)
      for _, v := range media.Period {
         for _, v := range v.AdaptationSet {
            for _, v := range v.Representation {
               fmt.Println(v.BaseUrl)
            }
         }
      }
      fmt.Println()
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

func TestHulu(t *testing.T) {
   media, err := new_mpd("testdata/hulu.mpd")
   if err != nil {
      t.Fatal(err)
   }
   for _, v := range media.Period {
      for _, v := range v.AdaptationSet {
         for _, v := range v.Representation {
            fmt.Print(v, "\n\n")
         }
      }
   }
}

func TestMax(t *testing.T) {
   media, err := new_mpd("testdata/max.mpd")
   if err != nil {
      t.Fatal(err)
   }
   set := map[string]struct{}{}
   for _, v := range media.Period {
      for _, v := range v.AdaptationSet {
         for _, v := range v.Representation {
            if _, ok := set[v.Id]; !ok {
               fmt.Print(v, "\n\n")
               set[v.Id] = struct{}{}
            }
         }
      }
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

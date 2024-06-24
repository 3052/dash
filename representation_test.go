package dash

import (
   "fmt"
   "testing"
)

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
            fmt.Println()
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

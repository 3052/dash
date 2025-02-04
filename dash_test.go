package dash

import (
   "net/url"
   "os"
   "testing"
)

func TestRepresentation(t *testing.T) {
   t.Run("itv", func(t *testing.T) {
      data, err := os.ReadFile("testdata/itv.mpd")
      if err != nil {
         t.Fatal(err)
      }
      var media Mpd
      err = media.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
      var represent Representation
      for represent = range media.Representation() {
         break
      }
      for segment := range represent.Segment() {
         if segment >= 1 {
            break
         }
      }
   })
   t.Run("pluto", func(t *testing.T) {
      data, err := os.ReadFile("testdata/pluto.mpd")
      if err != nil {
         t.Fatal(err)
      }
      var media Mpd
      err = media.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
      var represent Representation
      for represent = range media.Representation() {
         data := represent.String()
         if data == "" {
            t.Fatal(represent)
         }
      }
      for range represent.Representation() {
         break
      }
      for segment := range represent.Segment() {
         if segment >= 9 {
            break
         }
      }
   })
}
func TestListUrl(t *testing.T) {
   data, err := os.ReadFile("testdata/criterion.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = media.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   var represent Representation
   for represent = range media.Representation() {
      if *represent.MimeType == "video/mp4" {
         break
      }
   }
   _, err = represent.SegmentList.SegmentUrl[0].Media.Url(&represent)
   if err != nil {
      t.Fatal(err)
   }
   _, err = ListUrl{}.Url(&Representation{})
   if err != nil {
      t.Fatal(err)
   }
}

func TestUrl(t *testing.T) {
   data, err := os.ReadFile("testdata/criterion.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   media.BaseUrl = &Url{&url.URL{
      Path: "/0/1/2/3/4/5/6/7/8/9/10",
   }}
   err = media.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   if media.BaseUrl.UnmarshalText([]byte{'\n'}) == nil {
      t.Fatal("BaseUrl.UnmarshalText")
   }
}

var range_tests = []struct {
   in  string
   out string
   ok  bool
}{
   {"!-3", "", false},
   {"-", "", false},
   {"-3", "0-3", true},
   {"2-", "2-", true},
   {"2-3", "2-3", true},
}

func TestDuration(t *testing.T) {
   var d Duration
   if d.UnmarshalText(nil) == nil {
      t.Fatal("Duration.UnmarshalText")
   }
}

func TestMpd(t *testing.T) {
   data, err := os.ReadFile("testdata/pluto.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = media.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
}

func TestPeriod(t *testing.T) {
   data, err := os.ReadFile("testdata/max.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = media.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   var represent Representation
   for represent = range media.Representation() {
      if represent.Id == "images_1" {
         break
      }
   }
   for segment := range represent.Segment() {
      if segment >= 1 {
         break
      }
   }
}

func TestRange(t *testing.T) {
   for _, test := range range_tests {
      var r Range
      ok := r.UnmarshalText([]byte(test.in)) == nil
      if ok != test.ok {
         t.Fatal("Range.UnmarshalText")
      }
      if ok {
         out, err := r.MarshalText()
         if err != nil {
            t.Fatal(err)
         }
         if string(out) != test.out {
            t.Fatalf("Range.MarshalText")
         }
      }
   }
}

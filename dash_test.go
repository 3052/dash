package dash

import (
   "41.neocities.org/dash/url"
   "os"
   "testing"
)

func TestUrl(t *testing.T) {
   data, err := os.ReadFile("testdata/criterion.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   media.BaseUrl, err = url.Parse("/0/1/2/3/4/5/6/7/8/9/10")
   if err != nil {
      t.Fatal(err)
   }
   err = media.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   if media.BaseUrl.UnmarshalText([]byte{'\n'}) == nil {
      t.Fatal("BaseUrl.UnmarshalText")
   }
}

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

func TestDuration(t *testing.T) {
   var d Duration
   if d.UnmarshalText(nil) == nil {
      t.Fatal("Duration.UnmarshalText")
   }
}

func TestPssh(t *testing.T) {
   var p Pssh
   if p.UnmarshalText([]byte{0}) == nil {
      t.Fatal("Pssh.UnmarshalText")
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

func TestSchemeIdUri(t *testing.T) {
   data, err := os.ReadFile("testdata/pluto.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = media.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   for represent := range media.Representation() {
      for _, protect := range represent.ContentProtection {
         if protect.SchemeIdUri.Widevine() {
            return
         }
      }
   }
   t.Fatal("SchemeIdUri")
}

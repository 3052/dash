package dash

import (
   "encoding/xml"
   "net/url"
   "os"
   "testing"
)

func TestDuration(t *testing.T) {
   var d Duration
   if d.UnmarshalText(nil) == nil {
      t.Fatal("Duration.UnmarshalText")
   }
}

func TestInitialization(t *testing.T) {
   data, err := os.ReadFile("ignore/pluto.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = xml.Unmarshal(data, &media)
   if err != nil {
      t.Fatal(err)
   }
   base, err := url.Parse(pluto.mpd)
   if err != nil {
      t.Fatal(err)
   }
   media.Set(base)
   var represent Representation
   for represent = range media.Representation() {
      if *represent.MimeType == "video/mp4" {
         break
      }
   }
   t.Run("Url", func(t *testing.T) {
      initial, err := represent.SegmentTemplate.Initialization.Url(&represent)
      if err != nil {
         t.Fatal(err)
      }
      if initial.String() != pluto.init {
         t.Fatal(initial)
      }
      _, err = Initialization{"\n"}.Url(&Representation{})
      if err == nil {
         t.Fatal("Initialization.Url")
      }
   })
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
   data, err := os.ReadFile("ignore/pluto.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   err = xml.Unmarshal(data, &media)
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

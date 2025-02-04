package dash

import (
   "net/url"
   "os"
   "testing"
)

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

var initial_test = struct {
   init  string
   media string
   mpd   string
}{
   mpd:   "http://silo-hybrik.pluto.tv.s3.amazonaws.com/576_pluto/clip/64ff3987cecd3f001332df52_Memento/720pDRM/20230911_090007/dash/0-end/main.mpd",
   init:  "http://silo-hybrik.pluto.tv.s3.amazonaws.com/576_pluto/clip/64ff3987cecd3f001332df52_Memento/720pDRM/20230911_090007/dash/0-end/video/240p-300/init.mp4",
   media: "http://silo-hybrik.pluto.tv.s3.amazonaws.com/576_pluto/clip/64ff3987cecd3f001332df52_Memento/720pDRM/20230911_090007/dash/0-end/video/240p-300/01362.m4s",
}

func TestDuration(t *testing.T) {
   var d Duration
   if d.UnmarshalText(nil) == nil {
      t.Fatal("Duration.UnmarshalText")
   }
}

func TestInitialization(t *testing.T) {
   data, err := os.ReadFile("testdata/pluto.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   media.BaseUrl = &Url{}
   media.BaseUrl.Url, err = url.Parse(initial_test.mpd)
   if err != nil {
      t.Fatal(err)
   }
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
   initial, err := represent.SegmentTemplate.Initialization.Url(&represent)
   if err != nil {
      t.Fatal(err)
   }
   if initial.String() != initial_test.init {
      t.Fatal(initial)
   }
   _, err = Initialization{"\n"}.Url(&Representation{})
   if err == nil {
      t.Fatal("Initialization.Url")
   }
}

func TestMedia(t *testing.T) {
   _, err := Media{"\n"}.Url(&Representation{}, 0)
   if err == nil {
      t.Fatal("Media")
   }
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
   _, err = represent.SegmentTemplate.Media.Url(&represent, 0)
   if err != nil {
      t.Fatal(err)
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

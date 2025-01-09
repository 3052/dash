package dash

import (
   "41.neocities.org/dash/url"
   "os"
   "testing"
)

func TestInitialization(t *testing.T) {
   data, err := os.ReadFile("testdata/pluto.mpd")
   if err != nil {
      t.Fatal(err)
   }
   var media Mpd
   media.BaseUrl, err = url.Parse(pluto.mpd)
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
   if initial.String() != pluto.init {
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

var pluto = struct {
   init  string
   media string
   mpd   string
}{
   mpd:   "http://silo-hybrik.pluto.tv.s3.amazonaws.com/576_pluto/clip/64ff3987cecd3f001332df52_Memento/720pDRM/20230911_090007/dash/0-end/main.mpd",
   init:  "http://silo-hybrik.pluto.tv.s3.amazonaws.com/576_pluto/clip/64ff3987cecd3f001332df52_Memento/720pDRM/20230911_090007/dash/0-end/video/240p-300/init.mp4",
   media: "http://silo-hybrik.pluto.tv.s3.amazonaws.com/576_pluto/clip/64ff3987cecd3f001332df52_Memento/720pDRM/20230911_090007/dash/0-end/video/240p-300/01362.m4s",
}

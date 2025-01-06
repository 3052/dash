package dash

import (
   "encoding/xml"
   "os"
   "testing"
)

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
   err = xml.Unmarshal(data, &media)
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

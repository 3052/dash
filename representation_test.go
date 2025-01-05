package dash

import (
   "encoding/xml"
   "net/url"
   "os"
   "testing"
)

func TestMedia(t *testing.T) {
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
   var media_url *url.URL
   for segment := range represent.Segment() {
      media_url, err = represent.SegmentTemplate.Media.Url(&represent, segment)
      if err != nil {
         t.Fatal(err)
      }
   }
   if media_url.String() != pluto.media {
      t.Fatal(media_url)
   }
}

var pluto = struct{
   init string
   media string
   mpd string
}{
   init: "http://silo-hybrik.pluto.tv.s3.amazonaws.com/576_pluto/clip/64ff3987cecd3f001332df52_Memento/720pDRM/20230911_090007/dash/0-end/video/240p-300/init.mp4",
   media: "http://silo-hybrik.pluto.tv.s3.amazonaws.com/576_pluto/clip/64ff3987cecd3f001332df52_Memento/720pDRM/20230911_090007/dash/0-end/video/240p-300/01362.m4s",
   mpd: "http://silo-hybrik.pluto.tv.s3.amazonaws.com/576_pluto/clip/64ff3987cecd3f001332df52_Memento/720pDRM/20230911_090007/dash/0-end/main.mpd",
}

package hls

import (
   "fmt"
   "os"
   "testing"
)

var segment_names = []string{
   "audio_eng_aacl.m3u8",
   "video.m3u8",
}

func TestSegment(t *testing.T) {
   for _, name := range segment_names {
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      var segment MediaSegment
      segment.New(string(text))
      fmt.Printf("%+v\n", segment.Key)
      for _, uri := range segment.URI {
         fmt.Printf("%q\n", uri)
      }
   }
}

package hls

import (
   "fmt"
   "strconv"
   "strings"
   "testing"
)

const media = `#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="audio",NAME="English",LANGUAGE="eng",DEFAULT=YES,AUTOSELECT=YES,URI="QualityLevels(192000)/Manifest(audio_eng_aacl,format=m3u8-aapl,filter=desktop)"`

func Test_Strconv(t *testing.T) {
   for _, field := range split(media) {
      fmt.Printf("%q\n", field)
   }
}

func Benchmark_Strconv(b *testing.B) {
   for n := 0; n < b.N; n++ {
      _ = split(media)
   }
}

func Benchmark_Strconv_Cap(b *testing.B) {
   for n := 0; n < b.N; n++ {
      _ = split_cap(media)
   }
}

func split(s string) []string {
   var field []string
   key, after, ok := strings.Cut(s, ":")
   if !ok {
      return nil
   }
   field = append(field, key)
   for {
      key, after, ok = strings.Cut(after, "=")
      if !ok {
         return field
      }
      field = append(field, key)
      value, err := strconv.QuotedPrefix(after)
      if err != nil {
         value, after, _ = strings.Cut(after, ",")
      }
      field = append(field, value)
      if err == nil {
         after = after[len(value):]
         _, after, _ = strings.Cut(after, ",")
      }
   }
}

func split_cap(s string) []string {
   field := make([]string, 0, strings.Count(s, ","))
   key, after, ok := strings.Cut(s, ":")
   if !ok {
      return nil
   }
   field = append(field, key)
   for {
      key, after, ok = strings.Cut(after, "=")
      if !ok {
         return field
      }
      field = append(field, key)
      value, err := strconv.QuotedPrefix(after)
      if err != nil {
         value, after, _ = strings.Cut(after, ",")
      }
      field = append(field, value)
      if err == nil {
         after = after[len(value):]
         _, after, _ = strings.Cut(after, ",")
      }
   }
}

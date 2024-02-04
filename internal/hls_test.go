package hls

import (
   "regexp"
   "strings"
   "testing"
   "text/scanner"
   "unicode"
)

const media = `#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="audio",NAME="English",LANGUAGE="eng",DEFAULT=YES,AUTOSELECT=YES,URI="QualityLevels(192000)/Manifest(audio_eng_aacl,format=m3u8-aapl,filter=desktop)"`

var reKeyValue = regexp.MustCompile(`([a-zA-Z0-9_-]+)=("[^"]+"|[^",]+)`)

func Benchmark_RegExp(b *testing.B) {
   for n := 0; n < b.N; n++ {
      _ = reKeyValue.FindAllStringSubmatch(media, -1)
   }
}

func Benchmark_Scanner(b *testing.B) {
   var s scanner.Scanner
   s.IsIdentRune = func(r rune, i int) bool {
      return r == '-' || unicode.IsLetter(r)
   }
   for n := 0; n < b.N; n++ {
      s.Init(strings.NewReader(media))
      for s.Scan() != scanner.EOF {
         _ = s.TokenText()
      }
   }
}

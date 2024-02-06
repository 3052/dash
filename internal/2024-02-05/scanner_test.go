package hls

import (
   "strings"
   "testing"
   "text/scanner"
   "unicode"
)

const media = `#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="audio",NAME="English",LANGUAGE="eng",DEFAULT=YES,AUTOSELECT=YES,URI="QualityLevels(192000)/Manifest(audio_eng_aacl,format=m3u8-aapl,filter=desktop)"`

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

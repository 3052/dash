package hls

import (
   "fmt"
   "strings"
   "testing"
   "text/scanner"
   "unicode"
)

func Test_Scanner(t *testing.T) {
   var s scanner.Scanner
   s.IsIdentRune = func(r rune, _ int) bool {
      return r == '-' || unicode.IsLetter(r)
   }
   s.Init(strings.NewReader(media))
   for s.Scan() != scanner.EOF {
      fmt.Printf("%q\n", s.TokenText())
   }
}

func Benchmark_Scanner(b *testing.B) {
   var s scanner.Scanner
   s.IsIdentRune = func(r rune, _ int) bool {
      return r == '-' || unicode.IsLetter(r)
   }
   for n := 0; n < b.N; n++ {
      s.Init(strings.NewReader(media))
      for s.Scan() != scanner.EOF {
         _ = s.TokenText()
      }
   }
}

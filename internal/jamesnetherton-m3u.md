# jamesnetherton-m3u

move away from RegExp

I like this package because it has no third party imports:

https://github.com/jamesnetherton/m3u/blob/45aea2b60bc42cd569432224a6b9a97cdbf8c05f/go.mod#L1-L3

but I noticed today the parsing is based on regular expression:

https://github.com/jamesnetherton/m3u/blob/45aea2b60bc42cd569432224a6b9a97cdbf8c05f/m3u.go#L74

using this file:

~~~go
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
~~~

I get this result:

~~~
goos: windows
goarch: amd64
cpu: Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz
Benchmark_RegExp-12               282069   4286 ns/op   914 B/op   15 allocs/op
Benchmark_Scanner-12              858847   1491 ns/op   224 B/op   16 allocs/op
~~~

so the RegExp option is using 4 times the memory, while the scanner option is
double to triple the speed. its possible other options are better as well, the
above implementation is just what I came up with.

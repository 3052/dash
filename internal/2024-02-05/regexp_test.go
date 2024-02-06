package hls

import (
   "regexp"
   "testing"
)

var reKeyValue = regexp.MustCompile(`([a-zA-Z0-9_-]+)=("[^"]+"|[^",]+)`)

func Benchmark_RegExp(b *testing.B) {
   for n := 0; n < b.N; n++ {
      _ = reKeyValue.FindAllStringSubmatch(media, -1)
   }
}

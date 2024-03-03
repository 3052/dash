package hls

import (
   "fmt"
   "regexp"
   "testing"
)

var reKeyValue = regexp.MustCompile(`([a-zA-Z0-9_-]+)=("[^"]+"|[^",]+)`)

func TestRegExp(t *testing.T) {
   for _, match := range reKeyValue.FindAllStringSubmatch(media, -1) {
      fmt.Printf("%q\n", match)
   }
}

func BenchmarkRegExp(b *testing.B) {
   for n := 0; n < b.N; n++ {
      _ = reKeyValue.FindAllStringSubmatch(media, -1)
   }
}

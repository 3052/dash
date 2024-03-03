# HLS

~~~
goos: windows
goarch: amd64
cpu: Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz
BenchmarkScanner-12             796146   1594 ns/op     224 B/op   16 allocs/op
BenchmarkRegExp-12              290432   4553 ns/op     915 B/op   15 allocs/op
BenchmarkStrconv-12            2053395    582.3 ns/op   496 B/op    5 allocs/op
BenchmarkStrconvCap-12         2807541    429.8 ns/op   384 B/op    2 allocs/op
~~~

## why not `encoding/csv`?

with `Reader.LazyQuotes = false`, you get:

~~~
parse error on line 1, column 34: bare " in non-quoted-field
~~~

with `Reader.LazyQuotes = true`, you get:

~~~
URI="QualityLevels(192000)/Manifest(audio_eng_aacl
format=m3u8-aapl
filter=desktop)"
~~~

## why not `regexp`?

its slower and more memory

## why not `text/scanner`?

its slower and more allocations

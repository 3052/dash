# why

## why not `strconv`?

## why not `bufio`?

## why not `fmt`?

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

its slower and more memory:

~~~
goos: windows
goarch: amd64
cpu: Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz
Benchmark_RegExp-12               282069   4286 ns/op   914 B/op   15 allocs/op
Benchmark_Scanner-12              858847   1491 ns/op   224 B/op   16 allocs/op
~~~

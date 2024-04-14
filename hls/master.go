package hls

import (
   "strconv"
   "strings"
)

const ModeLine = `
{{- range $index, $_ := . -}}
   {{ with $index }}
index = {{ . }}
   {{- else -}}
index = 0
   {{- end }}
bandwidth = {{ .Bandwidth }}
codecs = {{ .Codecs }}
   {{- with .Resolution }}
resolution = {{ . }}
   {{- end }}
{{ end -}}
`

// datatracker.ietf.org/doc/html/rfc8216#section-4.3.4
type MasterPlaylist []VariantStream

func (m MasterPlaylist) Index(i int) (*VariantStream, bool) {
   if i >= 0 {
      if i < len(m) {
         return &m[i], true
      }
   }
   return nil, false
}

func (m *MasterPlaylist) New(s string) {
   for s != "" {
      var line string
      line, s, _ = strings.Cut(s, "\r\n")
      var key string
      key, line, _ = strings.Cut(line, ":")
      if key == "#EXT-X-STREAM-INF" {
         var stream VariantStream
         for {
            var ok bool
            key, line, ok = strings.Cut(line, "=")
            if !ok {
               break
            }
            value, err := strconv.QuotedPrefix(line)
            if err != nil {
               value, line, _ = strings.Cut(line, ",")
            } else {
               line = line[len(value):]
               _, line, _ = strings.Cut(line, ",")
               // github.com/golang/go/blob/go1.22.0/src/runtime/debug/mod.go#L240-L250
               value, _ = strconv.Unquote(value)
            }
            switch key {
            case "BANDWIDTH":
               stream.Bandwidth = value
            case "CODECS":
               stream.Codecs = value
            case "RESOLUTION":
               stream.Resolution = value
            }
         }
         stream.URI, s, _ = strings.Cut(s, "\r\n")
         *m = append(*m, stream)
      }
   }
}

func (v VariantStream) Ext() string {
   if v.Resolution != "" {
      return ".m4v"
   }
   return ".m4a"
}

// datatracker.ietf.org/doc/html/rfc8216#section-4.3.4.2
type VariantStream struct {
   Bandwidth string
   Codecs string
   Resolution string
   URI string
}

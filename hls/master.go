package hls

import (
   "strconv"
   "strings"
)

// datatracker.ietf.org/doc/html/rfc8216#section-4.3.4
type MasterPlaylist []VariantStream

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

// datatracker.ietf.org/doc/html/rfc8216#section-4.3.4.2
type VariantStream struct {
   Bandwidth string
   Codecs string
   Resolution string
   URI string
}

func (v VariantStream) Ext() string {
   if v.Resolution != "" {
      return ".m4v"
   }
   return ".m4a"
}

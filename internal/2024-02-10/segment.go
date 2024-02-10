package hls

import (
   "strconv"
   "strings"
)

// datatracker.ietf.org/doc/html/rfc8216#section-3
type MediaSegment struct {
   // datatracker.ietf.org/doc/html/rfc8216#section-4.3.2.4
   Key struct {
      IV string
      URI string
   }
   URI []string
}

func (m *MediaSegment) New(s string) {
   for s != "" {
      var line string
      line, s, _ = strings.Cut(s, "\r\n")
      var key string
      key, line, _ = strings.Cut(line, ":")
      switch key {
      case "#EXT-X-KEY":
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
            case "IV":
               m.Key.IV = value
            case "URI":
               m.Key.URI = value
            }
         }
      case "#EXTINF":
         line, s, _ = strings.Cut(s, "\r\n")
         m.URI = append(m.URI, line)
      }
   }
}

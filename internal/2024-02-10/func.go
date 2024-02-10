package hls

import (
   "io"
   "strconv"
   "strings"
   "text/scanner"
   "unicode"
)

func (m *MasterPlaylist) New(s string) error {
   for s != "" {
      var line string
      line, s, _ = strings.Cut(s, "\r\n")
      var key string
      key, line, _ = strings.Cut(line, ":")
      switch key {
      // rfc-editor.org/rfc/rfc8216#section-4.3.4.1
      case "#EXT-X-MEDIA":
         for {
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
            var media MediaPlaylist
            switch key {
            case "TYPE":
               media.Type = value
            case "URI":
               media.URI = value
            }
            m.Media = append(m.Media, media)
         }
      case "#EXT-X-STREAM-INF":
         var stream VariantStream
         for {
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
            case "RESOLUTION":
               stream.Resolution = value
            }
         }
         stream.URI, s, _ = strings.Cut(s, "\r\n")
         m.Stream = append(m.Stream, stream)
      }
   }
}

func (s Scanner) Segment() (*Segment, error) {
   var seg Segment
   for s.y.Scan() != scanner.EOF {
      line := s.y.TokenText()
      var err error
      switch {
      case len(line) >= 1 && !strings.HasPrefix(line, "#"):
         seg.URI = append(seg.URI, line)
      case line == "#EXT-X-DISCONTINUITY":
         if seg.Key != "" {
            return &seg, nil
         }
      case strings.HasPrefix(line, "#EXT-X-KEY:"):
         seg.URI = nil
         s.x.Init(strings.NewReader(line))
         for s.x.Scan() != scanner.EOF {
            switch s.x.TokenText() {
            case "IV":
               s.x.Scan()
               s.x.Scan()
               seg.RawIv = s.x.TokenText()
            case "URI":
               s.x.Scan()
               s.x.Scan()
               seg.Key, err = strconv.Unquote(s.x.TokenText())
               if err != nil {
                  return nil, err
               }
            }
         }
      case strings.HasPrefix(line, "#EXT-X-MAP:"):
         s.x.Init(strings.NewReader(line))
         for s.x.Scan() != scanner.EOF {
            switch s.x.TokenText() {
            case "URI":
               s.x.Scan()
               s.x.Scan()
               seg.Map, err = strconv.Unquote(s.x.TokenText())
               if err != nil {
                  return nil, err
               }
            }
         }
      }
   }
   return &seg, nil
}

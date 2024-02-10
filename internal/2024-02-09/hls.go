package hls

import (
   "io"
   "strconv"
   "strings"
   "text/scanner"
   "unicode"
)

// rfc-editor.org/rfc/rfc8216#section-4.3.4.1
type MediaPlaylist map[string]string

func (m MediaPlaylist) URI() (string, bool) {
   return m["URI"]
}

// datatracker.ietf.org/doc/html/rfc8216#section-4.3.4.2
type VariantStream map[string]string

func (v VariantStream) URI() string {
   return v["URI"]
}

// datatracker.ietf.org/doc/html/rfc8216#section-4.3.4
type MasterPlaylist struct {
   Media []MediaPlaylist
   Stream []VariantStream
}

func (m *MasterPlaylist) New(s string) error {
   for s != "" {
      var line string
      line, s, _ = strings.Cut(s, "\r\n")
      switch {
      // rfc-editor.org/rfc/rfc8216#section-4.3.4.1
      case strings.HasPrefix(line, "#EXT-X-MEDIA:"):
         var med MediaPlaylist
         for s.x.Scan() != scanner.EOF {
            switch s.x.TokenText() {
            case "CHARACTERISTICS":
               s.x.Scan()
               s.x.Scan()
               med.Characteristics, err = strconv.Unquote(s.x.TokenText())
            case "GROUP-ID":
               s.x.Scan()
               s.x.Scan()
               med.GroupId, err = strconv.Unquote(s.x.TokenText())
            case "NAME":
               s.x.Scan()
               s.x.Scan()
               med.Name, err = strconv.Unquote(s.x.TokenText())
            case "TYPE":
               s.x.Scan()
               s.x.Scan()
               med.Type = s.x.TokenText()
            case "URI":
               s.x.Scan()
               s.x.Scan()
               med.RawUri, err = strconv.Unquote(s.x.TokenText())
            }
            if err != nil {
               return nil, err
            }
         }
         mas.Media = append(mas.Media, med)
      case strings.HasPrefix(line, "#EXT-X-STREAM-INF:"):
         var str VariantStream
         for s.x.Scan() != scanner.EOF {
            switch s.x.TokenText() {
            case "AUDIO":
               s.x.Scan()
               s.x.Scan()
               str.Audio, err = strconv.Unquote(s.x.TokenText())
            case "BANDWIDTH":
               s.x.Scan()
               s.x.Scan()
               str.Bandwidth, err = strconv.ParseInt(s.x.TokenText(), 10, 64)
            case "CODECS":
               s.x.Scan()
               s.x.Scan()
               str.Codecs, err = strconv.Unquote(s.x.TokenText())
            case "RESOLUTION":
               s.x.Scan()
               s.x.Scan()
               str.Resolution = s.x.TokenText()
            }
            if err != nil {
               return nil, err
            }
         }
         s.y.Scan()
         str.RawUri = s.y.TokenText()
         mas.Stream = append(mas.Stream, str)
      }
   }
   return &mas, nil
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

type Segment struct {
   Key string
   Map string
   RawIv string
   URI []string
}

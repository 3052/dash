package hls

import (
   "io"
   "strconv"
   "strings"
   "text/scanner"
   "unicode"
)

type Scanner struct {
   x, y scanner.Scanner
}

func (s *Scanner) New(input io.Reader) {
   s.y.Init(input)
   s.y.IsIdentRune = func(r rune, _ int) bool {
      if r == '\n' {
         return false
      }
      if r == '\r' {
         return false
      }
      return true
   }
   s.x.IsIdentRune = func(r rune, _ int) bool {
      if r == '-' {
         return true
      }
      if unicode.IsDigit(r) {
         return true
      }
      if unicode.IsLetter(r) {
         return true
      }
      return false
   }
}

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

/////////////////////////////

func (s Scanner) Master() (*MasterPlaylist, error) {
   var mas MasterPlaylist
   for s.y.Scan() != scanner.EOF {
      var err error
      line := s.y.TokenText()
      s.x.Init(strings.NewReader(line))
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
               med.Group_ID, err = strconv.Unquote(s.x.TokenText())
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
               med.Raw_URI, err = strconv.Unquote(s.x.TokenText())
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
         str.Raw_URI = s.y.TokenText()
         mas.Stream = append(mas.Stream, str)
      }
   }
   return &mas, nil
}

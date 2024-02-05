package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "encoding/hex"
   "io"
   "strconv"
   "strings"
   "text/scanner"
   "unicode"
)

type Media struct {
   Group_ID string
   Type string
   Name string
   Characteristics string
   Raw_URI string
}

func (m Media) String() string {
   var b strings.Builder
   b.WriteString("group ID: ")
   b.WriteString(m.Group_ID)
   b.WriteString("\ntype: ")
   b.WriteString(m.Type)
   b.WriteString("\nname: ")
   b.WriteString(m.Name)
   if m.Characteristics != "" {
      b.WriteString("\ncharacteristics: ")
      b.WriteString(m.Characteristics)
   }
   return b.String()
}

type Stream struct {
   Bandwidth int64
   Raw_URI string
   Audio string
   Codecs string
   Resolution string
}

func (m Stream) String() string {
   var b []byte
   if m.Resolution != "" {
      b = append(b, "resolution: "...)
      b = append(b, m.Resolution...)
      b = append(b, '\n')
   }
   b = append(b, "bandwidth: "...)
   b = strconv.AppendInt(b, m.Bandwidth, 10)
   if m.Codecs != "" {
      b = append(b, "\ncodecs: "...)
      b = append(b, m.Codecs...)
   }
   if m.Audio != "" {
      b = append(b, "\naudio: "...)
      b = append(b, m.Audio...)
   }
   return string(b)
}
type Block struct {
   cipher.Block
   key []byte
}

func New_Block(key []byte) (*Block, error) {
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   return &Block{block, key}, nil
}

func (b Block) Decrypt(text, iv []byte) []byte {
   cipher.NewCBCDecrypter(b.Block, iv).CryptBlocks(text, text)
   if len(text) >= 1 {
      pad := text[len(text)-1]
      if len(text) >= int(pad) {
         text = text[:len(text)-int(pad)]
      }
   }
   return text
}

func (b Block) Decrypt_Key(text []byte) []byte {
   return b.Decrypt(text, b.key)
}

func (s Scanner) Segment() (*Segment, error) {
   var seg Segment
   for s.line.Scan() != scanner.EOF {
      line := s.line.TokenText()
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
         s.Init(strings.NewReader(line))
         for s.Scan() != scanner.EOF {
            switch s.TokenText() {
            case "IV":
               s.Scan()
               s.Scan()
               seg.Raw_IV = s.TokenText()
            case "URI":
               s.Scan()
               s.Scan()
               seg.Key, err = strconv.Unquote(s.TokenText())
               if err != nil {
                  return nil, err
               }
            }
         }
      case strings.HasPrefix(line, "#EXT-X-MAP:"):
         s.Init(strings.NewReader(line))
         for s.Scan() != scanner.EOF {
            switch s.TokenText() {
            case "URI":
               s.Scan()
               s.Scan()
               seg.Map, err = strconv.Unquote(s.TokenText())
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
   Raw_IV string
   URI []string
}

func (s Segment) IV() ([]byte, error) {
   up := strings.ToUpper(s.Raw_IV)
   return hex.DecodeString(strings.TrimPrefix(up, "0X"))
}

func (Media) Ext() string {
   return ".m4a"
}

func (m Media) URI() string {
   return m.Raw_URI
}

func (Stream) Ext() string {
   return ".m4v"
}

func (m Stream) URI() string {
   return m.Raw_URI
}

func (s Scanner) Master() (*Master, error) {
   var mas Master
   for s.line.Scan() != scanner.EOF {
      var err error
      line := s.line.TokenText()
      s.Init(strings.NewReader(line))
      switch {
      // rfc-editor.org/rfc/rfc8216#section-4.3.4.1
      case strings.HasPrefix(line, "#EXT-X-MEDIA:"):
         var med Media
         for s.Scan() != scanner.EOF {
            switch s.TokenText() {
            case "CHARACTERISTICS":
               s.Scan()
               s.Scan()
               med.Characteristics, err = strconv.Unquote(s.TokenText())
            case "GROUP-ID":
               s.Scan()
               s.Scan()
               med.Group_ID, err = strconv.Unquote(s.TokenText())
            case "NAME":
               s.Scan()
               s.Scan()
               med.Name, err = strconv.Unquote(s.TokenText())
            case "TYPE":
               s.Scan()
               s.Scan()
               med.Type = s.TokenText()
            case "URI":
               s.Scan()
               s.Scan()
               med.Raw_URI, err = strconv.Unquote(s.TokenText())
            }
            if err != nil {
               return nil, err
            }
         }
         mas.Media = append(mas.Media, med)
      case strings.HasPrefix(line, "#EXT-X-STREAM-INF:"):
         var str Stream
         for s.Scan() != scanner.EOF {
            switch s.TokenText() {
            case "AUDIO":
               s.Scan()
               s.Scan()
               str.Audio, err = strconv.Unquote(s.TokenText())
            case "BANDWIDTH":
               s.Scan()
               s.Scan()
               str.Bandwidth, err = strconv.ParseInt(s.TokenText(), 10, 64)
            case "CODECS":
               s.Scan()
               s.Scan()
               str.Codecs, err = strconv.Unquote(s.TokenText())
            case "RESOLUTION":
               s.Scan()
               s.Scan()
               str.Resolution = s.TokenText()
            }
            if err != nil {
               return nil, err
            }
         }
         s.line.Scan()
         str.Raw_URI = s.line.TokenText()
         mas.Stream = append(mas.Stream, str)
      }
   }
   return &mas, nil
}

type Mixed interface {
   Ext() string
   URI() string
}

type Master struct {
   Media []Media
   Stream []Stream
}

type Scanner struct {
   line scanner.Scanner
   scanner.Scanner
}

func New_Scanner(body io.Reader) Scanner {
   var scan Scanner
   scan.line.Init(body)
   scan.line.IsIdentRune = func(r rune, i int) bool {
      if r == '\n' {
         return false
      }
      if r == '\r' {
         return false
      }
      if r == scanner.EOF {
         return false
      }
      return true
   }
   scan.IsIdentRune = func(r rune, i int) bool {
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
   return scan
}

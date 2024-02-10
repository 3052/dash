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
   GroupId string
   Type string
   Name string
   Characteristics string
   RawUri string
}

type Stream struct {
   Bandwidth int64
   RawUri string
   Audio string
   Codecs string
   Resolution string
}

type Block struct {
   cipher.Block
   key []byte
}

func NewBlock(key []byte) (*Block, error) {
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

func (b Block) DecryptKey(text []byte) []byte {
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
               seg.RawIv = s.TokenText()
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
   RawIv string
   URI []string
}

func (s Segment) IV() ([]byte, error) {
   up := strings.ToUpper(s.RawIv)
   return hex.DecodeString(strings.TrimPrefix(up, "0X"))
}

func (Media) Ext() string {
   return ".m4a"
}

func (Stream) Ext() string {
   return ".m4v"
}

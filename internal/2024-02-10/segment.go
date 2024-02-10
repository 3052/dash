package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "encoding/hex"
   "strconv"
   "strings"
)

func (s Segment) IV() ([]byte, error) {
   up := strings.ToUpper(s.RawIv)
   return hex.DecodeString(strings.TrimPrefix(up, "0X"))
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

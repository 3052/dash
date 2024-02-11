package hls

import (
   "crypto/aes"
   "crypto/cipher"
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
               // github.com/golang/go/blob/go1.22.0/src/runtime/debug/mod.go#L240-L250
               value, _ = strconv.Unquote(value)
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

type BlockMode struct {
   block cipher.Block
   key []byte
}

func (b *BlockMode) New(key []byte) error {
   var err error
   b.block, err = aes.NewCipher(key)
   if err != nil {
      return err
   }
   b.key = key
   return nil
}

func (b BlockMode) Decrypt(text []byte) []byte {
   cipher.NewCBCDecrypter(b.block, b.key).CryptBlocks(text, text)
   if len(text) >= 1 {
      pad := text[len(text)-1]
      if len(text) >= int(pad) {
         text = text[:len(text)-int(pad)]
      }
   }
   return text
}

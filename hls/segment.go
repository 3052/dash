package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "encoding/hex"
   "strconv"
   "strings"
)

func Decrypt(b cipher.Block, iv, text []byte) []byte {
   cipher.NewCBCDecrypter(b, iv).CryptBlocks(text, text)
   if len(text) >= 1 {
      pad := text[len(text)-1]
      if len(text) >= int(pad) {
         text = text[:len(text)-int(pad)]
      }
   }
   return text
}

func NewCipher(key []byte) (cipher.Block, error) {
   return aes.NewCipher(key)
}

// datatracker.ietf.org/doc/html/rfc8216#section-3
type MediaSegment struct {
   // datatracker.ietf.org/doc/html/rfc8216#section-4.3.2.4
   Key struct {
      Iv string
      Uri string
   }
   Uri []string
}

func (m MediaSegment) Iv() ([]byte, error) {
   return hex.DecodeString(strings.TrimPrefix(m.Key.Iv, "0X"))
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
               m.Key.Iv = value
            case "URI":
               m.Key.Uri = value
            }
         }
      case "#EXTINF":
         line, s, _ = strings.Cut(s, "\r\n")
         m.Uri = append(m.Uri, line)
      }
   }
}

package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "encoding/hex"
   "strings"
)

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

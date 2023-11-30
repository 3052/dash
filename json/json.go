package json

import (
   "bytes"
   "encoding/json"
)

func Cut(s, before, after []byte) ([]byte, []byte) {
   i := bytes.Index(s, append(before, after...))
   if i == -1 {
      return s, nil
   }
   i += len(before)
   return s[:i], s[i:]
}

func Decode(data []byte, v any) error {
   return json.NewDecoder(bytes.NewReader(data)).Decode(v)
}

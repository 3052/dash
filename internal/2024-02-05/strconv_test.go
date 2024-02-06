package hls

import (
   "fmt"
   "strconv"
   "strings"
   "testing"
)

func split(s string) ([]string, bool) {
   var field []string
   key, after, ok := strings.Cut(s, ":")
   if !ok {
      return nil, false
   }
   field = append(field, key)
   for {
      key, after, ok = strings.Cut(after, "=")
      if !ok {
         return field, true
      }
      field = append(field, key)
      value, err := strconv.QuotedPrefix(after)
      if err != nil {
         value, after, _ = strings.Cut(after, ",")
      }
      field = append(field, value)
      if err == nil {
         after = after[len(value):]
         _, after, _ = strings.Cut(after, ",")
      }
   }
}

func Test_Strconv(t *testing.T) {
   fields, ok := split(media)
   if ok {
      for _, field := range fields {
         fmt.Printf("%q\n", field)
      }
   }
}

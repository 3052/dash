package hls

import (
   "encoding/csv"
   "fmt"
   "strings"
   "testing"
)

func Test_CSV(t *testing.T) {
   {
      read := csv.NewReader(strings.NewReader(media))
      read.LazyQuotes = false
      _, err := read.Read()
      fmt.Println(err)
   }
   {
      read := csv.NewReader(strings.NewReader(media))
      read.LazyQuotes = true
      record, err := read.Read()
      if err != nil {
         t.Fatal(err)
      }
      for _, field := range record {
         fmt.Println(field)
      }
   }
}

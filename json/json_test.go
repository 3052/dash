package json

import (
   "fmt"
   "testing"
)

var dirty = []byte(`one two {"year":12,"month":31} three`)

type date struct {
   Year int
   Month int
}

func Test_After(t *testing.T) {
   _, after := Cut(dirty, nil, []byte(`{"year"`))
   var value date
   err := Decode(after, &value)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", value)
}

func Test_Both(t *testing.T) {
   _, after := Cut(dirty, []byte(" two "), []byte(`{"year"`))
   var value date
   err := Decode(after, &value)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", value)
}

func Test_Before(t *testing.T) {
   _, after := Cut(dirty, []byte(" two "), nil)
   var value date
   err := Decode(after, &value)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", value)
}

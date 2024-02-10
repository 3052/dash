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

func TestAfter(t *testing.T) {
   _, after := Cut(dirty, nil, []byte(`{"year"`))
   var value date
   err := Decode(after, &value)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", value)
}

func TestBoth(t *testing.T) {
   _, after := Cut(dirty, []byte(" two "), []byte(`{"year"`))
   var value date
   err := Decode(after, &value)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", value)
}

func TestBefore(t *testing.T) {
   _, after := Cut(dirty, []byte(" two "), nil)
   var value date
   err := Decode(after, &value)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", value)
}

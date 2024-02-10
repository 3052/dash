package xml

import (
   "fmt"
   "testing"
)

const dirty = `one two <root>
  <Year>12</Year>
  <Month>31</Month>
</root> three`

func TestBefore(t *testing.T) {
   _, after := Cut([]byte(dirty), []byte(" two "), nil)
   var root struct { Year, Month int }
   if err := Decode(after, &root); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", root)
}

func TestAfter(t *testing.T) {
   _, after := Cut([]byte(dirty), nil, []byte("<root>"))
   var root struct { Year, Month int }
   if err := Decode(after, &root); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", root)
}

func TestBoth(t *testing.T) {
   _, after := Cut([]byte(dirty), []byte(" two "), []byte("<root>"))
   var root struct { Year, Month int }
   if err := Decode(after, &root); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", root)
}

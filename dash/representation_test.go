package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
   "text/template"
)

func TestModeLine(t *testing.T) {
   tmpl, err := new(template.Template).Parse(ModeLine)
   if err != nil {
      t.Fatal(err)
   }
   for i, name := range tests {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(name)
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      var media MPD
      if err := xml.Unmarshal(text, &media); err != nil {
         t.Fatal(err)
      }
      if err := tmpl.Execute(os.Stdout, media); err != nil {
         t.Fatal(err)
      }
   }
}

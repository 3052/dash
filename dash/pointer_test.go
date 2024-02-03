package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
   "text/template"
)

func Test_Pointer(t *testing.T) {
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media MPD
      if err := xml.Unmarshal(text, &media); err != nil {
         t.Fatal(err)
      }
      media.Every(func(p Pointer) {
         fmt.Printf("name:%v period:%q ", test, p.Period.ID)
         _, ok := p.PSSH()
         fmt.Printf("pssh:%v ", ok)
         fmt.Printf("mimeType:%q\n", p.MimeType())
      })
   }
}

func Test_Template(t *testing.T) {
   tmpl, err := new(template.Template).Parse(Template)
   if err != nil {
      t.Fatal(err)
   }
   file, err := os.Create("dash.html")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   for _, name := range tests {
      file.WriteString(name)
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      var media MPD
      if err := xml.Unmarshal(text, &media); err != nil {
         t.Fatal(err)
      }
      if err := tmpl.Execute(file, media); err != nil {
         t.Fatal(err)
      }
   }
}

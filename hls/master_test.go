package hls

import (
   "fmt"
   "os"
   "testing"
)

func TestTemplate(t *testing.T) {
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

func TestMaster(t *testing.T) {
   text, err := os.ReadFile("m3u8/desktop_master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   var master MasterPlaylist
   master.New(string(text))
   for _, stream := range master {
      fmt.Printf("%+v\n", stream)
   }
}

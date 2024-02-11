package hls

import (
   "fmt"
   "os"
   "testing"
   "text/template"
)

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

func TestTemplate(t *testing.T) {
   tmpl, err := new(template.Template).Parse(Template)
   if err != nil {
      t.Fatal(err)
   }
   file, err := os.Create("ignore.html")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   text, err := os.ReadFile("m3u8/desktop_master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   var master MasterPlaylist
   master.New(string(text))
   if err := tmpl.Execute(file, master); err != nil {
      t.Fatal(err)
   }
}

package hls

import (
   "fmt"
   "os"
   "testing"
   "text/template"
)

func TestModeLine(t *testing.T) {
   line, err := new(template.Template).Parse(ModeLine)
   if err != nil {
      t.Fatal(err)
   }
   text, err := os.ReadFile("m3u8/desktop_master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   var master MasterPlaylist
   master.New(string(text))
   if err := line.Execute(os.Stdout, master); err != nil {
      t.Fatal(err)
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

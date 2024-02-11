package hls

import (
   "fmt"
   "os"
   "testing"
)

func TestMaster(t *testing.T) {
   text, err := os.ReadFile("desktop_master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   var master MasterPlaylist
   master.New(string(text))
   for _, stream := range master {
      fmt.Printf("%+v\n", stream)
   }
}

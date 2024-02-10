package hls

import (
   "fmt"
   "os"
   "testing"
)

func TestStream(t *testing.T) {
   for _, name := range master_tests {
      text, err := reverse(name)
      if err != nil {
         t.Fatal(err)
      }
      master, err := NewScanner(bytes.NewReader(text)).Master()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      for i, value := range master.Stream {
         if i >= 1 {
            fmt.Println()
         }
         fmt.Println(value)
      }
      fmt.Println()
   }
}

var master_tests = []string{
   "m3u8/cbc-master.m3u8.txt",
   "m3u8/nbc-master.m3u8.txt",
   "m3u8/roku-master.m3u8.txt",
}

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

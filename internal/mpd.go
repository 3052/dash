package main

import (
   "os"
   "strings"
)

const max = 5

func main() {
   file, err := os.Create("out.mpd")
   if err != nil {
      panic(err)
   }
   defer file.Close()
   text, err := os.ReadFile("stream.mpd")
   if err != nil {
      panic(err)
   }
   var count int
   for _, line := range strings.SplitAfter(string(text), "\n") {
      if strings.Contains(line, "<S ") {
         count++
      }
      if strings.Contains(line, "</SegmentTimeline") {
         count = 0
      }
      if count <= max {
         file.WriteString(line)
      }
   }
}

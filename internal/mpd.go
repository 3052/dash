package main

import (
   "flag"
   "os"
   "strings"
)

func main() {
   max := flag.Int("m", 9, "max")
   input := flag.String("i", "", "input")
   output := flag.String("o", "", "output")
   flag.Parse()
   if *input != "" {
      file, err := os.Create(*output)
      if err != nil {
         panic(err)
      }
      defer file.Close()
      text, err := os.ReadFile(*input)
      if err != nil {
         panic(err)
      }
      var count int
      for _, line := range strings.SplitAfter(string(text), "\n") {
         if strings.Contains(line, "<S ") {
            count++
         }
         if strings.Contains(line, "</SegmentTimeline>") {
            count = 0
         }
         if count <= *max {
            file.WriteString(line)
         }
      }
   } else {
      flag.Usage()
   }
}

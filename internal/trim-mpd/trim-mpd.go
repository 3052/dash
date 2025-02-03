package main

import (
   "flag"
   "os"
   "strings"
)

func main() {
   max := flag.Int("m", 2, "max")
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
      var (
         period_count      int
         s_count           int
         segment_url_count int
      )
      for _, line := range strings.SplitAfter(string(text), "\n") {
         if strings.Contains(line, "<Period") {
            period_count++
         }
         if strings.Contains(line, "</MPD>") {
            period_count = 0
         }
         if strings.Contains(line, "<S ") {
            s_count++
         }
         if strings.Contains(line, "</SegmentTimeline>") {
            s_count = 0
         }
         if strings.Contains(line, "<SegmentURL") {
            segment_url_count++
         }
         if strings.Contains(line, "</SegmentList>") {
            segment_url_count = 0
         }
         if period_count <= *max {
            if s_count <= *max {
               if segment_url_count <= *max {
                  file.WriteString(line)
               }
            }
         }
      }
   } else {
      flag.Usage()
   }
}

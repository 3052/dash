package main

import (
   "flag"
   "os"
   "strings"
)

func ok(input, output string) bool {
   if input != "" {
      if output != "" {
         return true
      }
   }
   return false
}

func main() {
   newMax := flag.Int("m", 2, "max")
   input := flag.String("i", "", "input")
   output := flag.String("o", "", "output")
   flag.Parse()
   if ok(*input, *output) {
      err := do_trim(*input, *output, *newMax)
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

func do_trim(input, output string, newMax int) error {
   file, err := os.Create(output)
   if err != nil {
      return err
   }
   defer file.Close()
   data, err := os.ReadFile(input)
   if err != nil {
      return err
   }
   var (
      period_count      int
      s_count           int
      segment_url_count int
   )
   for _, line := range strings.SplitAfter(string(data), "\n") {
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
      if period_count <= newMax {
         if s_count <= newMax {
            if segment_url_count <= newMax {
               _, err = file.WriteString(line)
               if err != nil {
                  return err
               }
            }
         }
      }
   }
   return nil
}

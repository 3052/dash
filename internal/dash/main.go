package main

import (
   "flag"
   "fmt"
   "net/url"
   "sort"
   "time"
)

type flags struct {
   address string
   id string
   url *url.URL
   channels int
}

func main() {
   var f flags
   flag.IntVar(&f.channels, "c", 3, "channels")
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.id, "i", "", "ID")
   flag.Parse()
   if f.address != "" {
      reps, err := f.manifest()
      if err != nil {
         panic(err)
      }
      index := func() int {
         for i, rep := range reps {
            if rep.ID == f.id {
               return i
            }
         }
         return -1
      }()
      if index >= 0 {
         begin := time.Now()
         err := f.download(reps[index])
         if err != nil {
            panic(err)
         }
         fmt.Println(time.Since(begin))
      } else {
         sort.Slice(reps, func(i, j int) bool {
            return reps[i].Bandwidth < reps[j].Bandwidth
         })
         for i, rep := range reps {
            if i >= 1 {
               fmt.Println()
            }
            fmt.Println(rep)
         }
      }
   } else {
      flag.Usage()
   }
}

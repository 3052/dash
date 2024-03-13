package main

import (
   "flag"
   "fmt"
   "net/url"
   "sort"
)

type flags struct {
   address string
   id string
   url *url.URL
}

func main() {
   var f flags
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.id, "i", "", "ID")
   flag.Parse()
   if f.address != "" {
      reps, err := f.manifest()
      if err != nil {
         panic(err)
      }
      if f.id != "" {
         for _, rep := range reps {
            if rep.ID == f.id {
               if err := f.download(rep); err != nil {
                  panic(err)
               }
            }
         }
      } else {
         sort.Slice(reps, func(i, j int) bool {
            return reps[j].Bandwidth < reps[i].Bandwidth
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

package main

import (
   "flag"
   "fmt"
   "net/url"
)

type flags struct {
   address string
   id string
   url *url.URL
}

func main() {
   var f flags
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.id, "id", "", "ID")
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

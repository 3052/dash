package main

import (
   "flag"
   "fmt"
)

type flags struct {
   address string
   id string
   period string
}

func main() {
   var f flags
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.id, "id", "", "ID")
   flag.StringVar(&f.period, "p", "", "period")
   flag.Parse()
   if f.address != "" {
      reps, err := f.manifest()
      if err != nil {
         panic(err)
      }
      if f.id != "" {
         if rep, ok := f.pick(reps); ok {
            if err := download(rep); err != nil {
               panic(err)
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

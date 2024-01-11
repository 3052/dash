package main

import (
   "154.pages.dev/encoding/dash"
   "errors"
   "flag"
   "fmt"
   "net/http"
)

func main() {
   var f flags
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.period, "p", "", "period")
   flag.Parse()
   if f.address != "" {
      reps, err := f.get()
      if err != nil {
         panic(err)
      }
      for i, rep := range reps {
         if i >= 1 {
            fmt.Println()
         }
         fmt.Println(rep)
      }
   } else {
      flag.Usage()
   }
}

func (f flags) get() ([]*dash.Representation, error) {
   res, err := http.Get(f.address)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   var media dash.Media
   if err := media.Decode(res.Body); err != nil {
      return nil, err
   }
   return media.Representation(f.period)
}

type flags struct {
   address string
   period string
}

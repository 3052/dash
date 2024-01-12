package main

import (
   "flag"
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
      media, err := f.manifest()
      if err != nil {
         panic(err)
      }
      if f.id != "" {
         if rep, ok := f.pick(media); ok {
            if err := f.download(rep); err != nil {
               panic(err)
            }
         }
      } else {
         err := execute(media)
         if err != nil {
            panic(err)
         }
      }
   } else {
      flag.Usage()
   }
}

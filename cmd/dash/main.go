package main

import (
   "154.pages.dev/encoding/dash"
   "flag"
   "net/url"
   "os"
   "text/template"
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
         tmpl, err := new(template.Template).Parse(dash.Template)
         if err != nil {
            panic(err)
         }
         if err := tmpl.Execute(os.Stdout, media); err != nil {
            panic(err)
         }
      }
   } else {
      flag.Usage()
   }
}

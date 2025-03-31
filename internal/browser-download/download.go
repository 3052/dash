package main

import (
   "fmt"
   "io"
   "log"
   "net/http"
   "path"
   "strings"
)

const address = "https://cdimage.debian.org/debian-cd/current-live/amd64/iso-hybrid/"

const homepage = `
<ul>
   <li>
      <a href="/get/debian-live-12.10.0-amd64-standard.iso">debian-live-12.10.0-amd64-standard.iso</a> (1.4G)
   </li>
   <li>
      <a href="/get/debian-live-12.10.0-amd64-lxde.iso">debian-live-12.10.0-amd64-lxde.iso</a> (3.0G)
   </li>
</ul>
`

func handler(rw http.ResponseWriter, req *http.Request) {
   switch {
   case req.URL.Path == "/":
      rw.Header().Set("content-type", "text/html")
      fmt.Fprint(rw, homepage)
   case strings.HasPrefix(req.URL.Path, "/get/"):
      resp, err := http.Get(address + path.Base(req.URL.Path))
      if err != nil {
         fmt.Fprint(rw, err)
         return
      }
      defer resp.Body.Close()
      _, err = io.Copy(rw, resp.Body)
      if err != nil {
         fmt.Fprint(rw, err)
         return
      }
   }
}

const port = ":99"

func main() {
   log.SetFlags(log.Ltime)
   log.Print("localhost", port)
   err := http.ListenAndServe(port, http.HandlerFunc(handler))
   if err != nil {
      panic(err)
   }
}

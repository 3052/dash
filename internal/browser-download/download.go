package main

import (
   "fmt"
   "log"
   "net/http"
   "strings"
   "time"
)

const homepage = `
<ul>
   <li>
      <a href="/get/cdimage.debian.org/debian-cd/current-live/amd64/iso-hybrid/debian-live-12.10.0-amd64-standard.iso"
      >debian-live-12.10.0-amd64-standard.iso</a>
   </li>
   <li>
      <a target="_blank" href="/get/dl.google.com/go/go1.24.1.windows-amd64.zip"
      >go1.24.1.windows-amd64.zip</a>
   </li>
</ul>
`

func handler(rw http.ResponseWriter, req *http.Request) {
   switch {
   case req.URL.Path == "/":
      rw.Header().Set("content-type", "text/html")
      fmt.Fprint(rw, homepage)
   case strings.HasPrefix(req.URL.Path, "/get/"):
      log.SetOutput(rw)
      flush, ok := rw.(http.Flusher)
      for range 9 {
         log.Print("hello world")
         if ok {
            flush.Flush()
         }
         time.Sleep(time.Second)
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

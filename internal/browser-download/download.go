package main

import (
   "fmt"
   "io"
   "log"
   "net/http"
   "strings"
)

const mp4 = "https://vod-lf-oneapp2-prd.akamaized.net/prod/nbc/J39/rhk/9000405641/1743322955758-hqIHK/cmaf/mpeg_cenc_2sec/7830k_1080_cmaf/_237142984_1.mp4"

const homepage = `
<ul>
   <li>
      <a href="/get/alfa">alfa</a>
   </li>
   <li>
      <a href="/get/bravo">bravo</a>
   </li>
</ul>
`

func handler(rw http.ResponseWriter, req *http.Request) {
   switch {
   case req.URL.Path == "/":
      rw.Header().Set("content-type", "text/html")
      fmt.Fprint(rw, homepage)
   case strings.HasPrefix(req.URL.Path, "/get/"):
      resp, err := http.Get(mp4)
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

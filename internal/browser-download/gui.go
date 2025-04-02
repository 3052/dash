package main

import (
   "fmt"
   "log"
   "net/http"
   "strings"
   "time"
)

const homepage = `
<a target="_blank" href="/get/alfa">alfa</a>
`

const log_page = `
<script>
'use strict';

const observer = new MutationObserver(() => {
   window.scrollTop = window.scrollHeight;
});

observer.observe(window, {
  childList: true,
  subtree: true
});

</script>
`

func handler(rw http.ResponseWriter, req *http.Request) {
   rw.Header().Set("content-type", "text/html")
   switch {
   case req.URL.Path == "/":
      fmt.Fprint(rw, homepage)
   case strings.HasPrefix(req.URL.Path, "/get/"):
      logger := log.New(rw, "", log.Ltime)
      flush, ok := rw.(http.Flusher)
      logger.Print(log_page)
      if ok {
         flush.Flush()
      }
      for range 9 {
         logger.Print("hello world")
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

package main

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
)

func hello(u *url.URL, b []byte) {
   fmt.Println(string(b), u)
}

func main() {
   res, err := http.Get("http://example.com")
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   text, err := io.ReadAll(res.Body)
   if err != nil {
      panic(err)
   }
   hello(res.Request.URL, text)
}

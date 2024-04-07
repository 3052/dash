package main

import (
   "fmt"
   "io"
   "net/http"
)

func hello(s string, b []byte) {
   fmt.Println(string(b), s)
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
   hello(res.Request.URL.String(), text)
}

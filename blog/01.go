package main

import (
   "fmt"
   "io"
   "net/http"
)

func hello(s string, r io.Reader) error {
   text, err := io.ReadAll(r)
   if err != nil {
      return err
   }
   fmt.Println(string(text), s)
   return nil
}

func main() {
   res, err := http.Get("http://example.com")
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   hello(res.Request.URL.String(), res.Body)
}


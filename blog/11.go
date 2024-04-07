package main

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
)

func hello(u *url.URL, r io.Reader) error {
   text, err := io.ReadAll(r)
   if err != nil {
      return err
   }
   fmt.Println(string(text), u)
   return nil
}

func main() {
   res, err := http.Get("http://example.com")
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   hello(res.Request.URL, res.Body)
}


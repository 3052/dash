package main

import (
   "154.pages.dev/encoding/dash"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "os"
   "path"
)

func (f *flags) manifest() ([]dash.Representation, error) {
   res, err := http.Get(f.address)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   f.url = res.Request.URL
   text, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   return dash.Unmarshal(text)
}

func (f flags) download(rep dash.Representation) error {
   initialization, ok := rep.Initialization()
   if !ok {
      return errors.New("dash.Representation.Initialization")
   }
   address, err := f.url.Parse(initialization)
   if err != nil {
      return err
   }
   file, err := os.Create(path.Base(initialization))
   if err != nil {
      return err
   }
   defer file.Close()
   if err := write(address, file); err != nil {
      return err
   }
   media := rep.Media()
   for i, ref := range media {
      // with DASH, initialization and media URLs are relative to the MPD URL
      address, err := f.url.Parse(ref)
      if err != nil {
         return err
      }
      fmt.Println(len(media)-i)
      if err := write(address, file); err != nil {
         return err
      }
   }
   return nil
}

func write(address *url.URL, file *os.File) error {
   res, err := http.Get(address.String())
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return errors.New(res.Status)
   }
   if _, err := file.ReadFrom(res.Body); err != nil {
      return err
   }
   return nil
}

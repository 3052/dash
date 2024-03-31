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

func (f flags) download(rep dash.Representation) error {
   template, ok := rep.GetSegmentTemplate()
   if !ok {
      return errors.New("GetSegmentTemplate")
   }
   initial, ok := template.GetInitialization(rep)
   if !ok {
      return errors.New("GetInitialization")
   }
   address, err := f.url.Parse(initial)
   if err != nil {
      return err
   }
   if err := create(address); err != nil {
      return err
   }
   media, err := template.GetMedia(rep)
   if err != nil {
      return err
   }
   for i, medium := range media {
      fmt.Println(len(media)-i)
      // with DASH, initialization and media URLs are relative to the MPD URL
      url, err := f.url.Parse(medium)
      if err != nil {
         return err
      }
      if err := create(url); err != nil {
         return err
      }
   }
   return nil
}

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

func create(url *url.URL) error {
   res, err := http.Get(url.String())
   if err != nil {
      return err
   }
   defer res.Body.Close()
   file, err := os.Create(path.Base(url.Path))
   if err != nil {
      return err
   }
   defer file.Close()
   if _, err := file.ReadFrom(res.Body); err != nil {
      return err
   }
   return nil
}

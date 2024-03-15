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

func (f flags) download(rep dash.Representation) error {
   initialization, ok := rep.Initialization()
   if !ok {
      return errors.New("dash.Representation.Initialization")
   }
   address, err := f.url.Parse(initialization)
   if err != nil {
      return err
   }
   if err := create(address); err != nil {
      return err
   }
   media := rep.Media()
   errs := make(chan error)
   for channel := range f.channels {
      go f.loop(media, channel, errs)
   }
   for range f.channels {
      err := <-errs
      if err != nil {
         return err
      }
   }
   return nil
}

func (f flags) loop(media []string, channel int, errs chan error) {
   for i := channel; i < len(media); i += f.channels {
      // with DASH, initialization and media URLs are relative to the MPD URL
      url, err := f.url.Parse(media[i])
      if err != nil {
         errs <- err
         return
      }
      fmt.Printf(
         "channel:%v i:%v len:%v\n", channel, i, len(media),
      )
      if err := create(url); err != nil {
         errs <- err
         return
      }
   }
   errs <- nil
}

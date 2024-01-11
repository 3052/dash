package main

import (
   "154.pages.dev/encoding/dash"
   "errors"
   "net/http"
)

func download(*dash.Representation) error {
   file, err := os.Create(s.Name + ext)
   if err != nil {
      return err
   }
   defer file.Close()
   req, err := http.NewRequest("GET", initialization, nil)
   if err != nil {
      return err
   }
   req.URL = s.Base.ResolveReference(req.URL)
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if err := encode_init(file, res.Body); err != nil {
      return err
   }
   media, ok := item.Media()
   if !ok {
      return errors.New("Media")
   }
   src := log.New_Progress(len(media))
   log.Set_Transport(slog.LevelDebug)
   defer log.Set_Transport(slog.LevelInfo)
   for _, ref := range media {
      // with DASH, initialization and media URLs are relative to the MPD URL
      req.URL, err = s.Base.Parse(ref)
      if err != nil {
         return err
      }
      err := func() error {
         res, err := http.DefaultClient.Do(req)
         if err != nil {
            return err
         }
         defer res.Body.Close()
         if res.StatusCode != http.StatusOK {
            return errors.New(res.Status)
         }
         return encode_segment(file, src.Reader(res), key)
      }()
      if err != nil {
         return err
      }
   }
   return nil
}

func (f flags) pick(reps []*dash.Representation) (*dash.Representation, bool) {
   for _, rep := range reps {
      if rep.ID == f.id {
         return rep, true
      }
   }
   return nil, false
}

func (f flags) manifest() ([]*dash.Representation, error) {
   res, err := http.Get(f.address)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   var media dash.Media
   if err := media.Decode(res.Body); err != nil {
      return nil, err
   }
   return media.Representation(f.period)
}

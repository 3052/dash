package main

import (
   "154.pages.dev/dash"
   "errors"
   "flag"
   "fmt"
   "io"
   "net/http"
   "os"
   "path"
   "sort"
   "time"
)

func create(base string, i int, req *http.Request) error {
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return errors.New(resp.Status)
   }
   file, err := os.Create(
      base + fmt.Sprintf("%03v", i) + path.Ext(req.URL.Path),
   )
   if err != nil {
      return err
   }
   defer file.Close()
   if _, err := file.ReadFrom(resp.Body); err != nil {
      return err
   }
   return nil
}

func main() {
   var f flags
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.id, "i", "", "ID")
   flag.Parse()
   if f.address != "" {
      reps, err := f.manifest()
      if err != nil {
         panic(err)
      }
      index := func() int {
         for i, rep := range reps {
            if rep.Id == f.id {
               return i
            }
         }
         return -1
      }()
      if index >= 0 {
         begin := time.Now()
         err := f.download(reps[index])
         if err != nil {
            panic(err)
         }
         fmt.Println(time.Since(begin))
      } else {
         sort.Slice(reps, func(i, j int) bool {
            return reps[i].Bandwidth < reps[j].Bandwidth
         })
         for i, rep := range reps {
            if i >= 1 {
               fmt.Println()
            }
            fmt.Println(&rep)
         }
      }
   } else {
      flag.Usage()
   }
}

type flags struct {
   address string
   id string
}

func (f *flags) manifest() ([]dash.Representation, error) {
   resp, err := http.Get(f.address)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   text, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return dash.Unmarshal(text, resp.Request.URL)
}

func (f flags) download(rep dash.Representation) error {
   base, ok := rep.GetBaseUrl()
   if !ok {
      return errors.New("GetBaseUrl")
   }
   initial, ok := rep.Initialization()
   if !ok {
      return errors.New("Initialization")
   }
   req, err := http.NewRequest("", initial, nil)
   if err != nil {
      return err
   }
   req.URL = base.Url.ResolveReference(req.URL)
   err = create("init-", 0, req)
   if err != nil {
      return err
   }
   media := rep.Media()
   for i, medium := range media {
      fmt.Println(len(media)-i)
      req.URL, err = base.Url.Parse(medium)
      if err != nil {
         return err
      }
      err = create("segment-", i, req)
      if err != nil {
         return err
      }
   }
   return nil
}

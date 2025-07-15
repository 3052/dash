package main

import (
   "bytes"
   "flag"
   "mime/multipart"
   "net/http"
   "os"
   "os/exec"
   "path/filepath"
)

func main() {
   name := flag.String("n", "", "name")
   flag.Parse()
   if *name != "" {
      err := post(*name)
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

func (f *form_data) New(name string) error {
   writer := multipart.NewWriter(&f.body)
   defer writer.Close()
   field, err := writer.CreateFormField("purpose")
   if err != nil {
      return err
   }
   _, err = field.Write([]byte("file-extract"))
   if err != nil {
      return err
   }
   file, err := writer.CreateFormFile("file", filepath.Base(name))
   if err != nil {
      return err
   }
   data, err := os.ReadFile(name)
   if err != nil {
      return err
   }
   _, err = file.Write(data)
   if err != nil {
      return err
   }
   f.content_type = writer.FormDataContentType()
   return nil
}

type form_data struct {
   body bytes.Buffer
   content_type string
}

func post(name string) error {
   bearer, err := exec.Command("password", "-i", "moonshot.ai").Output()
   if err != nil {
      return err
   }
   var form form_data
   err = form.New(name)
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://api.moonshot.ai/v1/files", &form.body,
   )
   if err != nil {
      return err
   }
   req.Header.Set("authorization", "Bearer " + string(bearer))
   req.Header.Set("content-type", form.content_type)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   return resp.Write(os.Stdout)
}

package main

import (
   "41.neocities.org/dash"
   "fmt"
   "html/template"
   "io"
   "log"
   "net/http"
)

func world(rw http.ResponseWriter, req *http.Request) {
   switch req.Method {
   case "GET":
      if req.URL.Path == "/" {
         rw.Header().Set("content-type", "text/html")
         fmt.Fprint(rw, form)
      } else {
         log.Println(req.Method, req.URL)
      }
   case "POST":
      file, _, err := req.FormFile("mpd")
      if err != nil {
         fmt.Fprint(rw, err)
         return
      }
      defer file.Close()
      data, err := io.ReadAll(file)
      if err != nil {
         fmt.Fprint(rw, err)
         return
      }
      var mpd dash.Mpd
      err = mpd.Unmarshal(data)
      if err != nil {
         fmt.Fprint(rw, err)
         return
      }
      plate, err := template.New("table").Parse(table)
      if err != nil {
         fmt.Fprint(rw, err)
         return
      }
      err = plate.Execute(rw, mpd.Representation())
      if err != nil {
         fmt.Fprint(rw, err)
      }
   }
}

const form = `
<!doctype html>
<form method="post" enctype="multipart/form-data">
   <input type="file" name="mpd" />
   <input type="submit" />
</form>
`

const port = ":99"

func main() {
   log.SetFlags(log.Ltime)
   log.Print("localhost", port)
   err := http.ListenAndServe(port, http.HandlerFunc(world))
   if err != nil {
      panic(err)
   }
}

const table = `
<style>
</style>
<table>
   <tr>
      <th>bandwidth</th>
      <th>width</th>
      <th>height</th>
      <th>mimeType</th>
      <th>codecs</th>
      <th>id</th>
   </tr>
{{ range $element := . }}
   <tr>
      <td>{{ .Bandwidth }}</td>
      <td>{{ with .Width }}{{ . }}{{ end }}</td>
      <td>{{ with .Height }}{{ . }}{{ end }}</td>
      <td>{{ .MimeType }}</td>
      <td>{{ with .Codecs }}{{ . }}{{ end }}</td>
      <td>{{ .Id }}</td>
   </tr>
{{ end }}
</table>
`

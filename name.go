package encoding

import (
   "strconv"
   "strings"
   "text/template"
)

const Format = 
   "{{if .Show}}" +
      "{{.Show}}" +
      " - {{.Season}}" +
      " {{.Episode}}" +
      " - {{.Title}}" +
   "{{else}}" +
      "{{.Title}}" +
      " - {{.Year}}" +
   "{{end}}"

type Namer interface {
   Show() string
   Season() string
   Episode() string
   Title() string
   Year() string
}

func Name(format string, v Namer) (string, error) {
   text, err := new(template.Template).Parse(format)
   if err != nil {
      return "", err
   }
   var b strings.Builder
   if err := text.Execute(&b, v); err != nil {
      return "", err
   }
   return b.String(), nil
}

func clean(path []byte) {
   m := map[byte]bool{
      '"': true,
      '*': true,
      '/': true,
      ':': true,
      '<': true,
      '>': true,
      '?': true,
      '\\': true,
      '|': true,
   }
   for k, v := range path {
      if m[v] {
         path[k] = '-'
      }
   }
}

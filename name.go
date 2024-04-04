package encoding

import (
   "strings"
   "text/template"
)

var Format = 
   "{{if .Show}}" +
      "{{.Show}} - {{.Season}} {{.Episode}} - {{.Title}}" +
   "{{else}}" +
      "{{.Title}} - {{.Year}}" +
   "{{end}}"

func Clean(name string) string {
   mapping := func(r rune) rune {
      if strings.ContainsRune(`"*/:<>?\|`, r) {
         return '-'
      }
      return r
   }
   return strings.Map(mapping, name)
}

func Name(n Namer) (string, error) {
   text, err := new(template.Template).Parse(Format)
   if err != nil {
      return "", err
   }
   var b strings.Builder
   if err := text.Execute(&b, n); err != nil {
      return "", err
   }
   return b.String(), nil
}

type Namer interface {
   Show() string
   Season() int
   Episode() int
   Title() string
   Year() int
}

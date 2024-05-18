package encoding

import (
   "bytes"
   "strings"
   "text/template"
)

var Format = 
   "{{if .Show}}" +
      "{{.Show}} - {{.Season}} {{.Episode}} - {{.Title}}" +
   "{{else}}" +
      "{{.Title}} - {{.Year}}" +
   "{{end}}"

func Clean(s string) string {
   mapping := func(r rune) rune {
      if strings.ContainsRune(`"*/:<>?\|`, r) {
         return '-'
      }
      return r
   }
   return strings.Map(mapping, s)
}

func Cut(s, before, after []byte) ([]byte, []byte) {
   i := bytes.Index(s, append(before, after...))
   if i == -1 {
      return s, nil
   }
   i += len(before)
   return s[:i], s[i:]
}

func Name(n Namer) (string, error) {
   text, err := new(template.Template).Parse(Format)
   if err != nil {
      return "", err
   }
   var b strings.Builder
   err = text.Execute(&b, n)
   if err != nil {
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

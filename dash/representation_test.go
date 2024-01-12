package dash

import (
   "os"
   "testing"
   "text/template"
)

var tests = []string{
   "mpd/amc.mpd",
   "mpd/hulu.mpd",
   "mpd/nbc.mpd",
   "mpd/paramount.mpd",
   "mpd/roku.mpd",
}

func Test_Info(t *testing.T) {
   tmpl, err := template.ParseFiles("in.html")
   if err != nil {
      t.Fatal(err)
   }
   dst, err := os.Create("out.html")
   if err != nil {
      t.Fatal(err)
   }
   defer dst.Close()
   for _, name := range tests {
      dst.WriteString(name)
      func() {
         src, err := os.Open(name)
         if err != nil {
            t.Fatal(err)
         }
         defer src.Close()
         var m Media
         if err := m.Decode(src); err != nil {
            t.Fatal(err)
         }
         if err := tmpl.Execute(dst, m); err != nil {
            t.Fatal(err)
         }
      }()
   }
}

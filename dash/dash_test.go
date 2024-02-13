package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
   "text/template"
)

func TestRange(t *testing.T) {
   media, err := reader("mpd/hulu.mpd")
   if err != nil {
      t.Fatal(err)
   }
   media.Visit(func(p Pointer) {
      sb := p.Representation.SegmentBase
      start, end, err := sb.Initialization.Range.Scan()
      fmt.Printf("%v %v %v ", start, end, err)
      start, end, err = sb.IndexRange.Scan()
      fmt.Printf("%v %v %v\n", start, end, err)
   })
}

var tests = []string{
   "mpd/amc.mpd",
   "mpd/hulu.mpd",
   "mpd/nbc.mpd",
   "mpd/paramount.mpd",
   "mpd/roku.mpd",
}

func reader(name string) (*MPD, error) {
   text, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   media := new(MPD)
   if err := xml.Unmarshal(text, &media); err != nil {
      return nil, err
   }
   return media, nil
}

func TestModeLine(t *testing.T) {
   tmpl, err := new(template.Template).Parse(ModeLine)
   if err != nil {
      t.Fatal(err)
   }
   for i, name := range tests {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(name)
      text, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      var media MPD
      if err := xml.Unmarshal(text, &media); err != nil {
         t.Fatal(err)
      }
      if err := tmpl.Execute(os.Stdout, media); err != nil {
         t.Fatal(err)
      }
   }
}

func TestProtection(t *testing.T) {
   for _, test := range tests {
      text, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      var media MPD
      if err := xml.Unmarshal(text, &media); err != nil {
         t.Fatal(err)
      }
      media.Visit(func(p Pointer) {
         _, pssh := p.PSSH()
         _, kid := p.Default_KID()
         fmt.Printf(
            "mpd:%v period:%q type:%v pssh:%v kid:%v\n",
            test, p.Period.ID, p.MimeType(), pssh, kid,
         )
      })
   }
}

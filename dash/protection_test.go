package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

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
      media.Every(func(p Pointer) {
         _, pssh := p.PSSH()
         _, kid := p.Default_KID()
         fmt.Printf(
            "mpd:%v period:%q type:%v pssh:%v kid:%v\n",
            test, p.Period.ID, p.MimeType(), pssh, kid,
         )
      })
   }
}

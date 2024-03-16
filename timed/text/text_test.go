package text

import (
   "fmt"
   "os"
   "testing"
)

func TestText(t *testing.T) {
   text, err := os.ReadFile("text.xml")
   if err != nil {
      t.Fatal(err)
   }
   var mark Markup
   if err := mark.Unmarshal(text); err != nil {
      t.Fatal(err)
   }
   fmt.Println(WebVtt)
   for _, p := range mark.Body.Div.P {
      fmt.Println(p)
   }
}

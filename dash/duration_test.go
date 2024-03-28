package dash

import (
   "fmt"
   "os"
   "testing"
)

func reader(name string) ([]Representation, error) {
   text, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   return Unmarshal(text)
}

func TestMedia(t *testing.T) {
   for _, test := range media_tests {
      fmt.Println(test[0] + ":")
      reps, err := reader(test[0])
      if err != nil {
         t.Fatal(err)
      }
      for _, media := range reps[0].Media() {
         fmt.Println(test[1] + media)
      }
   }
}

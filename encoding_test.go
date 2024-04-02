package encoding

import (
   "fmt"
   "testing"
)

func TestName(t *testing.T) {
   namers := []Namer{
      episode{}, film{},
   }
   for _, namer := range namers {
      name, err := Name(NameFormat, namer)
      if err != nil {
         panic(err)
      }
      fmt.Printf("%q\n", name)
   }
}

type episode struct{}

func (episode) Show() string {
   return "show"
}

func (episode) Title() string {
   return "title"
}

type film struct{}

func (film) Show() string {
   return ""
}

func (film) Title() string {
   return "title"
}

func TestPercent(t *testing.T) {
   fmt.Println(Percent(1234) / 10000)
}

func (episode) Season() int {
   return 2
}

func (episode) Episode() int {
   return 3
}

func (episode) Year() int {
   return 2024
}

func (film) Season() int {
   return 0
}

func (film) Episode() int {
   return 0
}

func (film) Year() int {
   return 2024
}

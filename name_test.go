package encoding

import (
   "154.pages.dev/encoding/blog/encoding"
   "flag"
   "fmt"
)

func main() {
   format := flag.String("f", encoding.Format, "format")
   flag.Parse()
   namers := []encoding.Namer{
      episode{}, film{},
   }
   for _, namer := range namers {
      name, err := encoding.Name(*format, namer)
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

func (episode) Season() string {
   return "season"
}

func (episode) Episode() string {
   return "episode"
}

func (episode) Title() string {
   return "title"
}

func (episode) Year() string {
   return "year"
}

type film struct{}

func (film) Show() string {
   return ""
}

func (film) Season() string {
   return ""
}

func (film) Episode() string {
   return ""
}

func (film) Title() string {
   return "title"
}

func (film) Year() string {
   return "year"
}

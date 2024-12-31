package dash

import "strings"

type Initialization func(string) string

func (i *Initialization) UnmarshalText(data []byte) error {
   *i = func(id string) string {
      return strings.Replace(string(data), "$RepresentationID$", id, 1)
   }
   return nil
}

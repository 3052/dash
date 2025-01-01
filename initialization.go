package dash

import (
   "net/url"
   "strings"
)

type Initialization func(string) *url.URL

func (i *Initialization) UnmarshalText(data []byte) error {
   var u url.URL
   err := u.UnmarshalBinary(data)
   if err != nil {
      return err
   }
   *i = func(s string) *url.URL {
      u.Path = strings.Replace(u.Path, "$RepresentationID$", s, 1)
      return &u
   }
   return nil
}

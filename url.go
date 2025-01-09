package dash

import "net/url"

type Url struct {
   Url *url.URL
}

func (u *Url) UnmarshalText(data []byte) error {
   if u.Url == nil {
      u.Url = &url.URL{}
      return u.Url.UnmarshalBinary(data)
   }
   var err error
   u.Url, err = u.Url.Parse(string(data))
   if err != nil {
      return err
   }
   return nil
}

type LazyUrl struct {
   Url string
}

func (u *LazyUrl) UnmarshalText(data []byte) error {
   u.Url = string(data)
   return nil
}

//func (u Url) Parse(base *url.URL) (*url.URL, error) {
//   return base.Parse(u.S)
//}

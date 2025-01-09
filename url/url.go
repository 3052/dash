package url

import "net/url"

func Parse(s string) (*Url, error) {
   u, err := url.Parse(s)
   if err != nil {
      return nil, err
   }
   return &Url{u}, nil
}

func New() *Url {
   return &Url{&url.URL{}}
}

func (b *Url) UnmarshalText(data []byte) error {
   if b.Url == nil {
      b.Url = &url.URL{}
      return b.Url.UnmarshalBinary(data)
   }
   var err error
   b.Url, err = b.Url.Parse(string(data))
   if err != nil {
      return err
   }
   return nil
}

type Url struct {
   Url *url.URL
}

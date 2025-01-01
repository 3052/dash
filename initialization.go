package dash

import "net/url"

type Initialization func(*Representation) (*url.URL, error)

func (i *Initialization) UnmarshalText(data []byte) error {
   *i = func(r *Representation) (*url.URL, error) {
      u := &url.URL{}
      err := u.UnmarshalBinary(data)
      if err != nil {
         return nil, err
      }
      u.Path = execute(u.Path, r.Id, 0)
      if r.BaseUrl != nil {
         u = r.BaseUrl.Url.ResolveReference(u)
      }
      return u, nil
   }
   return nil
}

package dash

import (
   "fmt"
   "net/url"
   "strings"
)

type Media func(*Representation, int) (*url.URL, error)

func replace(s, from, to string) string {
   return strings.Replace(s, from, to, 1)
}

func execute(s, represent string, i int) string {
   s = replace(s, "$Number$", fmt.Sprint(i))
   s = replace(s, "$Number%02d$", fmt.Sprintf("%02d", i))
   s = replace(s, "$Number%03d$", fmt.Sprintf("%03d", i))
   s = replace(s, "$Number%04d$", fmt.Sprintf("%04d", i))
   s = replace(s, "$Number%05d$", fmt.Sprintf("%05d", i))
   s = replace(s, "$Number%06d$", fmt.Sprintf("%06d", i))
   s = replace(s, "$Number%07d$", fmt.Sprintf("%07d", i))
   s = replace(s, "$Number%08d$", fmt.Sprintf("%08d", i))
   s = replace(s, "$Number%09d$", fmt.Sprintf("%09d", i))
   s = replace(s, "$RepresentationID$", represent)
   return replace(s, "$Time$", fmt.Sprint(i))
}

func (m *Media) UnmarshalText(data []byte) error {
   *m = func(r *Representation, i int) (*url.URL, error) {
      u := &url.URL{}
      err := u.UnmarshalBinary(data)
      if err != nil {
         return nil, err
      }
      u.Path = execute(u.Path, r.Id, i)
      u.RawQuery = execute(u.RawQuery, "", i)
      if r.BaseUrl != nil {
         u = r.BaseUrl.Url.ResolveReference(u)
      }
      return u, nil
   }
   return nil
}

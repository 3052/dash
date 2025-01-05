package dash

import (
   "fmt"
   "net/url"
   "strings"
)

func (m Media) Url(r *Representation, i int) (*url.URL, error) {
   replace(&m.S, "$RepresentationID$", r.Id)
   if m.time() {
      replace(&m.S, "$Time$", fmt.Sprint(i))
   } else {
      replace(&m.S, "$Number$", fmt.Sprint(i))
      replace(&m.S, "$Number%02d$", fmt.Sprintf("%02d", i))
      replace(&m.S, "$Number%03d$", fmt.Sprintf("%03d", i))
      replace(&m.S, "$Number%04d$", fmt.Sprintf("%04d", i))
      replace(&m.S, "$Number%05d$", fmt.Sprintf("%05d", i))
      replace(&m.S, "$Number%06d$", fmt.Sprintf("%06d", i))
      replace(&m.S, "$Number%07d$", fmt.Sprintf("%07d", i))
      replace(&m.S, "$Number%08d$", fmt.Sprintf("%08d", i))
      replace(&m.S, "$Number%09d$", fmt.Sprintf("%09d", i))
   }
   u, err := url.Parse(m.S)
   if err != nil {
      return nil, err
   }
   if r.BaseUrl != nil {
      u = r.BaseUrl.Url.ResolveReference(u)
   }
   return u, nil
}

func (m Media) time() bool {
   return strings.Contains(m.S, "$Time$")
}

type Media struct {
   S string
}

func (m *Media) UnmarshalText(data []byte) error {
   m.S = string(data)
   return nil
}

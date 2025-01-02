package dash

import (
   "fmt"
   "net/url"
   "strings"
)

type Media struct {
   S string
}

func (m Media) Url(r *Representation, i int) (*url.URL, error) {
   m.execute(r.Id, i)
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

func (m *Media) replace(from, to string) {
   m.S = strings.Replace(m.S, from, to, 1)
}

func (m *Media) execute(represent string, i int) {
   m.replace("$RepresentationID$", represent)
   if m.time() {
      m.replace("$Time$", fmt.Sprint(i))
   } else {
      m.replace("$Number$", fmt.Sprint(i))
      m.replace("$Number%02d$", fmt.Sprintf("%02d", i))
      m.replace("$Number%03d$", fmt.Sprintf("%03d", i))
      m.replace("$Number%04d$", fmt.Sprintf("%04d", i))
      m.replace("$Number%05d$", fmt.Sprintf("%05d", i))
      m.replace("$Number%06d$", fmt.Sprintf("%06d", i))
      m.replace("$Number%07d$", fmt.Sprintf("%07d", i))
      m.replace("$Number%08d$", fmt.Sprintf("%08d", i))
      m.replace("$Number%09d$", fmt.Sprintf("%09d", i))
   }
}

func (m *Media) UnmarshalText(data []byte) error {
   m.S = string(data)
   return nil
}

package dash

import (
   "fmt"
   "strings"
)

func (m *Media) UnmarshalText(data []byte) error {
   *m = func(represent string, i int) string {
      s := string(data)
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
      s = replace(s, "$Time$", fmt.Sprint(i))
      return s
   }
   return nil
}

func replace(s, a, b string) string {
   return strings.Replace(s, a, b, 1)
}

type Media func(string, int) string

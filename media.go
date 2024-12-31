package dash

import (
   "fmt"
   "strings"
)

// %0[width]d
// 23009-1
// standards.iso.org/ittf/PubliclyAvailableStandards
func (m *Media) UnmarshalText(data []byte) error {
   format := string(data)
   *m = func(id string, value int) string {
      oldnew := []string{"%", "%%"} // first
      if strings.Contains(format, "$Time$") {
         oldnew = append(oldnew, "$Time$", "%d")
      } else {
         oldnew = append(oldnew,
            "$Number$", "%d",
            "$Number%02d$", "%02d",
            "$Number%03d$", "%03d",
            "$Number%04d$", "%04d",
            "$Number%05d$", "%05d",
            "$Number%06d$", "%06d",
            "$Number%07d$", "%07d",
            "$Number%08d$", "%08d",
            "$Number%09d$", "%09d",
         )
      }
      oldnew = append(oldnew, "$RepresentationID$", id) // last
      return fmt.Sprintf(strings.NewReplacer(oldnew...).Replace(format), value)
   }
   return nil
}

type Media func(string, int) string

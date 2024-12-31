package dash

import (
   "strings"
   "time"
)

func (d *Duration) UnmarshalText(data []byte) error {
   var err error
   d.Duration, err = time.ParseDuration(strings.ToLower(
      strings.TrimPrefix(string(data), "PT"),
   ))
   if err != nil {
      return err
   }
   return nil
}

type Duration struct {
   Duration time.Duration
}

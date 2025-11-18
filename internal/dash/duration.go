package dash

import (
   "fmt"
   "regexp"
   "strconv"
   "strings"
)

// This regex is a simplified parser for ISO 8601 durations of the format P...T...
// It handles hours (H), minutes (M), and seconds (S). It does not handle
// years, months, or days, as their duration is not constant.
var durationRegex = regexp.MustCompile(`^PT(?:([0-9.]+)H)?(?:([0-9.]+)M)?(?:([0-9.]+)S)?$`)

// AsSeconds parses an ISO 8601 duration string (e.g., "PT1M30.5S") and returns
// the total duration in seconds. It returns an error for invalid formats or
// formats that include Years, Months, or Days.
func (p *Period) AsSeconds() (float64, error) {
   if !strings.HasPrefix(p.Duration, "PT") {
      return 0, fmt.Errorf("invalid duration format: missing 'PT' prefix: %s", p.Duration)
   }

   matches := durationRegex.FindStringSubmatch(p.Duration)
   if matches == nil {
      return 0, fmt.Errorf("invalid or unsupported duration format: %s", p.Duration)
   }

   var totalSeconds float64

   // matches[1] is for Hours
   if matches[1] != "" {
      h, err := strconv.ParseFloat(matches[1], 64)
      if err != nil {
         return 0, fmt.Errorf("invalid hours in duration: %w", err)
      }
      totalSeconds += h * 3600
   }
   // matches[2] is for Minutes
   if matches[2] != "" {
      m, err := strconv.ParseFloat(matches[2], 64)
      if err != nil {
         return 0, fmt.Errorf("invalid minutes in duration: %w", err)
      }
      totalSeconds += m * 60
   }
   // matches[3] is for Seconds
   if matches[3] != "" {
      s, err := strconv.ParseFloat(matches[3], 64)
      if err != nil {
         return 0, fmt.Errorf("invalid seconds in duration: %w", err)
      }
      totalSeconds += s
   }

   return totalSeconds, nil
}

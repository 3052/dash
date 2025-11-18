package dash

import (
   "fmt"
   "strconv"
   "strings"
)

// AsSeconds parses an ISO 8601 duration string (e.g., "PT1M30.5S") and returns
// the total duration in seconds. It returns an error for invalid formats or
// formats that include Years, Months, or Days (Y, M, D).
func (p *Period) AsSeconds() (float64, error) {
   durationStr := p.Duration
   if !strings.HasPrefix(durationStr, "P") {
      return 0, fmt.Errorf("invalid duration format: must start with 'P': %s", durationStr)
   }

   // Find the time separator 'T'
   tIndex := strings.Index(durationStr, "T")

   // If 'T' is not found, we reject the format.
   // This correctly rejects date-only formats like "P1D" and invalid formats like "P10S".
   if tIndex == -1 {
      if strings.ContainsAny(durationStr, "YMD") {
         return 0, fmt.Errorf("unsupported duration format: date part is not supported: %s", durationStr)
      }
      return 0, fmt.Errorf("invalid duration format: missing 'T': %s", durationStr)
   }

   // Check for a date part. Anything between 'P' and 'T' is a date part, which we don't support.
   // The length of "P" is 1. If 'T' is at an index greater than 1, there's a date part.
   if tIndex > 1 {
      return 0, fmt.Errorf("unsupported duration format: date part is not supported: %s", durationStr)
   }

   timePart := durationStr[tIndex+1:]
   if timePart == "" {
      return 0, nil // A duration of "PT" is 0 seconds.
   }

   var totalSeconds float64
   var currentVal string

   for _, char := range timePart {
      switch char {
      case 'H':
         val, err := strconv.ParseFloat(currentVal, 64)
         if err != nil || currentVal == "" {
            return 0, fmt.Errorf("invalid value for hours in duration: %s", durationStr)
         }
         totalSeconds += val * 3600
         currentVal = ""
      case 'M':
         val, err := strconv.ParseFloat(currentVal, 64)
         if err != nil || currentVal == "" {
            return 0, fmt.Errorf("invalid value for minutes in duration: %s", durationStr)
         }
         totalSeconds += val * 60
         currentVal = ""
      case 'S':
         val, err := strconv.ParseFloat(currentVal, 64)
         if err != nil || currentVal == "" {
            return 0, fmt.Errorf("invalid value for seconds in duration: %s", durationStr)
         }
         totalSeconds += val
         currentVal = ""
      case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
         currentVal += string(char)
      default:
         return 0, fmt.Errorf("invalid character in duration: '%c' in %s", char, durationStr)
      }
   }

   // Check if there's a dangling number without a unit (e.g., "PT10")
   if currentVal != "" {
      return 0, fmt.Errorf("trailing number without unit in duration: %s", durationStr)
   }

   return totalSeconds, nil
}

package dash

import (
   "fmt"
   "net/url"
   "strconv"
   "strings"
   "time"
)

// Period represents a temporal part of the media content.
type Period struct {
   Duration       string           `xml:"duration,attr,omitempty"`
   ID             string           `xml:"id,attr,omitempty"`
   BaseURL        string           `xml:"BaseURL,omitempty"`
   AdaptationSets []*AdaptationSet `xml:"AdaptationSet"`

   // Navigation
   Parent *MPD `xml:"-"`
}

// ResolveBaseURL resolves the Period's BaseURL against the parent MPD's resolved BaseURL.
func (p *Period) ResolveBaseURL() (*url.URL, error) {
   parentBase, err := p.Parent.ResolveBaseURL()
   if err != nil {
      return nil, err
   }
   return resolveRef(parentBase, p.BaseURL)
}

// GetDuration parses the ISO 8601 Duration attribute.
// If Period@duration is missing, it falls back to MPD@mediaPresentationDuration.
func (p *Period) GetDuration() (time.Duration, error) {
   durStr := p.Duration
   if durStr == "" && p.Parent != nil {
      durStr = p.Parent.MediaPresentationDuration
   }

   if durStr == "" {
      return 0, nil
   }
   return parseIsoDuration(durStr)
}

func (p *Period) link() {
   for _, as := range p.AdaptationSets {
      // Req 10.1: AdaptationSet to Period
      as.Parent = p
      as.link()
   }
}

// parseIsoDuration parses a subset of ISO 8601 duration strings (P[n]Y[n]M[n]DT[n]H[n]M[n]S).
// It supports D (Days), H (Hours), M (Minutes), and S (Seconds).
func parseIsoDuration(iso string) (time.Duration, error) {
   if !strings.HasPrefix(iso, "P") {
      return 0, fmt.Errorf("invalid duration format: missing P prefix")
   }

   s := iso[1:] // Skip 'P'
   var total time.Duration

   // Split into Date and Time parts
   parts := strings.Split(s, "T")
   datePart := parts[0]
   timePart := ""
   if len(parts) > 1 {
      timePart = parts[1]
   }

   // Parse Date Part (Supports D for Days)
   if len(datePart) > 0 {
      d, err := parseDurationPart(datePart, map[byte]time.Duration{
         'D': 24 * time.Hour,
      })
      if err != nil {
         return 0, err
      }
      total += d
   }

   // Parse Time Part (H, M, S)
   if len(timePart) > 0 {
      t, err := parseDurationPart(timePart, map[byte]time.Duration{
         'H': time.Hour,
         'M': time.Minute,
         'S': time.Second,
      })
      if err != nil {
         return 0, err
      }
      total += t
   }

   return total, nil
}

func parseDurationPart(part string, units map[byte]time.Duration) (time.Duration, error) {
   var total time.Duration
   var numBuf string

   for i := 0; i < len(part); i++ {
      c := part[i]
      if (c >= '0' && c <= '9') || c == '.' {
         numBuf += string(c)
      } else {
         // Character is a unit designator
         if numBuf == "" {
            return 0, fmt.Errorf("missing value for unit %c", c)
         }
         val, err := strconv.ParseFloat(numBuf, 64)
         if err != nil {
            return 0, fmt.Errorf("invalid number %s: %v", numBuf, err)
         }
         numBuf = ""

         mult, ok := units[c]
         if !ok {
            return 0, fmt.Errorf("unsupported or invalid unit %c in component %s", c, part)
         }

         // Add to total (convert float seconds/minutes to nanoseconds)
         total += time.Duration(val * float64(mult))
      }
   }
   return total, nil
}

package dash

import (
   "net/url"
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

   // Simplified parsing logic using standard library
   return time.ParseDuration(strings.ToLower(strings.TrimPrefix(durStr, "PT")))
}

func (p *Period) link() {
   for _, as := range p.AdaptationSets {
      // Req 10.1: AdaptationSet to Period
      as.Parent = p
      as.link()
   }
}

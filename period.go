package dash

import (
   "net/url"
   "strings"
   "time"
)

// Period represents a temporal part of the media content.
type Period struct {
   Duration       string           `xml:"duration,attr"`
   Id             string           `xml:"id,attr"`
   BaseUrl        string           `xml:"BaseURL"`
   AdaptationSets []*AdaptationSet `xml:"AdaptationSet"`
   // Navigation
   Parent *Mpd `xml:"-"`
}

// ResolveBaseUrl resolves the Period's BaseURL against the parent Mpd's resolved BaseUrl.
func (p *Period) ResolveBaseUrl() (*url.URL, error) {
   parentBase, err := p.Parent.ResolveBaseUrl()
   if err != nil {
      return nil, err
   }
   return resolveRef(parentBase, p.BaseUrl)
}

// GetDuration parses the ISO 8601 Duration attribute.
func (p *Period) GetDuration() (time.Duration, error) {
   durStr := p.Duration
   if durStr == "" && p.Parent != nil {
      durStr = p.Parent.MediaPresentationDuration
   }
   if durStr == "" {
      return 0, nil
   }
   return time.ParseDuration(strings.ToLower(strings.TrimPrefix(durStr, "PT")))
}

func (p *Period) link() {
   for _, currentSet := range p.AdaptationSets {
      currentSet.Parent = p
      currentSet.link()
   }
}

package dash

import "net/url"

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

func (p *Period) link() {
   for _, as := range p.AdaptationSets {
      // Req 10.1: AdaptationSet to Period
      as.Parent = p
      as.link()
   }
}

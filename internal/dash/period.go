package dash

// Period represents a temporal part of the media content.
type Period struct {
   Duration       string           `xml:"duration,attr,omitempty"`
   ID             string           `xml:"id,attr,omitempty"`
   BaseURL        string           `xml:"BaseURL,omitempty"`
   AdaptationSets []*AdaptationSet `xml:"AdaptationSet"`

   // Navigation
   Parent *MPD `xml:"-"`
}

func (p *Period) link() {
   for _, as := range p.AdaptationSets {
      // Req 10.1: AdaptationSet to Period
      as.Parent = p
      as.link()
   }
}

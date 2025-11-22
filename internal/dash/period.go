package dash

// Period represents a distinct timing period within the media presentation.
type Period struct {
   Duration      string           `xml:"duration,attr,omitempty"`
   ID            string           `xml:"id,attr,omitempty"`
   BaseURL       string           `xml:"BaseURL,omitempty"`
   AdaptationSet []*AdaptationSet `xml:"AdaptationSet"`
}

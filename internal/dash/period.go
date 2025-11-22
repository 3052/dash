package dash

// Period represents the Period element.
type Period struct {
   ID             string          `xml:"id,attr,omitempty"`
   Duration       string          `xml:"duration,attr,omitempty"`
   BaseURL        string          `xml:"BaseURL,omitempty"`
   AdaptationSets []AdaptationSet `xml:"AdaptationSet,omitempty"`
}

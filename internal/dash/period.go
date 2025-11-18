package dash

import "encoding/xml"

// Period represents the Period element.
type Period struct {
   XMLName        xml.Name         `xml:"Period"`
   ID             string           `xml:"id,attr,omitempty"`
   Start          string           `xml:"start,attr,omitempty"`
   Duration       string           `xml:"duration,attr,omitempty"`
   AdaptationSets []*AdaptationSet `xml:"AdaptationSet"`
   BaseURL        string           `xml:"BaseURL,omitempty"`
}

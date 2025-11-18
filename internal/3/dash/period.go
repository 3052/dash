package dash

import "encoding/xml"

// Period describes a part of the content with a start time and duration.
type Period struct {
   XMLName        xml.Name        `xml:"Period"`
   ID             string          `xml:"id,attr,omitempty"`
   Start          string          `xml:"start,attr,omitempty"`
   Duration       string          `xml:"duration,attr,omitempty"`
   BaseURL        string          `xml:"BaseURL,omitempty"`
   AdaptationSets []AdaptationSet `xml:"AdaptationSet,omitempty"`
}

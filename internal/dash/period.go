package dash

import "encoding/xml"

// Period represents a Period element in the MPD.
type Period struct {
   XMLName        xml.Name         `xml:"Period"`
   AdaptationSets []*AdaptationSet `xml:"AdaptationSet,omitempty"`
}

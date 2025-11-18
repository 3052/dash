package dash

import "encoding/xml"

// Role specifies a media component role. It can be used to signal the purpose
// of the AdaptationSet, such as "main", "alternate", "caption", etc.
type Role struct {
   XMLName     xml.Name `xml:"Role"`
   SchemeIDURI string   `xml:"schemeIdUri,attr,omitempty"`
   Value       string   `xml:"value,attr,omitempty"`
   ID          string   `xml:"id,attr,omitempty"`
}

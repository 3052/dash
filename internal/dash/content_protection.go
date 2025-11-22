package dash

import "encoding/xml"

// ContentProtection specifies DRM scheme information.
type ContentProtection struct {
   XMLName     xml.Name `xml:"ContentProtection"`
   SchemeIdUri string   `xml:"schemeIdUri,attr,omitempty"`
   Value       string   `xml:"value,attr,omitempty"`
}

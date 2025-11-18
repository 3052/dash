package dash

import "encoding/xml"

// ContentProtection represents the ContentProtection element, which provides
// information about content encryption and DRM systems.
type ContentProtection struct {
   XMLName     xml.Name `xml:"ContentProtection"`
   SchemeIDURI string   `xml:"schemeIdUri,attr"`
   Value       string   `xml:"value,attr,omitempty"`
}

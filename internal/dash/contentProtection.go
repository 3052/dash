package dash

import "encoding/xml"

// ContentProtection represents the ContentProtection element, which provides
// information about content encryption and DRM systems.
type ContentProtection struct {
   XMLName     xml.Name `xml:"ContentProtection"`
   SchemeIDURI string   `xml:"schemeIdUri,attr"`
   Value       string   `xml:"value,attr,omitempty"`

   // PSSH is the Common Encryption PSSH box. The struct tag specifies the
   // full XML namespace ("urn:mpeg:cenc:2013") and the local element name ("pssh").
   PSSH *CencPSSH `xml:"urn:mpeg:cenc:2013 pssh,omitempty"`
}

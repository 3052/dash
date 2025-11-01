package dash

import "encoding/xml"

// ContentProtection represents a ContentProtection element in the MPD.
type ContentProtection struct {
   XMLName     xml.Name `xml:"ContentProtection"`
   SchemeIDURI *string  `xml:"schemeIdUri,attr,omitempty"`
   Value       *string  `xml:"value,attr,omitempty"`
   // Use the full namespace URI for the default_KID attribute
   DefaultKID *string `xml:"urn:mpeg:cenc:2013 default_KID,attr,omitempty"`
   // Use the full namespace URI for the pssh element
   PSSH *PSSH `xml:"urn:mpeg:cenc:2013 pssh,omitempty"`
}

// PSSH represents the cenc:pssh element.
type PSSH struct {
   // Define the XMLName with the full namespace URI and the local name
   XMLName xml.Name `xml:"urn:mpeg:cenc:2013 pssh"`
   Value   string   `xml:",innerxml"`
}

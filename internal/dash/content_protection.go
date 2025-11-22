package dash

// ContentProtection specifies DRM schemes.
type ContentProtection struct {
   SchemeIdUri string `xml:"schemeIdUri,attr,omitempty"`
   // cenc:pssh requires the specific namespace mapping
   Pssh string `xml:"urn:mpeg:cenc:2013 pssh,omitempty"`
}

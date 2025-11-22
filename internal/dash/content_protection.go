package dash

// ContentProtection represents the ContentProtection element.
type ContentProtection struct {
   SchemeIdUri string `xml:"schemeIdUri,attr,omitempty"`
   // Pssh handles the <cenc:pssh> element.
   // Note: We use the generic tag "pssh" here. If the XML uses namespaces,
   // encoding/xml usually matches the local name.
   Pssh string `xml:"pssh,omitempty"`
}

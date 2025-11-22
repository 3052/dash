package dash

// ContentProtection specifies information about the content protection schemes used.
type ContentProtection struct {
   SchemeIdUri string `xml:"schemeIdUri,attr,omitempty"`
   CencPssh    string `xml:"cenc:pssh,omitempty"`
}

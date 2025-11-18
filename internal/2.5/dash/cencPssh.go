package dash

// CencPSSH represents the cenc:pssh element. The struct tag on the parent
// ContentProtection's field handles the namespacing. This struct just
// captures the Base64 encoded data within the element.
type CencPSSH struct {
   Data string `xml:",chardata"`
}

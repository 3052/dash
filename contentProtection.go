package dash

type ContentProtection struct {
   Pssh Pssh `xml:"pssh"`
   SchemeIdUri SchemeIdUri `xml:"schemeIdUri,attr"`
}

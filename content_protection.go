package dash

import (
   "encoding/base64"
   "encoding/hex"
   "strings"
)

// ContentProtection specifies DRM schemes.
type ContentProtection struct {
   SchemeIdUri string `xml:"schemeIdUri,attr"`
   // DefaultKid maps to cenc:default_KID (urn:mpeg:cenc:2013)
   DefaultKid string `xml:"urn:mpeg:cenc:2013 default_KID,attr"`
   // Pssh maps to the cenc:pssh element (urn:mpeg:cenc:2013)
   Pssh string `xml:"urn:mpeg:cenc:2013 pssh"`
}

// GetDefaultKid returns the DefaultKid as a byte slice.
func (cp *ContentProtection) GetDefaultKid() ([]byte, error) {
   if cp.DefaultKid == "" {
      return nil, nil
   }
   clean := strings.ReplaceAll(cp.DefaultKid, "-", "")
   return hex.DecodeString(clean)
}

// GetPssh returns the PSSH data as a byte slice.
func (cp *ContentProtection) GetPssh() ([]byte, error) {
   if cp.Pssh == "" {
      return nil, nil
   }
   return base64.StdEncoding.DecodeString(cp.Pssh)
}

package dash

import (
   "encoding/base64"
   "encoding/hex"
   "strings"
)

// ContentProtection specifies DRM schemes.
type ContentProtection struct {
   SchemeIdUri string `xml:"schemeIdUri,attr"`
   // DefaultKID maps to cenc:default_KID (urn:mpeg:cenc:2013)
   DefaultKID string `xml:"urn:mpeg:cenc:2013 default_KID,attr"`
   // Pssh maps to the cenc:pssh element (urn:mpeg:cenc:2013)
   Pssh string `xml:"urn:mpeg:cenc:2013 pssh"`
}

// GetDefaultKID returns the DefaultKID as a byte slice.
// It strips hyphens from the UUID string before decoding the hex.
func (cp *ContentProtection) GetDefaultKID() ([]byte, error) {
   if cp.DefaultKID == "" {
      return nil, nil
   }
   clean := strings.ReplaceAll(cp.DefaultKID, "-", "")
   return hex.DecodeString(clean)
}

// GetPSSH returns the PSSH data as a byte slice.
// It decodes the Base64 content of the element.
func (cp *ContentProtection) GetPSSH() ([]byte, error) {
   if cp.Pssh == "" {
      return nil, nil
   }
   return base64.StdEncoding.DecodeString(cp.Pssh)
}

package dash

import (
   "encoding/base64"
   "encoding/hex"
   "errors"
   "strings"
)

func (p Pointer) contentProtection() []ContentProtection {
   if a := p.AdaptationSet; a.ContentProtection != nil {
      return a.ContentProtection
   }
   return p.Representation.ContentProtection
}

type ContentProtection struct {
   SchemeIdUri string `xml:"schemeIdUri,attr"`
   // this might not exist
   Default_KID string `xml:"default_KID,attr"`
   // this might not exist
   PSSH string `xml:"pssh"`
}

func (p Pointer) PSSH() ([]byte, error) {
   for _, c := range p.contentProtection() {
      if c.SchemeIdUri == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         if c.PSSH != "" {
            return base64.StdEncoding.DecodeString(c.PSSH)
         }
      }
   }
   return nil, errors.New("Pointer.PSSH")
}

func (p Pointer) Default_KID() ([]byte, error) {
   for _, c := range p.contentProtection() {
      if c.SchemeIdUri == "urn:mpeg:dash:mp4protection:2011" {
         c.Default_KID = strings.ReplaceAll(c.Default_KID, "-", "")
         return hex.DecodeString(c.Default_KID)
      }
   }
   return nil, errors.New("Pointer.Default_KID")
}

package dash

import (
   "strings"
   "time"
)

func (a *AdaptationSet) set(p *Period) {
   a.period = p
}

type AdaptationSet struct {
   Codecs            *string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height            *int64  `xml:"height,attr"`
   Lang              string  `xml:"lang,attr"`
   MimeType          *string `xml:"mimeType,attr"`
   Representation    []Representation
   Role              *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
   Width           *int64 `xml:"width,attr"`
   period          *Period
}

func (s SchemeIdUri) Widevine() bool {
   return s == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed"
}

type SchemeIdUri string

type ContentProtection struct {
   Pssh        Pssh        `xml:"pssh"`
   SchemeIdUri SchemeIdUri `xml:"schemeIdUri,attr"`
}

func (d *Duration) UnmarshalText(data []byte) error {
   var err error
   d.D, err = time.ParseDuration(strings.ToLower(
      strings.TrimPrefix(string(data), "PT"),
   ))
   if err != nil {
      return err
   }
   return nil
}

type Duration struct {
   D time.Duration
}

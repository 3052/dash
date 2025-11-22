package dash

import (
   "bytes"
   "encoding/xml"
)

// MPD represents the root element of the DASH Manifest.
type MPD struct {
   XMLName                   xml.Name `xml:"MPD"`
   MediaPresentationDuration string   `xml:"mediaPresentationDuration,attr,omitempty"`
   BaseURL                   string   `xml:"BaseURL,omitempty"`
   Period                    []*Period
}

// Parse parses a DASH MPD file from a byte slice.
func Parse(data []byte) (*MPD, error) {
   var m MPD
   reader := bytes.NewReader(data)
   decoder := xml.NewDecoder(reader)

   // Handle charset encodings if necessary, though standard utf-8 is default
   err := decoder.Decode(&m)
   if err != nil {
      return nil, err
   }

   return &m, nil
}

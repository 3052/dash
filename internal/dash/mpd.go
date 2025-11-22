package dash

import (
   "encoding/xml"
)

// Parse takes a byte slice of an MPD file and returns the parsed struct.
func Parse(input []byte) (*MPD, error) {
   var m MPD
   err := xml.Unmarshal(input, &m)
   if err != nil {
      return nil, err
   }
   return &m, nil
}

// MPD represents the root element.
type MPD struct {
   XMLName                   xml.Name `xml:"MPD"`
   MediaPresentationDuration string   `xml:"mediaPresentationDuration,attr,omitempty"`
   BaseURL                   string   `xml:"BaseURL,omitempty"`
   Periods                   []Period `xml:"Period,omitempty"`
}

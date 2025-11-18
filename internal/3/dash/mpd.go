package dash

import "encoding/xml"

// MPD represents the root element of a DASH Media Presentation Description.
type MPD struct {
   XMLName                   xml.Name `xml:"MPD"`
   Xmlns                     string   `xml:"xmlns,attr,omitempty"`
   Profiles                  string   `xml:"profiles,attr,omitempty"`
   Type                      string   `xml:"type,attr,omitempty"`
   MinBufferTime             string   `xml:"minBufferTime,attr"`
   MediaPresentationDuration string   `xml:"mediaPresentationDuration,attr,omitempty"`
   MaxSegmentDuration        string   `xml:"maxSegmentDuration,attr,omitempty"`
   BaseURL                   string   `xml:"BaseURL,omitempty"`
   Periods                   []Period `xml:"Period,omitempty"`
}

// Parse parses a DASH MPD from a byte slice.
func Parse(data []byte) (*MPD, error) {
   var m MPD
   err := xml.Unmarshal(data, &m)
   if err != nil {
      return nil, err
   }
   return &m, nil
}

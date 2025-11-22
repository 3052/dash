package dash

import (
   "encoding/xml"
)

// MPD represents the root element of the DASH MPD file.
// The XMLName field is omitted to avoid linter conflicts with child parent pointers.
// encoding/xml will automatically match the <MPD> element to this struct name.
type MPD struct {
   MediaPresentationDuration string    `xml:"mediaPresentationDuration,attr,omitempty"`
   BaseURL                   string    `xml:"BaseURL,omitempty"` // Requirement: Single element
   Periods                   []*Period `xml:"Period"`
}

// Parse takes a byte slice of an MPD file, unmarshals it,
// and populates the navigation parent pointers.
func Parse(data []byte) (*MPD, error) {
   var m MPD
   err := xml.Unmarshal(data, &m)
   if err != nil {
      return nil, err
   }

   // Initialize navigation links
   m.link()

   return &m, nil
}

// link establishes the parent-child relationships for navigation.
func (m *MPD) link() {
   for _, p := range m.Periods {
      // Req 10.3: Period to MPD
      p.Parent = m
      p.link()
   }
}

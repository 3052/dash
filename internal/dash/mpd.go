package dash

import "encoding/xml"

// MPD represents the Media Presentation Description (MPD) element.
type MPD struct {
   XMLName                   xml.Name  `xml:"MPD"`
   Type                      string    `xml:"type,attr"`
   MinBufferTime             string    `xml:"minBufferTime,attr"`
   MediaPresentationDuration string    `xml:"mediaPresentationDuration,attr"`
   Profiles                  string    `xml:"profiles,attr"`
   Periods                   []*Period `xml:"Period"`
}

// RepresentationsByID returns a map of all Representations in the MPD, keyed by their ID.
// Since IDs are not guaranteed to be unique across Periods, the value is a slice
// containing all representations that share the same ID.
func (m *MPD) RepresentationsByID() map[string][]*Representation {
   reps := make(map[string][]*Representation)
   for _, p := range m.Periods {
      for _, as := range p.AdaptationSets {
         for _, r := range as.Representations {
            if r.ID != "" {
               reps[r.ID] = append(reps[r.ID], r)
            }
         }
      }
   }
   return reps
}

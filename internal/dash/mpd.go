package dash

import (
   "encoding/xml"
)

// MPD represents the Media Presentation Description (MPD) element.
type MPD struct {
   XMLName                   xml.Name  `xml:"MPD"`
   Type                      string    `xml:"type,attr"`
   MinBufferTime             string    `xml:"minBufferTime,attr"`
   MediaPresentationDuration string    `xml:"mediaPresentationDuration,attr"`
   Profiles                  string    `xml:"profiles,attr"`
   Periods                   []*Period `xml:"Period"`
}

// QualityOptions returns a map where each key is a Representation ID. The value
// is a slice of Quality structs, where each struct contains the Representation
// itself plus the inherited context (like Lang and ContentType) from its
// parent AdaptationSet. This handles ID collisions across Periods.
func (m *MPD) QualityOptions() map[string][]*Quality {
   if m == nil {
      return nil
   }
   options := make(map[string][]*Quality)
   for _, p := range m.Periods {
      if p == nil {
         continue
      }
      for _, as := range p.AdaptationSets {
         if as == nil {
            continue
         }
         for _, r := range as.Representations {
            if r != nil && r.ID != "" {
               quality := &Quality{
                  Representation: r,
                  Lang:           as.Lang,
                  ContentType:    as.ContentType,
               }
               options[r.ID] = append(options[r.ID], quality)
            }
         }
      }
   }
   return options
}

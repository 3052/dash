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

// QualityOptions returns a map where each key is a unique Representation ID.
// The value is a single Quality object that contains the shared Representation
// data and a slice of all the different contexts (Period/AdaptationSet pairs)
// in which that Representation ID appears.
func (m *MPD) QualityOptions() map[string]*Quality {
   if m == nil {
      return nil
   }
   options := make(map[string]*Quality)
   for _, p := range m.Periods {
      if p == nil {
         continue
      }
      for _, as := range p.AdaptationSets {
         if as == nil {
            continue
         }
         for _, r := range as.Representations {
            if r == nil || r.ID == "" {
               continue
            }

            context := &RepresentationContext{
               Period:        p,
               AdaptationSet: as,
            }

            // If we've seen this ID before, just add the new context.
            if existingQuality, ok := options[r.ID]; ok {
               existingQuality.Contexts = append(existingQuality.Contexts, context)
            } else {
               // If this is the first time, create the new Quality object.
               options[r.ID] = &Quality{
                  Representation: r,
                  Contexts:       []*RepresentationContext{context},
               }
            }
         }
      }
   }
   return options
}

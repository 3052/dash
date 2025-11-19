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

// GetRepresentations returns all Representations wrapped in their context,
// grouped by the Representation ID.
func (m *MPD) GetRepresentations() map[string][]RepresentationContext {
   results := make(map[string][]RepresentationContext)

   // Iterate using indices to keep stable pointers to the slice elements
   for i := range m.Periods {
      period := &m.Periods[i]
      pCtx := PeriodContext{
         Period: period,
         MPD:    m,
      }

      for j := range period.AdaptationSets {
         as := &period.AdaptationSets[j]
         asCtx := AdaptationSetContext{
            AdaptationSet: as,
            Context:       pCtx,
         }

         for k := range as.Representations {
            rep := &as.Representations[k]
            rCtx := RepresentationContext{
               Representation: rep,
               Context:        asCtx,
            }

            results[rep.ID] = append(results[rep.ID], rCtx)
         }
      }
   }

   return results
}

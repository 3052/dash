package dash

import (
   "encoding/xml"
   "net/url"
   "strconv"
)

// Parse takes a byte slice of an MPD file, unmarshals it,
// links navigation parents, and normalizes Representation IDs.
func Parse(data []byte) (*Mpd, error) {
   var m Mpd
   err := xml.Unmarshal(data, &m)
   if err != nil {
      return nil, err
   }

   m.link()
   m.normalizeIds()

   return &m, nil
}

// Mpd represents the root element of the DASH MPD file.
type Mpd struct {
   MediaPresentationDuration string    `xml:"mediaPresentationDuration,attr"`
   BaseUrl                   string    `xml:"BaseURL"`
   Periods                   []*Period `xml:"Period"`
   MpdUrl                    *url.URL  `xml:"-"`
}

// ResolveBaseUrl resolves the MPD's BaseURL against the MpdUrl.
func (m *Mpd) ResolveBaseUrl() (*url.URL, error) {
   return resolveRef(m.MpdUrl, m.BaseUrl)
}

// GetRepresentations returns a map of all Representations keyed by their Id.
func (m *Mpd) GetRepresentations() map[string][]*Representation {
   grouped := make(map[string][]*Representation)
   for _, p := range m.Periods {
      for _, as := range p.AdaptationSets {
         for _, r := range as.Representations {
            grouped[r.Id] = append(grouped[r.Id], r)
         }
      }
   }
   return grouped
}

func (m *Mpd) normalizeIds() {
   reservedIds := make(map[string]bool)

   for _, p := range m.Periods {
      for _, as := range p.AdaptationSets {
         for _, r := range as.Representations {
            if r.requiresOriginalId() {
               reservedIds[r.Id] = true
            }
         }
      }
   }

   counter := 0
   patternToId := make(map[string]string)

   for _, p := range m.Periods {
      for _, as := range p.AdaptationSets {
         for _, r := range as.Representations {
            if r.requiresOriginalId() {
               continue
            }

            pattern := r.GetSegmentTemplate().Media
            if existingId, ok := patternToId[pattern]; ok {
               r.Id = existingId
               continue
            }

            var newId string
            for {
               newId = strconv.Itoa(counter)
               counter++
               if !reservedIds[newId] {
                  break
               }
            }

            patternToId[pattern] = newId
            r.Id = newId
         }
      }
   }
}

func (m *Mpd) link() {
   for _, p := range m.Periods {
      p.Parent = m
      p.link()
   }
}

func resolveRef(base *url.URL, relStr string) (*url.URL, error) {
   if relStr == "" {
      return base, nil
   }
   rel, err := url.Parse(relStr)
   if err != nil {
      return nil, err
   }
   if base == nil {
      return rel, nil
   }
   return base.ResolveReference(rel), nil
}

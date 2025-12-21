package dash

import (
   "encoding/xml"
   "net/url"
   "strconv"
)

// Parse takes a byte slice of an MPD file, unmarshals it,
// links navigation parents, and normalizes Representation IDs.
func Parse(data []byte) (*Mpd, error) {
   var manifest Mpd
   err := xml.Unmarshal(data, &manifest)
   if err != nil {
      return nil, err
   }
   manifest.link()
   manifest.normalizeIds()
   return &manifest, nil
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
   for _, manifestPeriod := range m.Periods {
      for _, currentSet := range manifestPeriod.AdaptationSets {
         for _, mediaRep := range currentSet.Representations {
            grouped[mediaRep.Id] = append(grouped[mediaRep.Id], mediaRep)
         }
      }
   }
   return grouped
}

// normalizeIds iterates through the MPD and rewrites Representation IDs.
func (m *Mpd) normalizeIds() {
   reservedIds := make(map[string]bool)
   for _, manifestPeriod := range m.Periods {
      for _, currentSet := range manifestPeriod.AdaptationSets {
         for _, mediaRep := range currentSet.Representations {
            if mediaRep.requiresOriginalId() {
               reservedIds[mediaRep.Id] = true
            }
         }
      }
   }
   counter := 0
   patternToId := make(map[string]string)
   for _, manifestPeriod := range m.Periods {
      for _, currentSet := range manifestPeriod.AdaptationSets {
         for _, mediaRep := range currentSet.Representations {
            if mediaRep.requiresOriginalId() {
               continue
            }
            currentTemplate := mediaRep.GetSegmentTemplate()
            if currentTemplate == nil {
               continue
            }
            pattern := currentTemplate.Media
            if existingId, ok := patternToId[pattern]; ok {
               mediaRep.Id = existingId
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
            mediaRep.Id = newId
         }
      }
   }
}

func (m *Mpd) link() {
   for _, manifestPeriod := range m.Periods {
      manifestPeriod.Parent = m
      manifestPeriod.link()
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

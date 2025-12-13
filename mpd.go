package dash

import (
   "encoding/xml"
   "net/url"
   "strconv"
)

// Parse takes a byte slice of an MPD file, unmarshals it,
// links navigation parents, and normalizes Representation IDs.
func Parse(data []byte) (*MPD, error) {
   var m MPD
   err := xml.Unmarshal(data, &m)
   if err != nil {
      return nil, err
   }

   // 1. Initialize navigation links
   m.link()

   // 2. Normalize IDs to simplified sequence numbers
   // (only where safe to do so without breaking URLs)
   m.normalizeIDs()

   return &m, nil
}

// MPD represents the root element of the DASH MPD file.
type MPD struct {
   MediaPresentationDuration string    `xml:"mediaPresentationDuration,attr"`
   BaseURL                   string    `xml:"BaseURL"`
   Periods                   []*Period `xml:"Period"`
   MPDURL                    *url.URL  `xml:"-"`
}

// ResolveBaseURL resolves the MPD's BaseURL against the MPDURL.
func (m *MPD) ResolveBaseURL() (*url.URL, error) {
   return resolveRef(m.MPDURL, m.BaseURL)
}

// GetRepresentations returns a map of all Representations keyed by their ID.
// Because Parse() calls normalizeIDs(), these IDs are guaranteed to be
// simplified ("0", "1") where possible, or the original IDs ("video_1") otherwise.
func (m *MPD) GetRepresentations() map[string][]*Representation {
   grouped := make(map[string][]*Representation)
   for _, p := range m.Periods {
      for _, as := range p.AdaptationSets {
         for _, r := range as.Representations {
            // We group simply by ID, relying on normalizeIDs to have
            // unified the IDs across periods for us.
            grouped[r.ID] = append(grouped[r.ID], r)
         }
      }
   }
   return grouped
}

// normalizeIDs iterates through the MPD and rewrites Representation IDs
// to simple counters ("0", "1", "2") IF the URL generation templates
// do not strictly require the original ID.
func (m *MPD) normalizeIDs() {
   counter := 0
   // assigned maps the unique continuity pattern/string -> the new short ID ("0", "1")
   assigned := make(map[string]string)

   for _, p := range m.Periods {
      for _, as := range p.AdaptationSets {
         for _, r := range as.Representations {
            // Check if the SegmentTemplate relies on the raw ID
            // (e.g. segments_$RepresentationID$.m4s)
            if r.requiresOriginalID() {
               // We cannot rename this safely. Keep original ID.
               // We add it to 'assigned' to ensure we don't accidentally reuse this ID string.
               assigned[r.ID] = r.ID
               continue
            }

            // Calculate the continuity pattern.
            // Since we ruled out ID-dependency above, this is usually strict template strings.
            pattern := r.getContinuityPattern()

            // Have we assigned a short ID to this pattern yet?
            newID, exists := assigned[pattern]
            if !exists {
               newID = strconv.Itoa(counter)
               assigned[pattern] = newID
               counter++
            }

            // Overwrite the ID with the simplified one
            r.ID = newID
         }
      }
   }
}

func (m *MPD) link() {
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

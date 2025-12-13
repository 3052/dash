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
// It detects collisions to ensure generated IDs do not clash with reserved IDs.
func (m *MPD) normalizeIDs() {
   // reservedIDs tracks IDs that cannot be changed and thus occupy that name.
   reservedIDs := make(map[string]bool)

   // Pass 1: Identify IDs that must be preserved
   for _, p := range m.Periods {
      for _, as := range p.AdaptationSets {
         for _, r := range as.Representations {
            if r.requiresOriginalID() {
               reservedIDs[r.ID] = true
            }
         }
      }
   }

   counter := 0
   // patternToID maps the unique continuity pattern -> the new simplified ID
   patternToID := make(map[string]string)

   // Pass 2: Assign new IDs to renamable representations
   for _, p := range m.Periods {
      for _, as := range p.AdaptationSets {
         for _, r := range as.Representations {
            if r.requiresOriginalID() {
               continue
            }

            pattern := r.getContinuityPattern()

            // Have we assigned an ID to this pattern (stream) yet?
            if existingID, ok := patternToID[pattern]; ok {
               r.ID = existingID
               continue
            }

            // Find the next available numeric ID that isn't reserved
            var newID string
            for {
               newID = strconv.Itoa(counter)
               counter++
               // If this number "0" is already taken by a preserved ID "0", skip it.
               if !reservedIDs[newID] {
                  break
               }
            }

            patternToID[pattern] = newID
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

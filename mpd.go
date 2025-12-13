package dash

import (
   "encoding/xml"
   "net/url"
)

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

// MPD represents the root element of the DASH MPD file.
// XMLName is omitted here to prevent SA5008 conflicts.
type MPD struct {
   MediaPresentationDuration string    `xml:"mediaPresentationDuration,attr"`
   BaseURL                   string    `xml:"BaseURL"`
   Periods                   []*Period `xml:"Period"`
   // MPDURL is the source URL of the MPD file itself.
   // It is used as the root for resolving relative BaseURLs.
   MPDURL *url.URL `xml:"-"`
}

// ResolveBaseURL resolves the MPD's BaseURL against the MPDURL.
func (m *MPD) ResolveBaseURL() (*url.URL, error) {
   // No parsing needed, MPDURL is already *url.URL
   return resolveRef(m.MPDURL, m.BaseURL)
}

// GetRepresentations returns a map of all Representations in the MPD,
// keyed by their continuity key (usually the ID, or a specific template string).
func (m *MPD) GetRepresentations() map[string][]*Representation {
   grouped := make(map[string][]*Representation)
   for _, p := range m.Periods {
      for _, as := range p.AdaptationSets {
         for _, r := range as.Representations {
            key := r.GetContinuityKey()
            grouped[key] = append(grouped[key], r)
         }
      }
   }
   return grouped
}

// link establishes the parent-child relationships for navigation.
func (m *MPD) link() {
   for _, p := range m.Periods {
      // Req 10.3: Period to MPD
      p.Parent = m
      p.link()
   }
}

// resolveRef is a helper shared within the package to resolve a relative URL string
// against a base *url.URL using RFC 3986 rules.
func resolveRef(base *url.URL, relStr string) (*url.URL, error) {
   // If the relative string is empty, return the base as is.
   if relStr == "" {
      return base, nil
   }
   rel, err := url.Parse(relStr)
   if err != nil {
      return nil, err
   }
   // Handle case where base is nil (e.g. MPDURL not set)
   if base == nil {
      return rel, nil
   }
   return base.ResolveReference(rel), nil
}

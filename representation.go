package dash

import (
   "fmt"
   "hash/crc32"
   "net/url"
   "strings"
)

// Representation describes a version of the media content.
type Representation struct {
   Bandwidth         int                  `xml:"bandwidth,attr"`
   Codecs            string               `xml:"codecs,attr"`
   Height            int                  `xml:"height,attr"`
   ID                string               `xml:"id,attr"`
   MimeType          string               `xml:"mimeType,attr"`
   Width             int                  `xml:"width,attr"`
   BaseURL           string               `xml:"BaseURL"`
   SegmentTemplate   *SegmentTemplate     `xml:"SegmentTemplate"`
   ContentProtection []*ContentProtection `xml:"ContentProtection"`
   SegmentBase       *SegmentBase         `xml:"SegmentBase"`
   SegmentList       *SegmentList         `xml:"SegmentList"`
   // Navigation
   Parent *AdaptationSet `xml:"-"`
}

// ResolveBaseURL resolves the Representation's BaseURL against the parent hierarchy.
func (r *Representation) ResolveBaseURL() (*url.URL, error) {
   parentBase, err := r.Parent.getAbsoluteBaseURL()
   if err != nil {
      return nil, err
   }
   return resolveRef(parentBase, r.BaseURL)
}

// GetContinuityKey returns the 32-bit CRC32 hash of the identity string.
// This hash is used to group Representations across Periods on a simplified logical timeline.
func (r *Representation) GetContinuityKey() uint32 {
   return crc32.ChecksumIEEE([]byte(r.getContinuityString()))
}

// getContinuityString returns the raw string used for grouping.
//
// Logic:
// 1. If using SegmentTemplate:
//   - If 'initialization' or 'media' contains "$RepresentationID$", the Representation.ID is sufficient for uniqueness.
//   - Otherwise, the template string itself ('initialization' or fallback to 'media') is the key.
//
// 2. All other cases (SegmentList, SegmentBase, simple BaseURL):
//   - The Representation.ID is used as the key.
func (r *Representation) getContinuityString() string {
   if st := r.GetSegmentTemplate(); st != nil {
      // If the template explicitly relies on the ID variable, the ID is the continuity key.
      if strings.Contains(st.Initialization, "$RepresentationID$") ||
         strings.Contains(st.Media, "$RepresentationID$") {
         return r.ID
      }

      // If the template does not use the ID variable, the template string itself defines the stream.
      if st.Initialization != "" {
         return st.Initialization
      }
      if st.Media != "" {
         return st.Media
      }
   }

   // For SegmentList, SegmentBase, or implicit BaseURL, we assume the Representation ID implies continuity.
   return r.ID
}

// GetCodecs returns the codecs for this Representation.
// If the Codecs attribute is empty on the Representation,
// it returns the Codecs attribute from the parent AdaptationSet.
func (r *Representation) GetCodecs() string {
   if r.Codecs != "" {
      return r.Codecs
   }
   if r.Parent != nil {
      return r.Parent.Codecs
   }
   return ""
}

// GetHeight returns the height for this Representation.
// If the Height attribute is 0 on the Representation,
// it returns the Height attribute from the parent AdaptationSet.
func (r *Representation) GetHeight() int {
   if r.Height != 0 {
      return r.Height
   }
   if r.Parent != nil {
      return r.Parent.Height
   }
   return 0
}

// GetWidth returns the width for this Representation.
// If the Width attribute is 0 on the Representation,
// it returns the Width attribute from the parent AdaptationSet.
func (r *Representation) GetWidth() int {
   if r.Width != 0 {
      return r.Width
   }
   if r.Parent != nil {
      return r.Parent.Width
   }
   return 0
}

// GetMimeType returns the mimeType for this Representation.
// If the MimeType attribute is empty on the Representation,
// it returns the MimeType attribute from the parent AdaptationSet.
func (r *Representation) GetMimeType() string {
   if r.MimeType != "" {
      return r.MimeType
   }
   if r.Parent != nil {
      return r.Parent.MimeType
   }
   return ""
}

// GetContentProtection returns the ContentProtection elements for this Representation.
// If the Representation has its own ContentProtection elements, they are returned.
// Otherwise, it returns the ContentProtection elements from the parent AdaptationSet.
func (r *Representation) GetContentProtection() []*ContentProtection {
   if len(r.ContentProtection) > 0 {
      return r.ContentProtection
   }
   if r.Parent != nil {
      return r.Parent.ContentProtection
   }
   return nil
}

// GetSegmentTemplate returns the SegmentTemplate for this Representation.
// If the SegmentTemplate is nil on the Representation,
// it returns the SegmentTemplate from the parent AdaptationSet.
func (r *Representation) GetSegmentTemplate() *SegmentTemplate {
   if r.SegmentTemplate != nil {
      return r.SegmentTemplate
   }
   if r.Parent != nil {
      return r.Parent.SegmentTemplate
   }
   return nil
}

// String returns a multi-line summary of the Representation.
// Fields: bandwidth, width, height, codecs, mimeType, lang, role, period, key (hex encoded hash).
// Optional fields are omitted if empty/zero.
func (r *Representation) String() string {
   var b []byte

   var periodID, lang, roleVal string
   if r.Parent != nil {
      lang = r.Parent.Lang
      if r.Parent.Role != nil {
         roleVal = r.Parent.Role.Value
      }
      if r.Parent.Parent != nil {
         periodID = r.Parent.Parent.ID
      }
   }

   // 1. Representation@bandwidth
   b = fmt.Appendf(b, "bandwidth = %d", r.Bandwidth)

   // 2. Representation.GetWidth
   if w := r.GetWidth(); w != 0 {
      b = fmt.Appendf(b, "\nwidth = %d", w)
   }

   // 3. Representation.GetHeight
   if h := r.GetHeight(); h != 0 {
      b = fmt.Appendf(b, "\nheight = %d", h)
   }

   // 4. Representation.GetCodecs
   if c := r.GetCodecs(); c != "" {
      b = fmt.Appendf(b, "\ncodecs = %s", c)
   }

   // 5. Representation.GetMimeType
   b = fmt.Appendf(b, "\nmimeType = %s", r.GetMimeType())

   // 6. AdaptationSet@lang
   if lang != "" {
      b = fmt.Appendf(b, "\nlang = %s", lang)
   }

   // 7. Role@value
   if roleVal != "" {
      b = fmt.Appendf(b, "\nrole = %s", roleVal)
   }

   // 8. Period@id
   if periodID != "" {
      b = fmt.Appendf(b, "\nperiod = %s", periodID)
   }

   // 9. Continuity Key (Hexadecimal Hash)
   // We print the hex value to ensure it is always safe for command line usage.
   b = fmt.Appendf(b, "\nkey = %x", r.GetContinuityKey())

   return string(b)
}

func (r *Representation) link() {
   if r.SegmentTemplate != nil {
      // Req 10.7: SegmentTemplate to Representation
      r.SegmentTemplate.ParentRepresentation = r
   }
   if r.SegmentList != nil {
      // Req 10.5: SegmentList to Representation
      r.SegmentList.Parent = r
      r.SegmentList.link()
   }
   if r.SegmentBase != nil {
      r.SegmentBase.link()
   }
}

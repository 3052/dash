package dash

import (
   "fmt"
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

// GetContentProtection returns a set of unique ContentProtection values
// from both the AdaptationSet and the Representation.
func (r *Representation) GetContentProtection() map[ContentProtection]struct{} {
   unique := make(map[ContentProtection]struct{})

   // 1. Add from Parent AdaptationSet
   if r.Parent != nil {
      for _, cp := range r.Parent.ContentProtection {
         if cp != nil {
            unique[*cp] = struct{}{}
         }
      }
   }

   // 2. Add from Self Representation
   for _, cp := range r.ContentProtection {
      if cp != nil {
         unique[*cp] = struct{}{}
      }
   }

   return unique
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
// Fields: bandwidth, width, height, codecs, mimeType, lang, role, period, id.
// Optional fields are omitted if empty/zero.
func (r *Representation) String() string {
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

   var parts []string

   // 1. Representation@bandwidth
   parts = append(parts, fmt.Sprintf("bandwidth = %d", r.Bandwidth))

   // 2. Representation.GetWidth
   if w := r.GetWidth(); w != 0 {
      parts = append(parts, fmt.Sprintf("width = %d", w))
   }

   // 3. Representation.GetHeight
   if h := r.GetHeight(); h != 0 {
      parts = append(parts, fmt.Sprintf("height = %d", h))
   }

   // 4. Representation.GetCodecs
   if c := r.GetCodecs(); c != "" {
      parts = append(parts, fmt.Sprintf("codecs = %s", c))
   }

   // 5. Representation.GetMimeType
   parts = append(parts, fmt.Sprintf("mimeType = %s", r.GetMimeType()))

   // 6. AdaptationSet@lang
   if lang != "" {
      parts = append(parts, fmt.Sprintf("lang = %s", lang))
   }

   // 7. Role@value
   if roleVal != "" {
      parts = append(parts, fmt.Sprintf("role = %s", roleVal))
   }

   // 8. Period@id
   if periodID != "" {
      parts = append(parts, fmt.Sprintf("period = %s", periodID))
   }

   // 9. Representation@id
   parts = append(parts, fmt.Sprintf("id = %s", r.ID))

   return strings.Join(parts, "\n")
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

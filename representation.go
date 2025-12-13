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
   Parent            *AdaptationSet       `xml:"-"`
}

// ResolveBaseURL resolves the Representation's BaseURL against the parent hierarchy.
func (r *Representation) ResolveBaseURL() (*url.URL, error) {
   parentBase, err := r.Parent.getAbsoluteBaseURL()
   if err != nil {
      return nil, err
   }
   return resolveRef(parentBase, r.BaseURL)
}

// requiresOriginalID checks if the SegmentTemplate strictly requires the
// original ID value to generate correct URLs.
func (r *Representation) requiresOriginalID() bool {
   st := r.GetSegmentTemplate()
   if st == nil {
      // If no template, we assume ID is non-functional info,
      // unless it's implicit in SegmentBase (rare to rely on specific ID string).
      // We safely allow renaming for basic BaseURL cases.
      return false
   }
   // If the template uses the $RepresentationID$ variable, we SHOULD NOT rename it.
   return strings.Contains(st.Initialization, "$RepresentationID$") ||
      strings.Contains(st.Media, "$RepresentationID$")
}

// getContinuityPattern returns the string used to group streams when IDs are ignored.
// This is typically the concatenation of the template fields.
func (r *Representation) getContinuityPattern() string {
   st := r.GetSegmentTemplate()
   if st != nil {
      // Use the template strings themselves as the unique signature
      return st.Initialization + "|" + st.Media
   }
   // Fallback for SegmentList/BaseURL: use the original ID as the pattern
   // before it gets overwritten.
   return r.ID
}

// GetCodecs returns the codecs for this Representation.
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
func (r *Representation) GetMimeType() string {
   if r.MimeType != "" {
      return r.MimeType
   }
   if r.Parent != nil {
      return r.Parent.MimeType
   }
   return ""
}

// GetContentProtection returns the ContentProtection elements.
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
// Since IDs are normalized in Parse(), 'id' will display the clean values ("0, 1") or original IDs.
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

   // 1. Bandwidth
   b = fmt.Appendf(b, "bandwidth = %d", r.Bandwidth)

   if w := r.GetWidth(); w != 0 {
      b = fmt.Appendf(b, "\nwidth = %d", w)
   }

   if h := r.GetHeight(); h != 0 {
      b = fmt.Appendf(b, "\nheight = %d", h)
   }

   if c := r.GetCodecs(); c != "" {
      b = fmt.Appendf(b, "\ncodecs = %s", c)
   }

   b = fmt.Appendf(b, "\nmimeType = %s", r.GetMimeType())

   if lang != "" {
      b = fmt.Appendf(b, "\nlang = %s", lang)
   }

   if roleVal != "" {
      b = fmt.Appendf(b, "\nrole = %s", roleVal)
   }

   if periodID != "" {
      b = fmt.Appendf(b, "\nperiod = %s", periodID)
   }

   // Last. ID (Normalized or Original)
   b = fmt.Appendf(b, "\nid = %s", r.ID)

   return string(b)
}

func (r *Representation) link() {
   if r.SegmentTemplate != nil {
      r.SegmentTemplate.ParentRepresentation = r
   }
   if r.SegmentList != nil {
      r.SegmentList.Parent = r
      r.SegmentList.link()
   }
   if r.SegmentBase != nil {
      r.SegmentBase.link()
   }
}

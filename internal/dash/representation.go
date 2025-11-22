package dash

import (
   "fmt"
   "net/url"
)

// Representation describes a version of the media content.
type Representation struct {
   Bandwidth         int                  `xml:"bandwidth,attr,omitempty"`
   Codecs            string               `xml:"codecs,attr,omitempty"`
   Height            int                  `xml:"height,attr,omitempty"`
   ID                string               `xml:"id,attr,omitempty"`
   MimeType          string               `xml:"mimeType,attr,omitempty"`
   Width             int                  `xml:"width,attr,omitempty"`
   BaseURL           string               `xml:"BaseURL,omitempty"`
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

// GetContentProtection returns the ContentProtection elements for this Representation.
// If the slice is empty on the Representation, it returns the slice from the parent AdaptationSet.
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

   return fmt.Sprintf(
      "Lang: %s\n"+
         "PeriodID: %s\n"+
         "Codecs: %s\n"+
         "Height: %d\n"+
         "MimeType: %s\n"+
         "Width: %d\n"+
         "Bandwidth: %d\n"+
         "Role: %s",
      lang,
      periodID,
      r.GetCodecs(),
      r.GetHeight(),
      r.GetMimeType(),
      r.GetWidth(),
      r.Bandwidth,
      roleVal,
   )
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

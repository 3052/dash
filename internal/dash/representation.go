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

// String returns a formatted summary of the Representation, including
// parent context (Period ID, Language, Role) and inherited attributes.
func (r *Representation) String() string {
   var periodID, lang, roleStr string

   if r.Parent != nil {
      lang = r.Parent.Lang

      if r.Parent.Role != nil {
         roleStr = r.Parent.Role.Value
      }

      if r.Parent.Parent != nil {
         periodID = r.Parent.Parent.ID
      }
   }

   return fmt.Sprintf("PeriodID=%s Lang=%s Role=%s Bandwidth=%d MimeType=%s Codecs=%s Width=%d Height=%d",
      periodID,
      lang,
      roleStr,
      r.Bandwidth,
      r.GetMimeType(),
      r.GetCodecs(),
      r.GetWidth(),
      r.GetHeight(),
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

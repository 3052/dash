package dash

import (
   "fmt"
   "net/url"
   "strings"
)

// String returns a multi-line summary of the Representation.
func (r *Representation) String() string {
   sb := &strings.Builder{}
   var periodId, lang, roleVal string
   if r.Parent != nil {
      lang = r.Parent.Lang
      if r.Parent.Role != nil {
         roleVal = r.Parent.Role.Value
      }
      if r.Parent.Parent != nil {
         periodId = r.Parent.Parent.Id
      }
   }
   sb.WriteString("bandwidth = ")
   fmt.Fprint(sb, r.Bandwidth)
   if w := r.GetWidth(); w != 0 {
      sb.WriteString("\nwidth = ")
      fmt.Fprint(sb, w)
   }
   if h := r.GetHeight(); h != 0 {
      sb.WriteString("\nheight = ")
      fmt.Fprint(sb, h)
   }
   if c := r.GetCodecs(); c != "" {
      sb.WriteString("\ncodecs = ")
      sb.WriteString(c)
   }
   sb.WriteString("\nmimeType = ")
   sb.WriteString(r.GetMimeType())
   if lang != "" {
      sb.WriteString("\nlang = ")
      sb.WriteString(lang)
   }
   if roleVal != "" {
      sb.WriteString("\nrole = ")
      sb.WriteString(roleVal)
   }
   if periodId != "" {
      sb.WriteString("\nperiod = ")
      sb.WriteString(periodId)
   }
   sb.WriteString("\nid = ")
   sb.WriteString(r.Id)
   return sb.String()
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

// Representation describes a version of the media content.
type Representation struct {
   Bandwidth         int                  `xml:"bandwidth,attr"`
   Codecs            string               `xml:"codecs,attr"`
   Height            int                  `xml:"height,attr"`
   Id                string               `xml:"id,attr"`
   MimeType          string               `xml:"mimeType,attr"`
   Width             int                  `xml:"width,attr"`
   BaseUrl           string               `xml:"BaseURL"`
   SegmentTemplate   *SegmentTemplate     `xml:"SegmentTemplate"`
   ContentProtection []*ContentProtection `xml:"ContentProtection"`
   SegmentBase       *SegmentBase         `xml:"SegmentBase"`
   SegmentList       *SegmentList         `xml:"SegmentList"`
   Parent            *AdaptationSet       `xml:"-"`
}

// ResolveBaseUrl resolves the Representation's BaseURL against the parent hierarchy.
func (r *Representation) ResolveBaseUrl() (*url.URL, error) {
   parentBase, err := r.Parent.getAbsoluteBaseUrl()
   if err != nil {
      return nil, err
   }
   return resolveRef(parentBase, r.BaseUrl)
}

// requiresOriginalId checks if the Representation ID should be preserved.
func (r *Representation) requiresOriginalId() bool {
   st := r.GetSegmentTemplate()
   if st == nil {
      return true
   }
   return strings.Contains(st.Media, "$RepresentationID$")
}

func (r *Representation) GetCodecs() string {
   if r.Codecs != "" {
      return r.Codecs
   }
   if r.Parent != nil {
      return r.Parent.Codecs
   }
   return ""
}

func (r *Representation) GetHeight() int {
   if r.Height != 0 {
      return r.Height
   }
   if r.Parent != nil {
      return r.Parent.Height
   }
   return 0
}

func (r *Representation) GetWidth() int {
   if r.Width != 0 {
      return r.Width
   }
   if r.Parent != nil {
      return r.Parent.Width
   }
   return 0
}

func (r *Representation) GetMimeType() string {
   if r.MimeType != "" {
      return r.MimeType
   }
   if r.Parent != nil {
      return r.Parent.MimeType
   }
   return ""
}

func (r *Representation) GetContentProtection() []*ContentProtection {
   if len(r.ContentProtection) > 0 {
      return r.ContentProtection
   }
   if r.Parent != nil {
      return r.Parent.ContentProtection
   }
   return nil
}

func (r *Representation) GetSegmentTemplate() *SegmentTemplate {
   if r.SegmentTemplate != nil {
      return r.SegmentTemplate
   }
   if r.Parent != nil {
      return r.Parent.SegmentTemplate
   }
   return nil
}

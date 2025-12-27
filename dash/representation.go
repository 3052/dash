package dash

import (
   "net/url"
   "strconv"
   "strings"
)

// String returns a multi-line summary of the Representation.
func (r *Representation) String() string {
   var data strings.Builder
   data.WriteString("bandwidth = ")
   data.WriteString(strconv.Itoa(r.Bandwidth))
   if width := r.GetWidth(); width != 0 {
      data.WriteString("\nwidth = ")
      data.WriteString(strconv.Itoa(width))
   }
   if height := r.GetHeight(); height != 0 {
      data.WriteString("\nheight = ")
      data.WriteString(strconv.Itoa(height))
   }
   if codecs := r.GetCodecs(); codecs != "" {
      data.WriteString("\ncodecs = ")
      data.WriteString(codecs)
   }
   data.WriteString("\nmimeType = ")
   data.WriteString(r.GetMimeType())
   if label := r.GetLabel(); label != "" {
      data.WriteString("\nlabel = ")
      data.WriteString(label)
   } else if lang := r.GetLang(); lang != "" {
      data.WriteString("\nlang = ")
      data.WriteString(lang)
   }
   if role := r.GetRole(); role != "" {
      data.WriteString("\nrole = ")
      data.WriteString(role)
   }
   if periodId := r.GetPeriodId(); periodId != "" {
      data.WriteString("\nperiod = ")
      data.WriteString(periodId)
   }
   data.WriteString("\nid = ")
   data.WriteString(r.Id)
   return data.String()
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
   currentTemplate := r.GetSegmentTemplate()
   if currentTemplate == nil {
      return true
   }
   return strings.Contains(currentTemplate.Media, "$RepresentationID$")
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

func (r *Representation) GetLang() string {
   if r.Parent != nil {
      return r.Parent.Lang
   }
   return ""
}

func (r *Representation) GetLabel() string {
   if r.Parent != nil {
      return r.Parent.Label
   }
   return ""
}

func (r *Representation) GetRole() string {
   if r.Parent != nil && r.Parent.Role != nil {
      return r.Parent.Role.Value
   }
   return ""
}

func (r *Representation) GetPeriodId() string {
   if r.Parent != nil && r.Parent.Parent != nil {
      return r.Parent.Parent.Id
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

func (r *Representation) link() {
   if r.SegmentTemplate != nil {
      r.SegmentTemplate.ParentRepresentation = r
   }
   if r.SegmentList != nil {
      r.SegmentList.Parent = r
      r.SegmentList.link()
   }
}

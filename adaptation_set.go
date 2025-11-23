package dash

import "net/url"

// AdaptationSet groups Representations.
type AdaptationSet struct {
   Codecs            string               `xml:"codecs,attr,omitempty"`
   Height            int                  `xml:"height,attr,omitempty"`
   Lang              string               `xml:"lang,attr,omitempty"`
   MimeType          string               `xml:"mimeType,attr,omitempty"`
   Width             int                  `xml:"width,attr,omitempty"`
   ContentProtection []*ContentProtection `xml:"ContentProtection"`
   Role              *Role                `xml:"Role"`
   SegmentTemplate   *SegmentTemplate     `xml:"SegmentTemplate"`
   Representations   []*Representation    `xml:"Representation"`

   // Navigation
   Parent *Period `xml:"-"`
}

// getAbsoluteBaseURL returns the resolved BaseURL of the parent Period.
func (as *AdaptationSet) getAbsoluteBaseURL() (*url.URL, error) {
   return as.Parent.ResolveBaseURL()
}

func (as *AdaptationSet) link() {
   if as.SegmentTemplate != nil {
      // Req 10.6: SegmentTemplate to AdaptationSet
      as.SegmentTemplate.ParentAdaptationSet = as
   }
   for _, rep := range as.Representations {
      // Req 10.4: Representation to AdaptationSet
      rep.Parent = as
      rep.link()
   }
}

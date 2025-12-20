package dash

import "net/url"

// AdaptationSet groups Representations.
type AdaptationSet struct {
   Codecs            string               `xml:"codecs,attr"`
   Height            int                  `xml:"height,attr"`
   Lang              string               `xml:"lang,attr"`
   MimeType          string               `xml:"mimeType,attr"`
   Width             int                  `xml:"width,attr"`
   ContentProtection []*ContentProtection `xml:"ContentProtection"`
   Role              *Role                `xml:"Role"`
   SegmentTemplate   *SegmentTemplate     `xml:"SegmentTemplate"`
   Representations   []*Representation    `xml:"Representation"`
   // Navigation
   Parent *Period `xml:"-"`
}

// getAbsoluteBaseUrl returns the resolved BaseUrl of the parent Period.
func (as *AdaptationSet) getAbsoluteBaseUrl() (*url.URL, error) {
   return as.Parent.ResolveBaseUrl()
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

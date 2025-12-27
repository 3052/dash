package dash

import "net/url"

// AdaptationSet groups Representations.
type AdaptationSet struct {
   Codecs            string               `xml:"codecs,attr"`
   ContentProtection []*ContentProtection `xml:"ContentProtection"`
   Height            int                  `xml:"height,attr"`
   Label             string               `xml:"Label"`
   Lang              string               `xml:"lang,attr"`
   MimeType          string               `xml:"mimeType,attr"`
   Representations   []*Representation    `xml:"Representation"`
   Role              *Role                `xml:"Role"`
   SegmentTemplate   *SegmentTemplate     `xml:"SegmentTemplate"`
   Width             int                  `xml:"width,attr"`
   // Navigation
   Parent *Period `xml:"-"`
}

// Role defines the role of the media content.
type Role struct {
   Value string `xml:"value,attr"`
}

// getAbsoluteBaseUrl returns the resolved BaseUrl of the parent Period.
func (as *AdaptationSet) getAbsoluteBaseUrl() (*url.URL, error) {
   return as.Parent.ResolveBaseUrl()
}

func (as *AdaptationSet) link() {
   if as.SegmentTemplate != nil {
      as.SegmentTemplate.ParentAdaptationSet = as
   }
   for _, mediaRep := range as.Representations {
      mediaRep.Parent = as
      mediaRep.link()
   }
}

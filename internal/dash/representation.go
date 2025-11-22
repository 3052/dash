package dash

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

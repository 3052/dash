package dash

type MPD struct {
   Period []struct {
      AdaptationSet []AdaptationSet
      ID string `xml:"id,attr"`
   }
}

func (m MPD) Some(f func(Representation) bool) {
   for _, period := range m.Period {
      for _, adapt := range period.AdaptationSet {
         for _, represent := range adapt.Representation {
            if represent.Codecs == "" {
               represent.Codecs = adapt.Codecs
            }
            if represent.ContentProtection == nil {
               represent.ContentProtection = adapt.ContentProtection
            }
            if represent.MimeType == "" {
               represent.MimeType = adapt.MimeType
            }
            if represent.SegmentTemplate == nil {
               represent.SegmentTemplate = adapt.SegmentTemplate
            }
            if !f(represent) {
               return
            }
         }
      }
   }
}

func (m MPD) Every(f func(Representation)) {
   m.Some(func(r Representation) bool {
      f(r)
      return true
   })
}

package dash

type Index struct {
   Period int
   AdaptationSet int
   Representation int
}

func (i Index) GetPeriod(m MPD) Period {
   return m.Period[i.Period]
}

func (i Index) GetAdaptation(m MPD) AdaptationSet {
   return i.GetPeriod(m).AdaptationSet[i.AdaptationSet]
}

func (i Index) GetRepresentation(m MPD) Representation {
   return i.GetAdaptation(m).Representation[i.Representation]
}

// media presentation description
// wikipedia.org/wiki/Dynamic_Adaptive_Streaming_over_HTTP
type MPD struct {
   Period []Period
}

type Period struct {
   AdaptationSet []AdaptationSet
   ID string `xml:"id,attr"`
}

type AdaptationSet struct {
   // this might be under Representation
   Codecs string `xml:"codecs,attr"`
   // this might be under Representation
   ContentProtection []ContentProtection
   // this might not exist
   Lang string `xml:"lang,attr"`
   // this might be under Representation
   MimeType string `xml:"mimeType,attr"`
   // pointer because we want to edit these
   Representation []Representation
   // this might not exist
   Role *struct {
      Value string `xml:"value,attr"`
   }
   // this might not exist, or might be under Representation
   SegmentTemplate *SegmentTemplate
}

type Representation struct {
   Bandwidth int `xml:"bandwidth,attr"`
   ID string `xml:"id,attr"`
   // this might not exist
   BaseURL string
   // this might be under AdaptationSet
   Codecs string `xml:"codecs,attr"`
   // this might be under AdaptationSet
   ContentProtection []ContentProtection
   // this might not exist
   Height *int `xml:"height,attr"`
   // this might be under AdaptationSet
   MimeType string `xml:"mimeType,attr"`
   // this might not exist
   SegmentBase *struct {
      Initialization struct {
         Range Range `xml:"range,attr"`
      }
      IndexRange Range `xml:"indexRange,attr"`
   }
   // this might not exist, or might be under AdaptationSet
   SegmentTemplate *SegmentTemplate
   // this might not exist
   Width *int `xml:"width,attr"`
}

func (i Index) Codecs(m MPD) string {
   if a := i.GetAdaptation(m); a.Codecs != "" {
      return a.Codecs
   }
   return i.GetRepresentation(m).Codecs
}

func (i Index) ContentProtection(m MPD) []ContentProtection {
   if a := i.GetAdaptation(m); a.ContentProtection != nil {
      return a.ContentProtection
   }
   return i.GetRepresentation(m).ContentProtection
}

func (i Index) MimeType(m MPD) string {
   if a := i.GetAdaptation(m); a.MimeType != "" {
      return a.MimeType
   }
   return i.GetRepresentation(m).MimeType
}

func (i Index) SegmentTemplate(m MPD) *SegmentTemplate {
   if a := i.GetAdaptation(m); a.SegmentTemplate != nil {
      return a.SegmentTemplate
   }
   return i.GetRepresentation(m).SegmentTemplate
}

func (m MPD) Some(f func(MPD, Index) bool) {
   for p, period := range m.Period {
      for a, adapt := range period.AdaptationSet {
         for r := range adapt.Representation {
            if !f(m, Index{p, a, r}) {
               return
            }
         }
      }
   }
}

func (m MPD) Every(f func(MPD, Index)) {
   m.Some(func(m MPD, i Index) bool {
      f(m, i)
      return true
   })
}

package dash

import "strings"

func (m MPD) GetAdaptation(i Index) AdaptationSet {
   return m.GetPeriod(i).AdaptationSet[i.AdaptationSet]
}

func (m MPD) GetRepresentation(i Index) Representation {
   return m.GetAdaptation(i).Representation[i.Representation]
}

func (m MPD) Codecs(i Index) string {
   if a := m.GetAdaptation(i); a.Codecs != "" {
      return a.Codecs
   }
   return m.GetRepresentation(i).Codecs
}

func (m MPD) MimeType(i Index) string {
   if a := m.GetAdaptation(i); a.MimeType != "" {
      return a.MimeType
   }
   return m.GetRepresentation(i).MimeType
}

func (m MPD) SegmentTemplate(i Index) *SegmentTemplate {
   if a := m.GetAdaptation(i); a.SegmentTemplate != nil {
      return a.SegmentTemplate
   }
   return m.GetRepresentation(i).SegmentTemplate
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

type Index struct {
   Period int
   AdaptationSet int
   Representation int
}

func (m MPD) Default_KID(i Index) (string, bool) {
   for _, cp := range m.ContentProtection(i) {
      if cp.SchemeIdUri == "urn:mpeg:dash:mp4protection:2011" {
         return strings.ReplaceAll(cp.Default_KID, "-", ""), true
      }
   }
   return "", false
}

func (m MPD) ContentProtection(i Index) []ContentProtection {
   if a := m.GetAdaptation(i); a.ContentProtection != nil {
      return a.ContentProtection
   }
   return m.GetRepresentation(i).ContentProtection
}

func (m MPD) GetPeriod(i Index) Period {
   return m.Period[i.Period]
}

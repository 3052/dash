package dash

import "strings"

func (m MPD) GetAdaptation(i Pointer) AdaptationSet {
   return m.GetPeriod(i).AdaptationSet[i.AdaptationSet]
}

func (m MPD) GetRepresentation(i Pointer) Representation {
   return m.GetAdaptation(i).Representation[i.Representation]
}

func (m MPD) Codecs(i Pointer) string {
   if a := m.GetAdaptation(i); a.Codecs != "" {
      return a.Codecs
   }
   return m.GetRepresentation(i).Codecs
}

func (m MPD) MimeType(i Pointer) string {
   if a := m.GetAdaptation(i); a.MimeType != "" {
      return a.MimeType
   }
   return m.GetRepresentation(i).MimeType
}

func (m MPD) SegmentTemplate(i Pointer) *SegmentTemplate {
   if a := m.GetAdaptation(i); a.SegmentTemplate != nil {
      return a.SegmentTemplate
   }
   return m.GetRepresentation(i).SegmentTemplate
}

func (m MPD) Some(f func(MPD, Pointer) bool) {
   for p, period := range m.Period {
      for a, adapt := range period.AdaptationSet {
         for r := range adapt.Representation {
            if !f(m, Pointer{p, a, r}) {
               return
            }
         }
      }
   }
}

func (m MPD) Every(f func(MPD, Pointer)) {
   m.Some(func(m MPD, i Pointer) bool {
      f(m, i)
      return true
   })
}

type Pointer struct {
   Period int
   AdaptationSet int
   Representation int
}

func (m MPD) Default_KID(i Pointer) (string, bool) {
   for _, cp := range m.ContentProtection(i) {
      if cp.SchemeIdUri == "urn:mpeg:dash:mp4protection:2011" {
         return strings.ReplaceAll(cp.Default_KID, "-", ""), true
      }
   }
   return "", false
}

func (m MPD) ContentProtection(i Pointer) []ContentProtection {
   if a := m.GetAdaptation(i); a.ContentProtection != nil {
      return a.ContentProtection
   }
   return m.GetRepresentation(i).ContentProtection
}

func (m MPD) GetPeriod(i Pointer) Period {
   return m.Period[i.Period]
}

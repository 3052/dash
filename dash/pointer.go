package dash

import "strings"

func (p Pointer) Codecs() string {
   if a := p.AdaptationSet; a.Codecs != "" {
      return a.Codecs
   }
   return p.Representation.Codecs
}

func (p Pointer) MimeType() string {
   if a := p.AdaptationSet; a.MimeType != "" {
      return a.MimeType
   }
   return p.Representation.MimeType
}

func (p Pointer) SegmentTemplate() *SegmentTemplate {
   if a := p.AdaptationSet; a.SegmentTemplate != nil {
      return a.SegmentTemplate
   }
   return p.Representation.SegmentTemplate
}

func (p Pointer) ContentProtection() []ContentProtection {
   if a := p.AdaptationSet; a.ContentProtection != nil {
      return a.ContentProtection
   }
   return p.Representation.ContentProtection
}

func (p Pointer) Default_KID() (string, bool) {
   for _, cp := range p.ContentProtection() {
      if cp.SchemeIdUri == "urn:mpeg:dash:mp4protection:2011" {
         return strings.ReplaceAll(cp.Default_KID, "-", ""), true
      }
   }
   return "", false
}

type Pointer struct {
   AdaptationSet *AdaptationSet
   Period *Period
   Representation *Representation
}

func (m MPD) Some(f func(Pointer) bool) {
   for _, period := range m.Period {
      for _, adapt := range period.AdaptationSet {
         for _, represent := range adapt.Representation {
            var p Pointer
            p.AdaptationSet = adapt
            p.Period = period
            p.Representation = represent
            if !f(p) {
               return
            }
         }
      }
   }
}

func (m MPD) Every(f func(Pointer)) {
   m.Some(func(p Pointer) bool {
      f(p)
      return true
   })
}

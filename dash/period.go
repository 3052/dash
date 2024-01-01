package dash

import (
   "encoding/xml"
   "io"
)

type Media struct {
   Period []struct {
      AdaptationSet []*AdaptationSet
      ID string `xml:"id,attr"`
   }
}

func (m *Media) Decode(r io.Reader) error {
   return xml.NewDecoder(r).Decode(m)
}

func (m Media) Representation(period string) ([]*Representation, error) {
   var rs []*Representation
   for _, p := range m.Period {
      if p.ID == period {
         for _, a := range p.AdaptationSet {
            for _, r := range a.Representation {
               if r.Codecs == "" {
                  r.Codecs = a.Codecs
               }
               if len(r.ContentProtection) == 0 {
                  r.ContentProtection = a.ContentProtection
               }
               if r.MimeType == "" {
                  r.MimeType = a.MimeType
               }
               if r.SegmentTemplate == nil {
                  r.SegmentTemplate = a.SegmentTemplate
               }
               r.adaptationSet = a
               rs = append(rs, r)
            }
         }
      }
   }
   return rs, nil
}

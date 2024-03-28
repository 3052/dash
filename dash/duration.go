package dash

import "encoding/xml"

func Unmarshal(b []byte) ([]Representation, error) {
   var s struct {
      Period []struct {
         AdaptationSet []adaptation_set
      }
   }
   err := xml.Unmarshal(b, &s)
   if err != nil {
      return nil, err
   }
   var rs []Representation
   for _, p := range s.Period {
      for _, a := range p.AdaptationSet {
         for _, r := range a.Representation {
            r.adaptation_set = &a
            rs = append(rs, r)
         }
      }
   }
   return rs, nil
}

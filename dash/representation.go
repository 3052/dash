package dash

import (
   "encoding/base64"
   "encoding/hex"
   "errors"
   "fmt"
   "strings"
)

func (r Representation) String() string {
   var b []byte
   if r.Width >= 1 {
      b = fmt.Append(b, "width: ", r.Width)
   }
   if r.Height >= 1 {
      if b != nil {
         b = append(b, '\n')
      }
      b = fmt.Append(b, "height: ", r.Height)
   }
   if r.Bandwidth >= 1 {
      if b != nil {
         b = append(b, '\n')
      }
      b = fmt.Append(b, "bandwidth: ", r.Bandwidth)
   }
   if r.Codecs != "" {
      if b != nil {
         b = append(b, '\n')
      }
      b = fmt.Append(b, "codecs: ", r.Codecs)
   }
   if b != nil {
      b = append(b, '\n')
   }
   b = fmt.Append(b, "type: ", r.MimeType)
   if v, ok := r.Role(); ok {
      b = fmt.Append(b, "\nrole: ", v)
   }
   if v, ok := r.Lang(); ok {
      b = fmt.Append(b, "\nlanguage: ", v)
   }
   b = fmt.Append(b, "\nid: ", r.ID)
   return string(b)
}

type Representation struct {
   Bandwidth int `xml:"bandwidth,attr"`
   ID string `xml:"id,attr"`
   adaptationSet *AdaptationSet
   // this might not exist
   BaseURL string
   // this might be under AdaptationSet
   Codecs string `xml:"codecs,attr"`
   // this might be under AdaptationSet
   ContentProtection []ContentProtection
   // this might not exist
   Height int `xml:"height,attr"`
   // this might be under AdaptationSet
   MimeType string `xml:"mimeType,attr"`
   // this might not exist
   SegmentBase *struct {
      IndexRange string `xml:"indexRange,attr"`
   }
   // this might not exist, or might be under AdaptationSet
   SegmentTemplate *SegmentTemplate
   // this might not exist
   Width int `xml:"width,attr"`
}

func (r Representation) Media() ([]string, bool) {
   t := r.SegmentTemplate
   if t == nil {
      return nil, false
   }
   var media []string
   for _, segment := range t.SegmentTimeline.S {
      for segment.R >= 0 {
         number := fmt.Sprint(t.StartNumber)
         medium := strings.Replace(t.Media, "$Number$", number, 1)
         medium = strings.Replace(medium, "$RepresentationID$", r.ID, 1)
         media = append(media, medium)
         segment.R--
         t.StartNumber++
      }
   }
   return media, true
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

func (r Representation) Sidx_Moof() (uint32, uint32, error) {
   if r.SegmentBase == nil {
      return 0, 0, errors.New("SegmentBase")
   }
   var start, end uint32
   _, err := fmt.Sscanf(r.SegmentBase.IndexRange, "%v-%v", &start, &end)
   if err != nil {
      return 0, 0, err
   }
   return start, end+1, nil
}

func (r Representation) Default_KID() ([]byte, error) {
   for _, c := range r.ContentProtection {
      if c.SchemeIdUri == "urn:mpeg:dash:mp4protection:2011" {
         c.Default_KID = strings.ReplaceAll(c.Default_KID, "-", "")
         return hex.DecodeString(c.Default_KID)
      }
   }
   return nil, errors.New("default_KID")
}

func (r Representation) Initialization() (string, bool) {
   if v := r.SegmentTemplate; v != nil {
      if v := v.Initialization; v != "" {
         return strings.Replace(v, "$RepresentationID$", r.ID, 1), true
      }
   }
   return "", false
}

func (r Representation) PSSH() ([]byte, error) {
   for _, c := range r.ContentProtection {
      if c.SchemeIdUri == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         return base64.StdEncoding.DecodeString(c.PSSH)
      }
   }
   return nil, errors.New("PSSH")
}

func (r Representation) Role() (string, bool) {
   if r := r.adaptationSet.Role; r != nil {
      return r.Value, true
   }
   return "", false
}

func (r Representation) Lang() (string, bool) {
   if v := r.adaptationSet.Lang; v != "" {
      return v, true
   }
   return "", false
}

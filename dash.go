package dash

import (
   "encoding/base64"
   "encoding/hex"
   "encoding/xml"
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
   if v := r.adaptationSet.Role; v != nil {
      b = fmt.Append(b, "\nrole: ", v.Value)
   }
   if r.adaptationSet.Lang != "" {
      b = fmt.Append(b, "\nlanguage: ", r.adaptationSet.Lang)
   }
   return string(b)
}

func Representations(b []byte) ([]*Representation, error) {
   var s struct {
      Period struct {
         // this need to be pointer so we can avoid loop bug
         AdaptationSet []*AdaptationSet
      }
   }
   err := xml.Unmarshal(b, &s)
   if err != nil {
      return nil, err
   }
   var rs []*Representation
   for _, a := range s.Period.AdaptationSet {
      for _, r := range a.Representation {
         if len(r.ContentProtection) == 0 {
            r.ContentProtection = a.ContentProtection
         }
         if r.MimeType == "" {
            r.MimeType = a.MimeType
         }
         r.adaptationSet = a
         rs = append(rs, r)
      }
   }
   return rs, nil
}

func (r Representation) Video() bool {
   return r.MimeType == "video/mp4"
}

func (r Representation) Audio() bool {
   return r.MimeType == "audio/mp4"
}

func (r Representation) Ext() (string, bool) {
   switch {
   case r.Audio():
      return ".m4a", true
   case r.Video():
      return ".m4v", true
   }
   return "", false
}

type ContentProtection struct {
   SchemeIdUri string `xml:"schemeIdUri,attr"`
   // this might not exist
   Default_KID string `xml:"default_KID,attr"`
   // this might not exist
   PSSH string `xml:"pssh"`
}

func (r Representation) Media() ([]string, bool) {
   if r.SegmentTemplate == nil {
      return nil, false
   }
   var refs []string
   for _, segment := range r.SegmentTemplate.SegmentTimeline.S {
      segment.T = r.SegmentTemplate.StartNumber
      for segment.R >= 0 {
         ref := func(s string) string {
            s = strings.Replace(s, "$Number$", fmt.Sprint(segment.T), 1)
            return strings.Replace(s, "$RepresentationID$", r.ID, 1)
         }(r.SegmentTemplate.Media)
         refs = append(refs, ref)
         r.SegmentTemplate.StartNumber++
         segment.R--
         segment.T++
      }
   }
   return refs, true
}

type AdaptationSet struct {
   // this might be under Representation
   ContentProtection []ContentProtection
   // pointer because we want to edit these
   Representation []*Representation
   // this might not exist
   Role *struct {
      Value string `xml:"value,attr"`
   }
   // this might not exist, or might be under Representation
   SegmentTemplate *SegmentTemplate
   // this might not exist
   Lang string `xml:"lang,attr"`
   // this might be under Representation
   MimeType string `xml:"mimeType,attr"`
}

func (r Representation) Initialization() (string, bool) {
   if r.SegmentTemplate == nil {
      return "", false
   }
   ref := func(s string) string {
      return strings.Replace(s, "$RepresentationID$", r.ID, 1)
   }(r.SegmentTemplate.Initialization)
   return ref, true
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

func (r Representation) PSSH() ([]byte, error) {
   for _, c := range r.ContentProtection {
      if c.SchemeIdUri == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         return base64.StdEncoding.DecodeString(c.PSSH)
      }
   }
   return nil, errors.New("PSSH")
}

func (r Representation) Index() (int64, error) {
   if r.SegmentBase == nil {
      return 0, errors.New("SegmentBase")
   }
   var i int64
   _, err := fmt.Sscan(r.SegmentBase.IndexRange, &i)
   if err != nil {
      return 0, err
   }
   return i, nil
}

type Representation struct {
   Bandwidth int `xml:"bandwidth,attr"`
   Codecs string `xml:"codecs,attr"`
   ID string `xml:"id,attr"`
   adaptationSet *AdaptationSet
   // this might not exist
   BaseURL string
   // this might be under AdaptationSet
   ContentProtection []ContentProtection
   // this might not exist
   Height int `xml:"height,attr"`
   // this might be under AdaptationSet
   MimeType string `xml:"mimeType,attr"`
   // this might not exist, or might be under AdaptationSet
   SegmentTemplate *SegmentTemplate
   // this might not exist
   Width int `xml:"width,attr"`
   // this might not exist
   SegmentBase *struct {
      IndexRange string `xml:"indexRange,attr"`
   }
}

////////////////////////////////////////

type SegmentTemplate struct {
   SegmentTimeline struct {
      S []struct {
         // duration
         D int `xml:"d,attr"`
         // repeat. this may not exist
         R int `xml:"r,attr"`
         // time. this may not exist
         T int `xml:"t,attr"`
      }
   }
   StartNumber int `xml:"startNumber,attr"`
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
}

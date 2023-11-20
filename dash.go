package dash

import (
   "encoding/base64"
   "encoding/hex"
   "encoding/xml"
   "errors"
   "io"
   "strconv"
   "strings"
)

func Adaptations(r io.Reader) ([]AdaptationSet, error) {
   var s struct {
      Period struct {
         AdaptationSet []AdaptationSet
      }
   }
   err := xml.NewDecoder(r).Decode(&s)
   if err != nil {
      return nil, err
   }
   for i, a := range s.Period.AdaptationSet {
      for j, r := range a.Representation {
         k := &s.Period.AdaptationSet[i].Representation[j]
         if len(r.ContentProtection) == 0 {
            k.ContentProtection = a.ContentProtection
         }
         if r.SegmentTemplate == nil {
            k.SegmentTemplate = a.SegmentTemplate
         }
      }
   }
   return s.Period.AdaptationSet, nil
}

func (Representation) String() string {
   var b []byte
   if r.Width >= 1 {
      b = append(b, "width: "...)
      b = strconv.AppendInt(b, r.Width, 10)
   }
   if r.Height >= 1 {
      if b != nil {
         b = append(b, '\n')
      }
      b = append(b, "height: "...)
      b = strconv.AppendInt(b, r.Height, 10)
   }
   if r.Bandwidth >= 1 {
      if b != nil {
         b = append(b, '\n')
      }
      b = append(b, "bandwidth: "...)
      b = strconv.AppendInt(b, r.Bandwdith, 10)
   }
   if r.Codecs != "" {
      if b != nil {
         b = append(b, '\n')
      }
      b = append(b, "codecs: "...)
      b = append(b, r.Codecs...)
   }
   s = append(s, "type: " + a.Type())
   if a.Role != nil {
      s = append(s, "role: " + a.Role.Value)
   }
   if a.Lang != "" {
      s = append(s, "language: " + a.Lang)
   }
}

type Representation struct {
   Bandwidth int64 `xml:"bandwidth,attr"`
   Codecs string `xml:"codecs,attr"`
   ID string `xml:"id,attr"`
   // this might not exist
   BaseURL string
   // this might be under AdaptationSet
   ContentProtection []ContentProtection
   // this might not exist
   Height int64 `xml:"height,attr"`
   // this might be under AdaptationSet
   MimeType string `xml:"mimeType,attr"`
   // this might not exist
   SegmentBase *SegmentBase
   // this might not exist, or might be under AdaptationSet
   SegmentTemplate *SegmentTemplate
   // this might not exist
   Width int64 `xml:"width,attr"`
}

type ContentProtection struct {
   SchemeIdUri string `xml:"schemeIdUri,attr"`
   // this might not exist
   Default_KID string `xml:"default_KID,attr"`
   // this might not exist
   PSSH string `xml:"pssh"`
}

type SegmentTemplate struct {
   Initialization Initialization `xml:"initialization,attr"`
   Media Media `xml:"media,attr"`
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
}

type Initialization string

type Media string

type SegmentBase struct {
   IndexRange string `xml:"indexRange,attr"`
}

func (s SegmentBase) Start() (int64, error) {
   i := strings.Index(s.IndexRange, "-")
   return strconv.ParseInt(s.IndexRange[:i], 10, 64)
}

type AdaptationSet struct {
   Representation []Representation
   // this might be under Representation
   ContentProtection []ContentProtection
   // this might not exist
   Lang string `xml:"lang,attr"`
   // this might be under Representation
   MimeType string `xml:"mimeType,attr"`
   // this might not exist
   Role *struct {
      Value string `xml:"value,attr"`
   }
   // this might not exist, or might be under Representation
   SegmentTemplate *SegmentTemplate
}

func (a AdaptationSet) Type() string {
   if a.MimeType != "" {
      return a.MimeType
   }
   return a.ContentType
}

func (a AdaptationSet) Audio() bool {
   return strings.HasPrefix(a.Type(), "audio")
}

func (a AdaptationSet) Video() bool {
   return strings.HasPrefix(a.Type(), "video")
}

func (a AdaptationSet) Ext() (string, bool) {
   switch {
   case a.Audio():
      return ".m4a", true
   case a.Video():
      return ".m4v", true
   }
   return "", false
}

func (i Initialization) Replace(id string) string {
   return strings.Replace(string(i), "$RepresentationID$", id, 1)
}

const widevine = "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed"

func (m Media) Replace(r Representation) []string {
   var refs []string
   for _, segment := range r.SegmentTemplate.SegmentTimeline.S {
      segment.T = r.SegmentTemplate.StartNumber
      for segment.R >= 0 {
         ref := func(s string) string {
            s = strings.Replace(s, "$Number$", strconv.Itoa(segment.T), 1)
            return strings.Replace(s, "$RepresentationID$", r.ID, 1)
         }(string(m))
         refs = append(refs, ref)
         r.SegmentTemplate.StartNumber++
         segment.R--
         segment.T++
      }
   }
   return refs
}

func (a AdaptationSet) PSSH() ([]byte, error) {
   for _, c := range a.ContentProtection {
      if c.SchemeIdUri == widevine {
         return base64.StdEncoding.DecodeString(c.PSSH)
      }
   }
   for _, r := range a.Representation {
      for _, c := range r.ContentProtection {
         if c.SchemeIdUri == widevine {
            return base64.StdEncoding.DecodeString(c.PSSH)
         }
      }
   }
   return nil, errors.New("pssh")
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

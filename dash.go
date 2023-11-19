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

type Adaptation struct {
   ContentType string `xml:"contentType,attr"`
   Lang string `xml:"lang,attr"`
   MimeType string `xml:"mimeType,attr"`
   ContentProtection []ContentProtection
   Representation []Representation
   Role *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
}

type Representation struct {
   Bandwidth int `xml:"bandwidth,attr"`
   BaseURL string
   Codecs string `xml:"codecs,attr"`
   Height int `xml:"height,attr"`
   ID string `xml:"id,attr"`
   SegmentBase *SegmentBase
   Width int `xml:"width,attr"`
   SegmentTemplate *SegmentTemplate
   ContentProtection []ContentProtection
}

type ContentProtection struct {
   Default_KID string `xml:"default_KID,attr"`
   PSSH string `xml:"pssh"`
   Scheme_ID_URI string `xml:"schemeIdUri,attr"`
}

type SegmentBase struct {
   Index_Range string `xml:"indexRange,attr"`
}

type SegmentTemplate struct {
   Start_Number int `xml:"startNumber,attr"`
   Segment_Timeline struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R int `xml:"r,attr"` // repeat
         T int `xml:"t,attr"` // time
      }
   } `xml:"SegmentTimeline"`
   Initialization Initialization `xml:"initialization,attr"`
   Media Media `xml:"media,attr"`
}

type Initialization string

type Media string

func (s SegmentBase) Start() (int64, error) {
   i := strings.Index(s.Index_Range, "-")
   return strconv.ParseInt(s.Index_Range[:i], 10, 64)
}

//////////////////////////////////////////////

func (a Adaptation) String() string {
   var s []string
   for i, r := range a.Representation {
      if i >= 1 {
         s = append(s, "")
      }
      if r.Width >= 1 {
         s = append(s, "width: " + strconv.Itoa(r.Width))
      }
      if r.Height >= 1 {
         s = append(s, "height: " + strconv.Itoa(r.Height))
      }
      if r.Bandwidth >= 1 {
         s = append(s, "bandwidth: " + strconv.Itoa(r.Bandwidth))
      }
      if r.Codecs != "" {
         s = append(s, "codecs: " + r.Codecs)
      }
      s = append(s, "type: " + a.Type())
      if a.Role != nil {
         s = append(s, "role: " + a.Role.Value)
      }
      if a.Lang != "" {
         s = append(s, "language: " + a.Lang)
      }
   }
   return strings.Join(s, "\n")
}

func (a Adaptation) Type() string {
   if a.MimeType != "" {
      return a.MimeType
   }
   return a.ContentType
}

func (a Adaptation) Audio() bool {
   return strings.HasPrefix(a.Type(), "audio")
}

func (a Adaptation) Video() bool {
   return strings.HasPrefix(a.Type(), "video")
}

func (a Adaptation) Ext() (string, bool) {
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
   for _, segment := range r.SegmentTemplate.Segment_Timeline.S {
      segment.T = r.SegmentTemplate.Start_Number
      for segment.R >= 0 {
         ref := func(s string) string {
            s = strings.Replace(s, "$Number$", strconv.Itoa(segment.T), 1)
            return strings.Replace(s, "$RepresentationID$", r.ID, 1)
         }(string(m))
         refs = append(refs, ref)
         r.SegmentTemplate.Start_Number++
         segment.R--
         segment.T++
      }
   }
   return refs
}

func Adaptations(r io.Reader) ([]Adaptation, error) {
   var s struct {
      Period struct {
         Adaptation []Adaptation `xml:"AdaptationSet"`
      }
   }
   err := xml.NewDecoder(r).Decode(&s)
   if err != nil {
      return nil, err
   }
   for i, a := range s.Period.Adaptation {
      for j, r := range a.Representation {
         k := &s.Period.Adaptation[i].Representation[j]
         if len(r.ContentProtection) == 0 {
            k.ContentProtection = a.ContentProtection
         }
         if r.SegmentTemplate == nil {
            k.SegmentTemplate = a.SegmentTemplate
         }
      }
   }
   return s.Period.Adaptation, nil
}

func (a Adaptation) PSSH() ([]byte, error) {
   for _, c := range a.ContentProtection {
      if c.Scheme_ID_URI == widevine {
         return base64.StdEncoding.DecodeString(c.PSSH)
      }
   }
   for _, r := range a.Representation {
      for _, c := range r.ContentProtection {
         if c.Scheme_ID_URI == widevine {
            return base64.StdEncoding.DecodeString(c.PSSH)
         }
      }
   }
   return nil, errors.New("pssh")
}

func (r Representation) Default_KID() ([]byte, error) {
   for _, c := range r.ContentProtection {
      if c.Scheme_ID_URI == "urn:mpeg:dash:mp4protection:2011" {
         c.Default_KID = strings.ReplaceAll(c.Default_KID, "-", "")
         return hex.DecodeString(c.Default_KID)
      }
   }
   return nil, errors.New("default_KID")
}


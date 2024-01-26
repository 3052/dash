package dash

import (
   "encoding/base64"
   "encoding/hex"
   "errors"
   "fmt"
   "strings"
)

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
      IndexRange string `xml:"indexRange,attr"`
   } `xml:"SegmentBase"`
   // this might not exist, or might be under AdaptationSet
   SegmentTemplate *SegmentTemplate
   // this might not exist
   Width *int `xml:"width,attr"`
}

// wikipedia.org/wiki/Mutator_method
func (r Representation) GetSegmentTemplate(a Adaptation) *SegmentTemplate {
   if r.SegmentTemplate != nil {
      return r.SegmentTemplate
   }
   return a.SegmentTemplate
}

// wikipedia.org/wiki/Mutator_method
func (r Representation) GetMimeType(a Adaptation) string {
   if r.MimeType != "" {
      return r.MimeType
   }
   return a.MimeType
}

func (r Representation) Ext(a Adaptation) (string, bool) {
   switch r.GetMimeType(a) {
   case "audio/mp4":
      return ".m4a", true
   case "video/mp4":
      return ".m4v", true
   }
   return "", false
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

func (r Representation) Initialization(a Adaptation) (string, bool) {
   if v := r.GetSegmentTemplate(a); v != nil {
      if v := v.Initialization; v != "" {
         return strings.Replace(v, "$RepresentationID$", r.ID, 1), true
      }
   }
   return "", false
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

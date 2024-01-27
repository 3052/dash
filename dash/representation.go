package dash

import (
   "encoding/base64"
   "encoding/hex"
   "errors"
   "fmt"
   "strings"
)

type IndexRange struct {
   Start int
   End int
}

func (i *IndexRange) UnmarshalText(b []byte) error {
   _, err := fmt.Sscanf(string(b), "%v-%v", &i.Start, &i.End)
   if err != nil {
      return err
   }
   return nil
}

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
      IndexRange IndexRange `xml:"indexRange,attr"`
   }
   // this might not exist, or might be under AdaptationSet
   SegmentTemplate *SegmentTemplate
   // this might not exist
   Width *int `xml:"width,attr"`
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

func (r Representation) Ext() (string, bool) {
   switch r.MimeType {
   case "audio/mp4":
      return ".m4a", true
   case "video/mp4":
      return ".m4v", true
   }
   return "", false
}

func (r Representation) Initialization() (string, bool) {
   if v := r.SegmentTemplate; v != nil {
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

func (r Representation) PSSH() ([]byte, error) {
   for _, c := range r.ContentProtection {
      if c.SchemeIdUri == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         return base64.StdEncoding.DecodeString(c.PSSH)
      }
   }
   return nil, errors.New("PSSH")
}

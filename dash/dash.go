package dash

import (
   "encoding/base64"
   "encoding/hex"
   "errors"
   "fmt"
   "strings"
)

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

// return a slice so we can measure progress
func (r Representation) Media() []string {
   st := r.SegmentTemplate
   if st == nil {
      return nil
   }
   var media []string
   for _, segment := range st.SegmentTimeline.S {
      for segment.R >= 0 {
         number := fmt.Sprint(st.StartNumber)
         medium := strings.Replace(st.Media, "$Number$", number, 1)
         medium = strings.Replace(medium, "$RepresentationID$", r.ID, 1)
         media = append(media, medium)
         segment.R--
         st.StartNumber++
      }
   }
   return media
}
func (r Range) Scan(start, end *uint32) error {
   _, err := fmt.Sscanf(string(r), "%v-%v", start, end)
   if err != nil {
      return err
   }
   return nil
}

type Range string

const Template = `<style>
table {
   border-collapse: collapse;
   margin: 9px;
}
td {
   border-style: solid;
   border-width: thin;
}
td,
th {
   padding-bottom: 9px;
   padding-left: 9px;
   padding-right: 9px;
   padding-top: 9px;
}
</style>
<table>
<tr>
   <th>width</th>
   <th>height</th>
   <th>bandwidth</th>
   <th>codecs</th>
   <th>type</th>
   <th>role</th>
   <th>language</th>
   <th>ID</th>
   <th>period</th>
</tr>
{{ range $period := .Period -}}
   {{ range $adaptation := .AdaptationSet -}}
      {{ range .Representation -}}
<tr>
   <td>{{ .Width }}</td>
   <td>{{ .Height }}</td>
   <td>{{ .Bandwidth }}</td>
         {{ with .Codecs -}}
   <td>{{ . }}</td>
         {{ else -}}
   <td>{{ $adaptation.Codecs }}</td>
         {{ end -}}
         {{ with .MimeType -}}
   <td>{{ . }}</td>
         {{ else -}}
   <td>{{ $adaptation.MimeType }}</td>
         {{ end -}}
         {{ with $adaptation.Role -}}
   <td>{{ .Value }}</td>
         {{ else -}}
   <td></td>
         {{ end -}}
   <td>{{ $adaptation.Lang }}</td>
   <td>{{ .ID }}</td>
   <td>{{ $period.ID }}</td>
</tr>
      {{ end -}}
   {{ end -}}
{{ end -}}
</table>
`

type ContentProtection struct {
   SchemeIdUri string `xml:"schemeIdUri,attr"`
   // this might not exist
   Default_KID string `xml:"default_KID,attr"`
   // this might not exist
   PSSH string `xml:"pssh"`
}

type SegmentTemplate struct {
   Media string `xml:"media,attr"`
   SegmentTimeline struct {
      S []struct {
         // duration
         D int `xml:"d,attr"`
         // repeat. this may not exist
         R int `xml:"r,attr"`
      }
   }
   StartNumber int `xml:"startNumber,attr"`
   // this may not exist
   Initialization string `xml:"initialization,attr"`
}

package dash

import "fmt"

type Range string

func (r Range) Scan(start, end any) error {
   _, err := fmt.Sscanf(string(r), "%v-%v", start, end)
   if err != nil {
      return err
   }
   return nil
}

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

type AdaptationSet struct {
   // this might be under Representation
   Codecs string `xml:"codecs,attr"`
   // this might be under Representation
   ContentProtection []ContentProtection
   // this might not exist
   Lang string `xml:"lang,attr"`
   // this might be under Representation
   MimeType string `xml:"mimeType,attr"`
   Representation []*Representation
   // this might not exist
   Role *struct {
      Value string `xml:"value,attr"`
   }
   // this might not exist, or might be under Representation
   SegmentTemplate *SegmentTemplate
}

type ContentProtection struct {
   SchemeIdUri string `xml:"schemeIdUri,attr"`
   // this might not exist
   PSSH string `xml:"pssh"`
}

// media presentation description
// wikipedia.org/wiki/Dynamic_Adaptive_Streaming_over_HTTP
type MPD struct {
   Period []*Period
}

type Period struct {
   AdaptationSet []*AdaptationSet
   ID string `xml:"id,attr"`
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
      Initialization struct {
         Range Range `xml:"range,attr"`
      }
      IndexRange Range `xml:"indexRange,attr"`
   }
   // this might not exist, or might be under AdaptationSet
   SegmentTemplate *SegmentTemplate
   // this might not exist
   Width *int `xml:"width,attr"`
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

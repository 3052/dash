package dash

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
   // pointer because we want to edit these
   Representation []Representation
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
   Default_KID string `xml:"default_KID,attr"`
   // this might not exist
   PSSH string `xml:"pssh"`
}

type MPD struct {
   Period []struct {
      AdaptationSet []AdaptationSet
      ID string `xml:"id,attr"`
   }
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

func (m MPD) Some(f func(Representation) bool) {
   for _, period := range m.Period {
      for _, adapt := range period.AdaptationSet {
         for _, represent := range adapt.Representation {
            if represent.Codecs == "" {
               represent.Codecs = adapt.Codecs
            }
            if represent.ContentProtection == nil {
               represent.ContentProtection = adapt.ContentProtection
            }
            if represent.MimeType == "" {
               represent.MimeType = adapt.MimeType
            }
            if represent.SegmentTemplate == nil {
               represent.SegmentTemplate = adapt.SegmentTemplate
            }
            if !f(represent) {
               return
            }
         }
      }
   }
}

func (m MPD) Every(f func(Representation)) {
   m.Some(func(r Representation) bool {
      f(r)
      return true
   })
}

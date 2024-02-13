package dash

const ModeLine = `{{ $first := true -}}
{{ range $period := .Period -}}
{{ range $adaptation := .AdaptationSet -}}
{{ range .Representation -}}
{{ if $first -}}
   {{ $first = false -}}
{{ else }}
{{ end -}}
{{ with .Width -}}
width = {{ . }}
{{ end -}}
{{ with .Height -}}
height = {{ . }}
{{ end -}}
bandwidth = {{ .Bandwidth }}
{{ with .Codecs -}}
codecs = {{ . }}
{{ end -}}
{{ with $adaptation.Codecs -}}
codecs = {{ . }}
{{ end -}}
{{ with .MimeType -}}
type = {{ . }}
{{ else -}}
type = {{ $adaptation.MimeType }}
{{ end -}}
{{ with $adaptation.Role -}}
role = {{ .Value }}
{{ end -}}
{{ with $adaptation.Lang -}}
lang = {{ . }}
{{ end -}}
id = {{ .ID }}
{{ with $period.ID -}}
period = {{ . }}
{{ end -}}
{{ end -}}
{{ end -}}
{{ end }}`

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
   Height int `xml:"height,attr"`
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
   Width int `xml:"width,attr"`
}

type AdaptationSet struct {
   // this might be under Representation
   Codecs string `xml:"codecs,attr"`
   // this might be under Representation
   ContentProtection []ContentProtection
   // this might not exist
   Lang string `xml:"lang,attr"`
   // this might be under Representation
   MimeType string `xml:"mimeType,attr"`
   Representation []Representation
   // this might not exist
   Role *struct {
      Value string `xml:"value,attr"`
   }
   // this might not exist, or might be under Representation
   SegmentTemplate *SegmentTemplate
}

type Period struct {
   AdaptationSet []AdaptationSet
   ID string `xml:"id,attr"`
}

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

// media presentation description
// wikipedia.org/wiki/Dynamic_Adaptive_Streaming_over_HTTP
type MPD struct {
   Period []Period
}

// godocs.io/flag#Visit
func (m MPD) Visit(f func(Pointer)) {
   m.Contains(func(p Pointer) bool {
      f(p)
      return false
   })
}

func (m MPD) Contains(f func(Pointer) bool) bool {
   for _, period := range m.Period {
      for _, adapt := range period.AdaptationSet {
         for _, represent := range adapt.Representation {
            var p Pointer
            p.AdaptationSet = &adapt
            p.Period = &period
            p.Representation = &represent
            if f(p) {
               return true
            }
         }
      }
   }
   return false
}

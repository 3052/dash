package dash

const ModeLine = 
"{{ $first := true }}" +
"{{ range $period := .Period }}" +
   "{{ range $adaptation := .AdaptationSet }}" +
      "{{ range .Representation }}" +
         "{{ if $first }}" +
            "{{ $first = false }}" +
         "{{ else }}" +
"\n" +
         "{{ end }}" +
         "{{ with .Width }}" +
"width = {{ . }}\n" +
         "{{ end }}" +
         "{{ with .Height }}" +
"height = {{ . }}\n" +
         "{{ end }}" +
"bandwidth = {{ .Bandwidth }}\n" +
         "{{ with .Codecs }}" +
"codecs = {{ . }}\n" +
         "{{ end }}" +
         "{{ with $adaptation.Codecs }}" +
"codecs = {{ . }}\n" +
         "{{ end }}" +
         "{{ with .MimeType }}" +
"type = {{ . }}\n" +
         "{{ else }}" +
"type = {{ $adaptation.MimeType }}\n" +
         "{{ end }}" +
         "{{ with $adaptation.Role }}" +
"role = {{ .Value }}\n" +
         "{{ end }}" +
         "{{ with $adaptation.Lang }}" +
"lang = {{ . }}\n" +
         "{{ end }}" +
"id = {{ .ID }}\n" +
         "{{ with $period.ID }}" +
"period = {{ . }}\n" +
         "{{ end }}" +
      "{{ end }}" +
   "{{ end }}" +
"{{ end }}"

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

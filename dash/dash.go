package dash

import (
   "fmt"
   "strconv"
   "strings"
)

func (m MPD) Every(f func(Pointer)) {
   m.Some(func(p Pointer) bool {
      f(p)
      return true
   })
}

func (m MPD) Some(f func(Pointer) bool) {
   for _, period := range m.Period {
      for _, adapt := range period.AdaptationSet {
         for _, represent := range adapt.Representation {
            var p Pointer
            p.AdaptationSet = adapt
            p.Period = period
            p.Representation = represent
            if !f(p) {
               return
            }
         }
      }
   }
}

type Pointer struct {
   AdaptationSet *AdaptationSet
   Period *Period
   Representation *Representation
}

func (p Pointer) Codecs() string {
   if a := p.AdaptationSet; a.Codecs != "" {
      return a.Codecs
   }
   return p.Representation.Codecs
}

func (p Pointer) Ext() (string, bool) {
   switch p.MimeType() {
   case "audio/mp4":
      return ".m4a", true
   case "video/mp4":
      return ".m4v", true
   }
   return "", false
}

func (p Pointer) Initialization() (string, bool) {
   if st := p.segmentTemplate(); st != nil {
      if i := st.Initialization; i != "" {
         i = strings.Replace(i, "$RepresentationID$", p.Representation.ID, 1)
         return i, true
      }
   }
   return "", false
}

// return a slice so we can measure progress
func (p Pointer) Media() []string {
   replace := func(s string, i int) string {
      s = strings.Replace(s, "$RepresentationID$", p.Representation.ID, 1)
      return strings.Replace(s, "$Number$", strconv.Itoa(i), 1)
   }
   var media []string
   if st := p.segmentTemplate(); st != nil {
      for _, segment := range st.SegmentTimeline.S {
         for segment.R >= 0 {
            medium := replace(st.Media, st.StartNumber)
            media = append(media, medium)
            segment.R--
            st.StartNumber++
         }
      }
   }
   return media
}

func (p Pointer) MimeType() string {
   if a := p.AdaptationSet; a.MimeType != "" {
      return a.MimeType
   }
   return p.Representation.MimeType
}

func (p Pointer) segmentTemplate() *SegmentTemplate {
   if a := p.AdaptationSet; a.SegmentTemplate != nil {
      return a.SegmentTemplate
   }
   return p.Representation.SegmentTemplate
}

type Range string

func (r Range) Scan() (int, int, error) {
   var start, end int
   _, err := fmt.Sscanf(string(r), "%v-%v", &start, &end)
   if err != nil {
      return 0, 0, err
   }
   return start, end, nil
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

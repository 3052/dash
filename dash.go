package dash

import (
   "encoding/base64"
   "encoding/xml"
   "math"
   "net/url"
   "strconv"
   "strings"
   "text/template"
   "time"
)

type Template struct {
   Template *template.Template
}

func (r Representation) Initialization(t SegmentTemplate) (string, error) {
   var medium strings.Builder
   var hello struct {
      Representation struct {
         Id string
      }
   }
   hello.Representation.Id = r.Id
   err := t.Initialization.Template.Execute(&medium, hello)
   if err != nil {
      return "", err
   }
   return medium.String(), nil
}

func (r Representation) Media(t SegmentTemplate) ([]string, error) {
   var media []string
   var hello struct {
      Number uint
      Representation struct {
         Id string
      }
      Time uint
   }
   hello.Number = t.StartNumber
   hello.Time = t.PresentationTimeOffset
   hello.Representation.Id = r.Id
   if t.SegmentTimeline != nil {
      for _, segment := range t.SegmentTimeline.S {
         for range 1 + segment.R {
            var medium strings.Builder
            err := t.Media.Template.Execute(&medium, hello)
            if err != nil {
               return nil, err
            }
            media = append(media, medium.String())
            hello.Number++
            hello.Time += segment.D
         }
      }
   } else {
      seconds := r.adaptation_set.period.get_duration().Duration.Seconds()
      for range t.segment_count(seconds) {
         var medium strings.Builder
         err := t.Media.Template.Execute(&medium, hello)
         if err != nil {
            return nil, err
         }
         media = append(media, medium.String())
         hello.Number++
      }
   }
   return media, nil
}

type SegmentTemplate struct {
   Media Template `xml:"media,attr"`
   Initialization *Template `xml:"initialization,attr"`
   StartNumber uint `xml:"startNumber,attr"`
   Duration               uint64 `xml:"duration,attr"`
   PresentationTimeOffset uint   `xml:"presentationTimeOffset,attr"`
   Timescale              uint64 `xml:"timescale,attr"`
   SegmentTimeline        *struct {
      S []struct {
         D uint `xml:"d,attr"` // duration
         R uint `xml:"r,attr"` // repeat
      }
   }
}

func (r Representation) String() string {
   var b []byte
   if v := r.get_width(); v >= 1 {
      b = append(b, "width = "...)
      b = strconv.AppendUint(b, v, 10)
   }
   if v := r.get_height(); v >= 1 {
      if b != nil {
         b = append(b, '\n')
      }
      b = append(b, "height = "...)
      b = strconv.AppendUint(b, v, 10)
   }
   if b != nil {
      b = append(b, '\n')
   }
   b = append(b, "bandwidth = "...)
   b = strconv.AppendUint(b, r.Bandwidth, 10)
   if v := r.get_codecs(); v != "" {
      b = append(b, "\ncodecs = "...)
      b = append(b, v...)
   }
   b = append(b, "\nmimeType = "...)
   b = append(b, r.get_mime_type()...)
   if v := r.adaptation_set.Role; v != nil {
      b = append(b, "\nrole = "...)
      b = append(b, v.Value...)
   }
   if v := r.adaptation_set.Lang; v != "" {
      b = append(b, "\nlang = "...)
      b = append(b, v...)
   }
   if v := r.adaptation_set.period.Id; v != "" {
      b = append(b, "\nperiod = "...)
      b = append(b, v...)
   }
   b = append(b, "\nid = "...)
   b = append(b, r.Id...)
   return string(b)
}

func (r Representation) get_segment_template() (*SegmentTemplate, bool) {
   if r.SegmentTemplate != nil {
      return r.SegmentTemplate, true
   }
   if r.adaptation_set.SegmentTemplate != nil {
      return r.adaptation_set.SegmentTemplate, true
   }
   return nil, false
}

func (r Representation) get_mime_type() string {
   if r.MimeType != "" {
      return r.MimeType
   }
   return r.adaptation_set.MimeType
}

func (r Representation) get_width() uint64 {
   if r.Width >= 1 {
      return r.Width
   }
   return r.adaptation_set.Width
}

func (r Representation) get_height() uint64 {
   if r.Height >= 1 {
      return r.Height
   }
   return r.adaptation_set.Height
}

////////

func (r Representation) get_content_protection() []ContentProtection {
   if len(r.ContentProtection) >= 1 {
      return r.ContentProtection
   }
   return r.adaptation_set.ContentProtection
}

func (r Representation) get_codecs() string {
   if r.Codecs != "" {
      return r.Codecs
   }
   return r.adaptation_set.Codecs
}

func (r Representation) GetBaseUrl() (*BaseUrl, bool) {
   var u *url.URL
   if v := r.adaptation_set.period.mpd.BaseUrl; v != nil {
      u = new(url.URL)
      *u = *v.Url
   }
   if v := r.adaptation_set.period.BaseUrl; v != nil {
      if u == nil {
         u = new(url.URL)
      }
      u = u.ResolveReference(v.Url)
   }
   if v := r.BaseUrl; v != nil {
      if u == nil {
         u = new(url.URL)
      }
      u = u.ResolveReference(v.Url)
   }
   if u != nil {
      return &BaseUrl{u}, true
   }
   return nil, false
}

type AdaptationSet struct {
   Codecs            string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height            uint64 `xml:"height,attr"`
   Lang              string `xml:"lang,attr"`
   MimeType          string `xml:"mimeType,attr"`
   Representation    []Representation
   Role              *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
   Width           uint64 `xml:"width,attr"`
   period          *Period
}

type BaseUrl struct {
   Url *url.URL
}

func (b *BaseUrl) UnmarshalText(text []byte) error {
   b.Url = new(url.URL)
   return b.Url.UnmarshalBinary(text)
}

type ContentProtection struct {
   Pssh        Pssh   `xml:"pssh"`
   SchemeIdUri string `xml:"schemeIdUri,attr"`
}

func (d *Duration) UnmarshalText(text []byte) error {
   var err error
   d.Duration, err = time.ParseDuration(strings.ToLower(
      strings.TrimPrefix(string(text), "PT"),
   ))
   if err != nil {
      return err
   }
   return nil
}

type Duration struct {
   Duration time.Duration
}

type Mpd struct {
   BaseUrl *BaseUrl `xml:"BaseURL"`
   MediaPresentationDuration *Duration `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

type Period struct {
   AdaptationSet []AdaptationSet
   BaseUrl *BaseUrl `xml:"BaseURL"`
   Duration      *Duration `xml:"duration,attr"`
   Id            string    `xml:"id,attr"`
   mpd           *Mpd
}

func (p Period) get_duration() *Duration {
   if p.Duration != nil {
      return p.Duration
   }
   return p.mpd.MediaPresentationDuration
}

type Pssh []byte

func (p *Pssh) UnmarshalText(src []byte) error {
   var err error
   *p, err = base64.StdEncoding.AppendDecode(nil, src)
   if err != nil {
      return err
   }
   return nil
}

// SegmentIndexBox uses:
// unsigned int(32) subsegment_duration;
// but range values can exceed 32 bits
type Range struct {
   Start uint64
   End   uint64
}

func (r Range) MarshalText() ([]byte, error) {
   b := strconv.AppendUint(nil, r.Start, 10)
   b = append(b, '-')
   return strconv.AppendUint(b, r.End, 10), nil
}

func (r *Range) UnmarshalText(text []byte) error {
   // the current testdata always has `-`, so lets assume for now
   start, end, _ := strings.Cut(string(text), "-")
   var err error
   r.Start, err = strconv.ParseUint(start, 10, 64)
   if err != nil {
      return err
   }
   r.End, err = strconv.ParseUint(end, 10, 64)
   if err != nil {
      return err
   }
   return nil
}

type SegmentBase struct {
   Initialization struct {
      Range Range `xml:"range,attr"`
   }
   IndexRange Range `xml:"indexRange,attr"`
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#addressing-simple-to-explicit
func (s SegmentTemplate) segment_count(seconds float64) uint64 {
   seconds /= float64(s.Duration) / float64(s.get_timescale())
   return uint64(math.Ceil(seconds))
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#timing-sampletimeline
func (s SegmentTemplate) get_timescale() uint64 {
   if s.Timescale >= 1 {
      return s.Timescale
   }
   return 1
}

type Representation struct {
   Bandwidth         uint64 `xml:"bandwidth,attr"`
   BaseUrl           *BaseUrl   `xml:"BaseURL"`
   Codecs            string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height            uint64 `xml:"height,attr"`
   Id                string `xml:"id,attr"`
   MimeType          string `xml:"mimeType,attr"`
   SegmentBase       *SegmentBase
   SegmentTemplate   *SegmentTemplate
   Width             uint64 `xml:"width,attr"`
   adaptation_set    *AdaptationSet
}

func (t *Template) UnmarshalText(text []byte) error {
   var (
      err error
      str = string(text)
   )
   t.Template, err = t.Template.Parse(strings.NewReplacer(
      "$Number$", "{{.Number}}",
      "$Number%02d$", `{{printf "%02d" .Number}}`,
      "$Number%03d$", `{{printf "%03d" .Number}}`,
      "$Number%04d$", `{{printf "%04d" .Number}}`,
      "$Number%05d$", `{{printf "%05d" .Number}}`,
      "$Number%06d$", `{{printf "%06d" .Number}}`,
      "$Number%07d$", `{{printf "%07d" .Number}}`,
      "$Number%08d$", `{{printf "%08d" .Number}}`,
      "$Number%09d$", `{{printf "%09d" .Number}}`,
      "$RepresentationID$", "{{.Representation.Id}}",
      "$Time$", "{{.Time}}",
   ).Replace(str))
   if err != nil {
      return err
   }
   return nil
}

func Unmarshal(text []byte, base *url.URL) ([]Representation, error) {
   var media Mpd
   err := xml.Unmarshal(text, &media)
   if err != nil {
      return nil, err
   }
   if media.BaseUrl == nil {
      if base != nil {
         media.BaseUrl = &BaseUrl{base}
      }
   }
   var reps []Representation
   for _, per := range media.Period {
      per.mpd = &media
      for _, ada := range per.AdaptationSet {
         ada.period = &per
         for _, rep := range ada.Representation {
            rep.adaptation_set = &ada
            reps = append(reps, rep)
         }
      }
   }
   return reps, nil
}

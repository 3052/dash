package dash

import (
   "math"
   "net/url"
   "strconv"
   "strings"
   "text/template"
   "time"
)

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
   Pssh        string   `xml:"pssh"`
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

func (p Period) get_duration() *Duration {
   if p.Duration != nil {
      return p.Duration
   }
   return p.mpd.MediaPresentationDuration
}

type Period struct {
   AdaptationSet []AdaptationSet
   BaseUrl *BaseUrl `xml:"BaseURL"`
   Duration      *Duration `xml:"duration,attr"`
   Id            string    `xml:"id,attr"`
   mpd           *Mpd
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

func (Template) Error() string {
   return "Template"
}

type Template struct {
   Template *template.Template
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

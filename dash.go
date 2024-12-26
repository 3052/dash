package dash

import (
   "encoding/base64"
   "encoding/xml"
   "fmt"
   "math"
   "net/url"
   "strconv"
   "strings"
   "time"
)

type ContentProtection struct {
   Pssh        Pssh   `xml:"pssh"`
   SchemeIdUri string `xml:"schemeIdUri,attr"`
}

type Pssh []byte

func (p *Pssh) UnmarshalText(data []byte) error {
   var err error
   *p, err = base64.StdEncoding.AppendDecode(nil, data)
   if err != nil {
      return err
   }
   return nil
}

type AdaptationSet struct {
   Codecs            string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height            uint64 `xml:"height,attr"`
   Lang              string `xml:"lang,attr"`
   MaxHeight         int    `xml:"maxHeight,attr"`
   MaxWidth          int    `xml:"maxWidth,attr"`
   MimeType          string `xml:"mimeType,attr"`
   Representation    []Representation
   Role              *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
   Width           uint64 `xml:"width,attr"`
}

///

func (r *Range) MarshalText() ([]byte, error) {
   b := strconv.AppendUint(nil, r.Start, 10)
   b = append(b, '-')
   return strconv.AppendUint(b, r.End, 10), nil
}

// SegmentIndexBox uses:
// unsigned int(32) subsegment_duration;
// but range values can exceed 32 bits
type Range struct {
   Start uint64
   End   uint64
}

func (r *Range) UnmarshalText(data []byte) error {
   // the current testdata always has `-`, so lets assume for now
   start, end, _ := strings.Cut(string(data), "-")
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

func (r *Representation) Initialization() (string, bool) {
   if v, ok := r.get_segment_template(); ok {
      if v := v.Initialization; v != "" {
         return r.id(v), true
      }
   }
   return "", false
}

func (r *Representation) Widevine() (Pssh, bool) {
   for _, v := range r.get_content_protection() {
      if v.SchemeIdUri == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         if len(v.Pssh) >= 1 {
            return v.Pssh, true
         }
      }
   }
   return nil, false
}

func (r *Representation) id(value string) string {
   return strings.Replace(value, "$RepresentationID$", r.Id, 1)
}

type SegmentBase struct {
   Initialization struct {
      Range Range `xml:"range,attr"`
   }
   IndexRange Range `xml:"indexRange,attr"`
}

func (s *SegmentTemplate) time(value uint) string {
   format := strings.Replace(s.Media, "$Time$", "%d", 1)
   return fmt.Sprintf(format, value)
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#timing-sampletimeline
func (s *SegmentTemplate) get_timescale() uint64 {
   if s.Timescale >= 1 {
      return s.Timescale
   }
   return 1
}

func (s *SegmentTemplate) number(value uint) string {
   format := strings.NewReplacer(
      "%", "%%",
      "$Number$", "%d",
      "$Number%02d$", "%02d",
      "$Number%03d$", "%03d",
      "$Number%04d$", "%04d",
      "$Number%05d$", "%05d",
      "$Number%06d$", "%06d",
      "$Number%07d$", "%07d",
      "$Number%08d$", "%08d",
      "$Number%09d$", "%09d",
   ).Replace(s.Media)
   return fmt.Sprintf(format, value)
}

// If using `$Number$` addressing, the number of the first segment reference is
// defined by `SegmentTemplate@startNumber` (default value 1)
//
// If using `$Time$` addressing, the value for each [=segment reference=] is the
// [=segment start point=] on the [=sample timeline=], in [=timescale units=]
//
// github.com/Dash-Industry-Forum/Guidelines-TimingModel/blob/master/22-Addressing.inc.md
func (s *SegmentTemplate) start() uint {
   if strings.Contains(s.Media, "$Time$") {
      return s.PresentationTimeOffset
   }
   if s.StartNumber != nil {
      return *s.StartNumber
   }
   return 1
}

///

type BaseUrl struct {
   Url *url.URL
}

func (b *BaseUrl) UnmarshalText(data []byte) error {
   b.Url = &url.URL{}
   return b.Url.UnmarshalBinary(data)
}

func (r *Representation) GetBaseUrl() (*BaseUrl, bool) {
   var u *url.URL
   if v := r.adaptation_set.period.mpd.BaseUrl; v != nil {
      u = &url.URL{}
      *u = *v.Url
   }
   if v := r.adaptation_set.period.BaseUrl; v != nil {
      if u == nil {
         u = &url.URL{}
      }
      u = u.ResolveReference(v.Url)
   }
   if v := r.BaseUrl; v != nil {
      if u == nil {
         u = &url.URL{}
      }
      u = u.ResolveReference(v.Url)
   }
   if u != nil {
      return &BaseUrl{u}, true
   }
   return nil, false
}

func Unmarshal(data []byte, base *url.URL) ([]Representation, error) {
   var media Mpd
   err := xml.Unmarshal(data, &media)
   if err != nil {
      return nil, err
   }
   if base != nil {
      if media.BaseUrl != nil {
         base = base.ResolveReference(media.BaseUrl.Url)
      }
      media.BaseUrl = &BaseUrl{base}
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

type Representation struct {
   Bandwidth         uint64   `xml:"bandwidth,attr"`
   BaseUrl           *BaseUrl `xml:"BaseURL"`
   Codecs            string   `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height            uint64 `xml:"height,attr"`
   Id                string `xml:"id,attr"`
   MimeType          string `xml:"mimeType,attr"`
   SegmentBase       *SegmentBase
   SegmentTemplate   *SegmentTemplate
   Width             uint64 `xml:"width,attr"`
   adaptation_set    *AdaptationSet
}

type SegmentTemplate struct {
   Duration               float64 `xml:"duration,attr"`
   Initialization         string  `xml:"initialization,attr"`
   Media                  string  `xml:"media,attr"`
   PresentationTimeOffset uint    `xml:"presentationTimeOffset,attr"`
   SegmentTimeline        *struct {
      S []struct {
         D uint `xml:"d,attr"` // duration
         R uint `xml:"r,attr"` // repeat
      }
   }
   Timescale   uint64 `xml:"timescale,attr"`
   StartNumber *uint  `xml:"startNumber,attr"`
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#addressing-simple-to-explicit
func (s *SegmentTemplate) segment_count(seconds float64) uint64 {
   seconds /= s.Duration / float64(s.get_timescale())
   return uint64(math.Ceil(seconds))
}

func (d *Duration) UnmarshalText(data []byte) error {
   var err error
   d.Duration, err = time.ParseDuration(strings.ToLower(
      strings.TrimPrefix(string(data), "PT"),
   ))
   if err != nil {
      return err
   }
   return nil
}

type Duration struct {
   Duration time.Duration
}

func (r *Representation) Media() []string {
   template, ok := r.get_segment_template()
   if !ok {
      return nil
   }
   number := template.start()
   template.Media = r.id(template.Media)
   var media []string
   if template.SegmentTimeline != nil {
      for _, segment := range template.SegmentTimeline.S {
         for range 1 + segment.R {
            var medium string
            if strings.Contains(template.Media, "$Time$") {
               medium = template.time(number)
               number += segment.D
            } else {
               medium = template.number(number)
               number++
            }
            media = append(media, medium)
         }
      }
   } else {
      seconds := r.adaptation_set.period.get_duration().Seconds()
      for range template.segment_count(seconds) {
         media = append(media, template.number(number))
         number++
      }
   }
   return media
}

func (p *Period) get_duration() time.Duration {
   if p.Duration != nil {
      return p.Duration.Duration
   }
   return p.mpd.MediaPresentationDuration.Duration
}

type Period struct {
   AdaptationSet []AdaptationSet
   BaseUrl       *BaseUrl  `xml:"BaseURL"`
   Duration      *Duration `xml:"duration,attr"`
   Id            string    `xml:"id,attr"`
   mpd           *Mpd
}

type Mpd struct {
   BaseUrl                   *BaseUrl  `xml:"BaseURL"`
   MediaPresentationDuration *Duration `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

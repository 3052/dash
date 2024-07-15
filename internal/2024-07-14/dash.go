package dash

import (
   "encoding/xml"
   "fmt"
   "math"
   "net/url"
   "strconv"
   "strings"
   "time"
)

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

func (r Representation) Media() []string {
   // `template` is a pointer, so if we edit `template.Media` it is permanent
   template, ok := r.get_segment_template()
   if !ok {
      return nil
   }
   id := r.id(template.Media)
   number := template.start()
   var media []string
   if template.SegmentTimeline != nil {
      for _, segment := range template.SegmentTimeline.S {
         for range 1 + segment.R {
            var medium string
            if strings.Contains(template.Media, "$Time$") {
               medium = replace_time(id, number)
               number += segment.D
            } else {
               medium = replace_number(id, number)
               number++
            }
            media = append(media, medium)
         }
      }
   } else {
      seconds := r.adaptation_set.period.get_duration().Duration.Seconds()
      for range template.segment_count(seconds) {
         media = append(media, replace_number(id, number))
         number++
      }
   }
   return media
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

func (r Representation) id(value string) string {
   return strings.Replace(value, "$RepresentationID$", r.Id, 1)
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

type ContentProtection struct {
   Pssh        string   `xml:"pssh"`
   SchemeIdUri string `xml:"schemeIdUri,attr"`
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

type Representation struct {
   Bandwidth         uint64 `xml:"bandwidth,attr"`
   BaseUrl           *BaseUrl   `xml:"BaseURL"`
   Codecs            string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height            uint64 `xml:"height,attr"`
   Id                string `xml:"id,attr"`
   MimeType          string `xml:"mimeType,attr"`
   SegmentBase       *struct {
      Initialization struct {
         Range Range `xml:"range,attr"`
      }
      IndexRange Range `xml:"indexRange,attr"`
   }
   SegmentTemplate   *SegmentTemplate
   Width             uint64 `xml:"width,attr"`
   adaptation_set    *AdaptationSet
}

func (p Period) get_duration() *Duration {
   if p.Duration != nil {
      return p.Duration
   }
   return p.mpd.MediaPresentationDuration
}

func (s SegmentTemplate) start() uint {
   if s.StartNumber >= 1 {
      return s.StartNumber
   }
   return s.PresentationTimeOffset
}

type SegmentTemplate struct {
   StartNumber            uint   `xml:"startNumber,attr"`
   Duration               uint64 `xml:"duration,attr"`
   Initialization         string `xml:"initialization,attr"`
   Media                  string `xml:"media,attr"`
   PresentationTimeOffset uint   `xml:"presentationTimeOffset,attr"`
   Timescale              uint64 `xml:"timescale,attr"`
   SegmentTimeline        *struct {
      S []struct {
         D uint `xml:"d,attr"` // duration
         R uint `xml:"r,attr"` // repeat
      }
   }
}

func replace_number(format string, a uint) string {
   format = strings.NewReplacer(
      "$Number$", "%d",
      "$Number%02d$", "%02d",
      "$Number%03d$", "%03d",
      "$Number%04d$", "%04d",
      "$Number%05d$", "%05d",
      "$Number%06d$", "%06d",
      "$Number%07d$", "%07d",
      "$Number%08d$", "%08d",
      "$Number%09d$", "%09d",
   ).Replace(format)
   return fmt.Sprintf(format, a)
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#addressing-simple-to-explicit
func (s SegmentTemplate) segment_count(seconds float64) uint64 {
   seconds /= float64(s.Duration) / float64(s.get_timescale())
   return uint64(math.Ceil(seconds))
}

func replace_time(format string, a uint) string {
   format = strings.Replace(format, "$Time$", "%d", 1)
   return fmt.Sprintf(format, a)
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#timing-sampletimeline
func (s SegmentTemplate) get_timescale() uint64 {
   if s.Timescale >= 1 {
      return s.Timescale
   }
   return 1
}

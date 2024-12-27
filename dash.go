package dash

import (
   "fmt"
   "math"
   "strconv"
   "strings"
   "time"
)

type ContentProtection struct {
   Pssh        string   `xml:"pssh"`
   SchemeIdUri string `xml:"schemeIdUri,attr"`
}

type Timescale float64

type Representation struct {
   Bandwidth         uint64   `xml:"bandwidth,attr"`
   BaseUrl           string `xml:"BaseURL"`
   Codecs            string   `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height            uint64 `xml:"height,attr"`
   Id                string `xml:"id,attr"`
   MimeType          string `xml:"mimeType,attr"`
   Width             uint64 `xml:"width,attr"`
   SegmentBase       *struct {
      Initialization struct {
         Range string `xml:"range,attr"`
      }
      IndexRange string `xml:"indexRange,attr"`
   }
   SegmentTemplate   *SegmentTemplate
   adaptation_set    *AdaptationSet
}

type AdaptationSet struct {
   Codecs            string `xml:"codecs,attr"`
   Height            uint64 `xml:"height,attr"`
   Lang              string `xml:"lang,attr"`
   MaxHeight         int    `xml:"maxHeight,attr"`
   MaxWidth          int    `xml:"maxWidth,attr"`
   MimeType          string `xml:"mimeType,attr"`
   Role              *struct {
      Value string `xml:"value,attr"`
   }
   Width           uint64 `xml:"width,attr"`
   Representation    []Representation
   ContentProtection []ContentProtection
   SegmentTemplate *SegmentTemplate
}

func (r Representation) initialization() string {
   return strings.Replace(
      r.SegmentTemplate.Initialization, "$RepresentationID$", r.Id, 1,
   )
}

// While optional in [DASH], the presence of the @timescale attribute is required
// by the interoperable timing model because the default value of 1 is unlikely
// to match any real-world content and is far more likely to indicate an
// unintentional content authoring error.
// dashif.org/Guidelines-TimingModel#timing-sampletimeline
func (t Timescale) value() float64 {
   if t >= 1 {
      return float64(t)
   }
   return 1
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

func (s *SegmentTemplate) time(value uint) string {
   format := strings.Replace(s.Media, "$Time$", "%d", 1)
   return fmt.Sprintf(format, value)
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

type Mpd struct {
   BaseUrl                   string  `xml:"BaseURL"`
   Period                    []Period
   MediaPresentationDuration Duration `xml:"mediaPresentationDuration,attr"`
}

type Period struct {
   mpd           *Mpd
   BaseUrl       string  `xml:"BaseURL"`
   Id            string    `xml:"id,attr"`
   Duration      Duration `xml:"duration,attr"`
   AdaptationSet []AdaptationSet
}

type SegmentTemplate struct {
   PresentationTimeOffset uint    `xml:"presentationTimeOffset,attr"`
   SegmentTimeline        *struct {
      S []struct {
         D uint `xml:"d,attr"` // duration
         R uint `xml:"r,attr"` // repeat
      }
   }
   StartNumber *uint  `xml:"startNumber,attr"`
   Initialization         string  `xml:"initialization,attr"`
   Media                  string  `xml:"media,attr"`
   Duration               float64 `xml:"duration,attr"`
   Timescale   Timescale `xml:"timescale,attr"`
}

func (d *Duration) UnmarshalText(data []byte) error {
   v, err := time.ParseDuration(strings.ToLower(
      strings.TrimPrefix(string(data), "PT"),
   ))
   if err != nil {
      return err
   }
   *d = Duration(v)
   return nil
}

type Duration time.Duration

// dashif.org/Guidelines-TimingModel#addressing-simple-to-explicit
func (r *Representation) segment_count() uint64 {
   seconds := time.Duration(r.adaptation_set.period.Duration).Seconds()
   seconds /= r.SegmentTemplate.Duration / r.SegmentTemplate.Timescale.value()
   return uint64(math.Ceil(seconds))
}

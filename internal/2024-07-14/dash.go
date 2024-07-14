package dash

import (
   "fmt"
   "math"
   "strings"
   "time"
)

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

type Mpd struct {
   BaseUrl string `xml:"BaseURL"`
   MediaPresentationDuration *Duration `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

func (r Representation) id(value string) string {
   return strings.Replace(value, "$RepresentationID$", r.Id, 1)
}

type Representation struct {
   Bandwidth         uint64 `xml:"bandwidth,attr"`
   BaseUrl           string   `xml:"BaseURL"`
   Codecs            string `xml:"codecs,attr"`
   Height            uint64 `xml:"height,attr"`
   Id                string `xml:"id,attr"`
   MimeType          string `xml:"mimeType,attr"`
   SegmentBase       *struct {
      Initialization struct {
         Range string `xml:"range,attr"`
      }
      IndexRange string `xml:"indexRange,attr"`
   }
   SegmentTemplate   *SegmentTemplate
   Width             uint64 `xml:"width,attr"`
   adaptation_set    *AdaptationSet
}

type AdaptationSet struct {
   Codecs            string `xml:"codecs,attr"`
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

type Duration struct {
   Duration time.Duration
}

type Period struct {
   AdaptationSet []AdaptationSet
   BaseUrl string `xml:"BaseURL"`
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

func (s SegmentTemplate) start() uint {
   if s.StartNumber >= 1 {
      return s.StartNumber
   }
   return s.PresentationTimeOffset
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

func replace_time(format string, a uint) string {
   format = strings.Replace(format, "$Time$", "%d", 1)
   return fmt.Sprintf(format, a)
}

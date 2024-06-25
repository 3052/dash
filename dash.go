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

func (m *Mpd) Unmarshal(text []byte) error {
   return xml.Unmarshal(text, m)
}

func (m Mpd) Representation() chan Representation {
   channel := make(chan Representation)
   go func() {
      for _, period := range m.Period {
         period.mpd = &m
         for _, adapt := range period.AdaptationSet {
            adapt.period = &period
            for _, represent := range adapt.Representation {
               represent.adaptation_set = &adapt
               channel <- represent
            }
         }
      }
      close(channel)
   }()
   return channel
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

func (a AdaptationSet) GetPeriod() *Period {
   return a.period
}

type ContentProtection struct {
   Pssh        Pssh   `xml:"pssh"`
   SchemeIdUri string `xml:"schemeIdUri,attr"`
}

func (d *Duration) UnmarshalText(text []byte) error {
   var err error
   d.D, err = time.ParseDuration(strings.ToLower(
      strings.TrimPrefix(string(text), "PT"),
   ))
   if err != nil {
      return err
   }
   return nil
}

type Duration struct {
   D time.Duration
}

type Mpd struct {
   BaseUrl                   *Url      `xml:"BaseURL"`
   MediaPresentationDuration *Duration `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

type Period struct {
   AdaptationSet []AdaptationSet
   Duration      *Duration `xml:"duration,attr"`
   Id            string    `xml:"id,attr"`
   mpd           *Mpd
}

func (p Period) GetMpd() *Mpd {
   return p.mpd
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

func (s SegmentTemplate) number(value uint) string {
   s.Media = strings.NewReplacer(
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
   return fmt.Sprintf(s.Media, value)
}

type SegmentTemplate struct {
   Duration               uint64 `xml:"duration,attr"`
   Initialization         string `xml:"initialization,attr"`
   Media                  string `xml:"media,attr"`
   StartNumber            uint   `xml:"startNumber,attr"`
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

func (s SegmentTemplate) time(value uint) string {
   f := strings.Replace(s.Media, "$Time$", "%d", 1)
   return fmt.Sprintf(f, value)
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#timing-sampletimeline
func (s SegmentTemplate) get_timescale() uint64 {
   if s.Timescale >= 1 {
      return s.Timescale
   }
   return 1
}

func (s SegmentTemplate) start() uint {
   if s.PresentationTimeOffset >= 1 {
      return s.PresentationTimeOffset
   }
   return s.StartNumber
}

func (u *Url) UnmarshalText(text []byte) error {
   u.U = new(url.URL)
   return u.U.UnmarshalBinary(text)
}

type Url struct {
   U *url.URL
}

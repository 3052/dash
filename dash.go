package dash

import (
   "encoding/base64"
   "iter"
   "math"
   "net/url"
   "strconv"
   "strings"
   "time"
)

func replace(s *string, from, to string) {
   *s = strings.Replace(*s, from, to, 1)
}

func (a *AdaptationSet) set(p *Period) {
   a.period = p
}

type AdaptationSet struct {
   Codecs            *string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height            *int64  `xml:"height,attr"`
   Lang              string  `xml:"lang,attr"`
   MimeType          *string `xml:"mimeType,attr"`
   Representation    []Representation
   Role              *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
   Width           *int64 `xml:"width,attr"`
   period          *Period
}

type ContentProtection struct {
   Pssh        Pssh        `xml:"pssh"`
   SchemeIdUri SchemeIdUri `xml:"schemeIdUri,attr"`
}

func (d *Duration) UnmarshalText(data []byte) error {
   var err error
   d.D, err = time.ParseDuration(strings.ToLower(
      strings.TrimPrefix(string(data), "PT"),
   ))
   if err != nil {
      return err
   }
   return nil
}

type Duration struct {
   D time.Duration
}

type Initialization struct {
   S string
}

func (i *Initialization) UnmarshalText(data []byte) error {
   i.S = string(data)
   return nil
}

func (i Initialization) Url(r *Representation) (*url.URL, error) {
   replace(&i.S, "$RepresentationID$", r.Id)
   u, err := url.Parse(i.S)
   if err != nil {
      return nil, err
   }
   if r.BaseUrl != nil {
      u = r.BaseUrl.Url.ResolveReference(u)
   }
   return u, nil
}

type Mpd struct {
   BaseUrl                   *Url      `xml:"BaseURL"`
   MediaPresentationDuration *Duration `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

func (m *Mpd) representation() iter.Seq[Representation] {
   return func(yield func(Representation) bool) {
      for _, p := range m.Period {
         p.set(m)
         for _, adapt := range p.AdaptationSet {
            adapt.set(&p)
            for _, represent := range adapt.Representation {
               represent.set(&adapt)
               if !yield(represent) {
                  return
               }
            }
         }
      }
   }
}

// dashif.org/Guidelines-TimingModel#addressing-simple-to-explicit
func (p *Period) segment_count(template *SegmentTemplate) float64 {
   return math.Ceil(
      *template.Timescale * p.Duration.D.Seconds() / template.Duration,
   )
}

type Period struct {
   AdaptationSet []AdaptationSet
   BaseUrl       *Url      `xml:"BaseURL"`
   Duration      *Duration `xml:"duration,attr"`
   Id            string    `xml:"id,attr"`
   mpd           *Mpd
}

func (p *Period) set(media *Mpd) {
   p.mpd = media
   if v := p.mpd.BaseUrl; v != nil {
      if p.BaseUrl == nil {
         p.BaseUrl = &Url{&url.URL{}}
      }
      p.BaseUrl.Url = v.Url.ResolveReference(p.BaseUrl.Url)
   }
   if p.Duration == nil {
      p.Duration = p.mpd.MediaPresentationDuration
   }
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

// the current testdata always has `-`, so lets assume for now
func (r *Range) UnmarshalText(data []byte) error {
   before, after, _ := strings.Cut(string(data), "-")
   var err error
   (*r)[0], err = strconv.ParseUint(before, 10, 64)
   if err != nil {
      return err
   }
   (*r)[1], err = strconv.ParseUint(after, 10, 64)
   if err != nil {
      return err
   }
   return nil
}

func (r Range) MarshalText() ([]byte, error) {
   b := strconv.AppendUint(nil, r[0], 10)
   b = append(b, '-')
   return strconv.AppendUint(b, r[1], 10), nil
}

// SegmentIndexBox uses:
// unsigned int(32) subsegment_duration;
// but range values can exceed 32 bits
type Range [2]uint64

func (s SchemeIdUri) Widevine() bool {
   return s == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed"
}

type SchemeIdUri string

type SegmentTemplate struct {
   Initialization         Initialization `xml:"initialization,attr"`
   Media                  Media          `xml:"media,attr"`
   Duration               float64         `xml:"duration,attr"`
   Timescale              *float64           `xml:"timescale,attr"`
   StartNumber            *int           `xml:"startNumber,attr"`
   PresentationTimeOffset int            `xml:"presentationTimeOffset,attr"`
   SegmentTimeline        *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R int `xml:"r,attr"` // repeat
      }
   }
}

func (s *SegmentTemplate) set() {
   // dashif.org/Guidelines-TimingModel#addressing-simple
   if s.StartNumber == nil {
      value := 1
      s.StartNumber = &value
   }
   // dashif.org/Guidelines-TimingModel#timing-sampletimeline
   if s.Timescale == nil {
      var value float64 = 1
      s.Timescale = &value
   }
}

type Url struct {
   Url *url.URL
}

func (b *Url) UnmarshalText(data []byte) error {
   b.Url = &url.URL{}
   return b.Url.UnmarshalBinary(data)
}

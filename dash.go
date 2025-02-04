package dash

import (
   "encoding/base64"
   "encoding/xml"
   "fmt"
   "iter"
   "math"
   "net/url"
   "strconv"
   "strings"
   "time"
)

func (r *Range) UnmarshalText(data []byte) error {
   before, after, _ := strings.Cut(string(data), "-")
   var err error
   if before != "" {
      (*r)[0], err = strconv.ParseUint(before, 10, 64)
      if err != nil {
         return err
      }
   }
   if before != "" {
      if after == "" {
         return nil
      }
   }
   (*r)[1], err = strconv.ParseUint(after, 10, 64)
   if err != nil {
      return err
   }
   return nil
}

func (r Range) MarshalText() ([]byte, error) {
   data := strconv.AppendUint(nil, r[0], 10)
   data = append(data, '-')
   if r[1] != 0 {
      data = strconv.AppendUint(data, r[1], 10)
   }
   return data, nil
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
   StartNumber *int `xml:"startNumber,attr"`
   // This can be any frequency but typically is the media clock frequency of
   // one of the media streams (or a positive integer multiple thereof).
   Timescale              *uint64         `xml:"timescale,attr"`
   Media                  Media           `xml:"media,attr"`
   Initialization         *Initialization `xml:"initialization,attr"`
   Duration               float64         `xml:"duration,attr"`
   PresentationTimeOffset int             `xml:"presentationTimeOffset,attr"`
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
      start := 1
      s.StartNumber = &start
   }
   // dashif.org/Guidelines-TimingModel#timing-sampletimeline
   if s.Timescale == nil {
      var scale uint64 = 1
      s.Timescale = &scale
   }
}

type Url struct {
   Url *url.URL
}

func (u *Url) UnmarshalText(data []byte) error {
   var err error
   if u.Url != nil {
      u.Url, err = u.Url.Parse(string(data))
   } else {
      u.Url = &url.URL{}
      err = u.Url.UnmarshalBinary(data)
   }
   if err != nil {
      return err
   }
   return nil
}

func replace(data *string, from, to string) {
   *data = strings.Replace(*data, from, to, 1)
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

func (i *Initialization) UnmarshalText(data []byte) error {
   i.S = string(data)
   return nil
}

type Initialization struct {
   S string
}

func (m Media) time() bool {
   return strings.Contains(m.S, "$Time$")
}

func (m *Media) UnmarshalText(data []byte) error {
   m.S = string(data)
   return nil
}

type Media struct {
   S string
}

func (m *Mpd) Unmarshal(data []byte) error {
   return xml.Unmarshal(data, m)
}

type Mpd struct {
   BaseUrl                   *Url      `xml:"BaseURL"`
   MediaPresentationDuration *Duration `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

// dashif.org/Guidelines-TimingModel#addressing-simple-to-explicit
func (p *Period) segment_count(template *SegmentTemplate) float64 {
   return math.Ceil(
      p.Duration.D.Seconds() * float64(*template.Timescale) / template.Duration,
   )
}

type Period struct {
   AdaptationSet []AdaptationSet
   BaseUrl       *Url      `xml:"BaseURL"`
   Duration      *Duration `xml:"duration,attr"`
   Id            string    `xml:"id,attr"`
   mpd           *Mpd
}

func (u *ListUrl) UnmarshalText(data []byte) error {
   u.S = string(data)
   return nil
}

type ListUrl struct {
   S string
}

func (a *AdaptationSet) set(period0 *Period) {
   a.period = period0
}

func (p *Period) set(media *Mpd) {
   p.mpd = media
   if base := p.mpd.BaseUrl; base != nil {
      if p.BaseUrl == nil {
         p.BaseUrl = &Url{&url.URL{}}
      }
      p.BaseUrl.Url = base.Url.ResolveReference(p.BaseUrl.Url)
   }
   if p.Duration == nil {
      p.Duration = p.mpd.MediaPresentationDuration
   }
}

func (m *Mpd) Representation() iter.Seq[Representation] {
   id := map[string]struct{}{}
   return func(yield func(Representation) bool) {
      for _, period0 := range m.Period {
         for _, adapt := range period0.AdaptationSet {
            for _, represent := range adapt.Representation {
               _, ok := id[represent.Id]
               if !ok {
                  if adapt.period == nil {
                     period0.set(m)
                     adapt.set(&period0)
                  }
                  represent.set(&adapt)
                  if !yield(represent) {
                     return
                  }
                  id[represent.Id] = struct{}{}
               }
            }
         }
      }
   }
}

func (u ListUrl) Url(r *Representation) (*url.URL, error) {
   if r.BaseUrl != nil {
      return r.BaseUrl.Url.Parse(u.S)
   }
   return url.Parse(u.S)
}

func (i Initialization) Url(r *Representation) (*url.URL, error) {
   replace(&i.S, "$RepresentationID$", r.Id)
   url0, err := url.Parse(i.S)
   if err != nil {
      return nil, err
   }
   if r.BaseUrl != nil {
      url0 = r.BaseUrl.Url.ResolveReference(url0)
   }
   return url0, nil
}

func (m Media) Url(r *Representation, i int) (*url.URL, error) {
   replace(&m.S, "$RepresentationID$", r.Id)
   if m.time() {
      replace(&m.S, "$Time$", fmt.Sprint(i))
   } else {
      replace(&m.S, "$Number$", fmt.Sprint(i))
      replace(&m.S, "$Number%02d$", fmt.Sprintf("%02d", i))
      replace(&m.S, "$Number%03d$", fmt.Sprintf("%03d", i))
      replace(&m.S, "$Number%04d$", fmt.Sprintf("%04d", i))
      replace(&m.S, "$Number%05d$", fmt.Sprintf("%05d", i))
      replace(&m.S, "$Number%06d$", fmt.Sprintf("%06d", i))
      replace(&m.S, "$Number%07d$", fmt.Sprintf("%07d", i))
      replace(&m.S, "$Number%08d$", fmt.Sprintf("%08d", i))
      replace(&m.S, "$Number%09d$", fmt.Sprintf("%09d", i))
   }
   url0, err := url.Parse(m.S)
   if err != nil {
      return nil, err
   }
   if r.BaseUrl != nil {
      url0 = r.BaseUrl.Url.ResolveReference(url0)
   }
   return url0, nil
}

type ContentProtection struct {
   Pssh        string      `xml:"pssh"`
   SchemeIdUri SchemeIdUri `xml:"schemeIdUri,attr"`
}

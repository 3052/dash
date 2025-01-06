package dash

import (
   "encoding/base64"
   "iter"
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

func (m *Mpd) Representation() iter.Seq[Representation] {
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

type Mpd struct {
   BaseUrl                   *Url      `xml:"BaseURL"`
   MediaPresentationDuration *Duration `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

func (m *Mpd) Set(base *url.URL) {
   if m.BaseUrl == nil {
      m.BaseUrl = &Url{&url.URL{}}
   }
   m.BaseUrl.Url = base.ResolveReference(m.BaseUrl.Url)
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
   b := strconv.AppendUint(nil, r[0], 10)
   b = append(b, '-')
   if r[1] != 0 {
      b = strconv.AppendUint(b, r[1], 10)
   }
   return b, nil
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
   Initialization Initialization `xml:"initialization,attr"`
   Media          Media          `xml:"media,attr"`
   Duration       float64        `xml:"duration,attr"`
   // This can be any frequency but typically is the media clock frequency of
   // one of the media streams (or a positive integer multiple thereof).
   Timescale              *uint64 `xml:"timescale,attr"`
   StartNumber            *int    `xml:"startNumber,attr"`
   PresentationTimeOffset int     `xml:"presentationTimeOffset,attr"`
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
      var value uint64 = 1
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

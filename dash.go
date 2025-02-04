package dash

import (
   "encoding/xml"
   "fmt"
   "iter"
   "math"
   "net/url"
   "strconv"
   "strings"
   "time"
)

func (m Media) time() bool {
   return strings.Contains(string(m), "$Time$")
}

func (r *Representation) Segment() iter.Seq[int] {
   template := r.SegmentTemplate
   var address int
   if template.Media.time() {
      address = template.PresentationTimeOffset
   } else {
      address = *template.StartNumber
   }
   return func(yield func(int) bool) {
      if template.SegmentTimeline != nil {
         for _, segment := range template.SegmentTimeline.S {
            for range 1 + segment.R {
               if !yield(address) {
                  return
               }
               if template.Media.time() {
                  address += segment.D
               } else {
                  address++
               }
            }
         }
      } else {
         segment_count := r.adaptation_set.period.segment_count(template)
         for range int64(segment_count) {
            if !yield(address) {
               return
            }
            address++
         }
      }
   }
}

type SegmentTemplate struct {
   StartNumber *int `xml:"startNumber,attr"`
   // This can be any frequency but typically is the media clock frequency of
   // one of the media streams (or a positive integer multiple thereof).
   Timescale              *uint64         `xml:"timescale,attr"`
   Media                  Media           `xml:"media,attr"`
   Initialization         Initialization `xml:"initialization,attr"`
   Duration               float64         `xml:"duration,attr"`
   PresentationTimeOffset int             `xml:"presentationTimeOffset,attr"`
   SegmentTimeline        *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R int `xml:"r,attr"` // repeat
      }
   }
}

func (i Initialization) Url(r *Representation) (*url.URL, error) {
   raw := replace(string(i), "$RepresentationID$", r.Id)
   url0, err := url.Parse(raw)
   if err != nil {
      return nil, err
   }
   if r.BaseUrl != nil {
      url0 = r.BaseUrl.Url.ResolveReference(url0)
   }
   return url0, nil
}

type Media string

func (m Media) Url(r *Representation, index int) (*url.URL, error) {
   raw := replace(string(m), "$RepresentationID$", r.Id)
   if m.time() {
      raw = replace(raw, "$Time$", fmt.Sprint(index))
   } else {
      raw = replace(raw, "$Number$", fmt.Sprint(index))
      raw = replace(raw, "$Number%02d$", fmt.Sprintf("%02d", index))
      raw = replace(raw, "$Number%03d$", fmt.Sprintf("%02d", index))
      raw = replace(raw, "$Number%04d$", fmt.Sprintf("%02d", index))
      raw = replace(raw, "$Number%05d$", fmt.Sprintf("%02d", index))
      raw = replace(raw, "$Number%06d$", fmt.Sprintf("%02d", index))
      raw = replace(raw, "$Number%07d$", fmt.Sprintf("%02d", index))
      raw = replace(raw, "$Number%08d$", fmt.Sprintf("%02d", index))
      raw = replace(raw, "$Number%09d$", fmt.Sprintf("%02d", index))
   }
   url0, err := url.Parse(raw)
   if err != nil {
      return nil, err
   }
   if r.BaseUrl != nil {
      url0 = r.BaseUrl.Url.ResolveReference(url0)
   }
   return url0, nil
}

func (m *Mpd) Unmarshal(data []byte) error {
   return xml.Unmarshal(data, m)
}

type Mpd struct {
   BaseUrl                   *Url      `xml:"BaseURL"`
   MediaPresentationDuration *Duration `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

type Period struct {
   AdaptationSet []AdaptationSet
   BaseUrl       *Url      `xml:"BaseURL"`
   Duration      *Duration `xml:"duration,attr"`
   Id            string    `xml:"id,attr"`
   mpd           *Mpd
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

type ContentProtection struct {
   Pssh        string `xml:"pssh"`
   SchemeIdUri string `xml:"schemeIdUri,attr"`
}

func (r *Representation) String() string {
   var b []byte
   if r.Width != nil {
      b = append(b, "width = "...)
      b = strconv.AppendInt(b, *r.Width, 10)
   }
   if r.Height != nil {
      if b != nil {
         b = append(b, '\n')
      }
      b = append(b, "height = "...)
      b = strconv.AppendInt(b, *r.Height, 10)
   }
   if b != nil {
      b = append(b, '\n')
   }
   b = append(b, "bandwidth = "...)
   b = strconv.AppendInt(b, int64(r.Bandwidth), 10)
   if r.Codecs != nil {
      b = append(b, "\ncodecs = "...)
      b = append(b, *r.Codecs...)
   }
   b = append(b, "\nmimeType = "...)
   b = append(b, *r.MimeType...)
   if role := r.adaptation_set.Role; role != nil {
      b = append(b, "\nrole = "...)
      b = append(b, role.Value...)
   }
   if lang := r.adaptation_set.Lang; lang != "" {
      b = append(b, "\nlang = "...)
      b = append(b, lang...)
   }
   if id := r.adaptation_set.period.Id; id != "" {
      b = append(b, "\nperiod = "...)
      b = append(b, id...)
   }
   b = append(b, "\nid = "...)
   b = append(b, r.Id...)
   return string(b)
}

type Representation struct {
   Bandwidth         int     `xml:"bandwidth,attr"`
   BaseUrl           *Url    `xml:"BaseURL"`
   Codecs            *string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height            *int64  `xml:"height,attr"`
   Id                string  `xml:"id,attr"`
   MimeType          *string `xml:"mimeType,attr"`
   SegmentBase       *struct {
      Initialization struct {
         Range Range `xml:"range,attr"`
      }
      IndexRange Range `xml:"indexRange,attr"`
   }
   SegmentList *struct {
      Initialization struct {
         SourceUrl ListUrl `xml:"sourceURL,attr"`
      }
      SegmentUrl []struct {
         Media ListUrl `xml:"media,attr"`
      } `xml:"SegmentURL"`
   }
   SegmentTemplate *SegmentTemplate
   Width           *int64 `xml:"width,attr"`
   adaptation_set  *AdaptationSet
}

func (r *Representation) Representation() iter.Seq[Representation] {
   return func(yield func(Representation) bool) {
      for _, period0 := range r.adaptation_set.period.mpd.Period {
         for _, adapt := range period0.AdaptationSet {
            for _, represent := range adapt.Representation {
               if represent.Id == r.Id {
                  if adapt.period == nil {
                     period0.set(r.adaptation_set.period.mpd)
                     adapt.set(&period0)
                  }
                  represent.set(&adapt)
                  if !yield(represent) {
                     return
                  }
               }
            }
         }
      }
   }
}

///

func (r *Representation) set(adapt *AdaptationSet) {
   r.adaptation_set = adapt
   if base := r.adaptation_set.period.BaseUrl; base != nil {
      if r.BaseUrl == nil {
         r.BaseUrl = &Url{&url.URL{}}
      }
      r.BaseUrl.Url = base.Url.ResolveReference(r.BaseUrl.Url)
   }
   if r.Codecs == nil {
      r.Codecs = r.adaptation_set.Codecs
   }
   if len(r.ContentProtection) == 0 {
      r.ContentProtection = r.adaptation_set.ContentProtection
   }
   if r.Height == nil {
      r.Height = r.adaptation_set.Height
   }
   if r.MimeType == nil {
      r.MimeType = r.adaptation_set.MimeType
   }
   if r.SegmentTemplate == nil {
      r.SegmentTemplate = r.adaptation_set.SegmentTemplate
   }
   if r.SegmentTemplate != nil {
      r.SegmentTemplate.set()
   }
   if r.Width == nil {
      r.Width = r.adaptation_set.Width
   }
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

// SegmentIndexBox uses:
// unsigned int(32) subsegment_duration;
// but range values can exceed 32 bits
type Range [2]uint64

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

func (r Range) MarshalText() ([]byte, error) {
   data := strconv.AppendUint(nil, r[0], 10)
   data = append(data, '-')
   if r[1] != 0 {
      data = strconv.AppendUint(data, r[1], 10)
   }
   return data, nil
}

// dashif.org/Guidelines-TimingModel#addressing-simple-to-explicit
func (p *Period) segment_count(template *SegmentTemplate) float64 {
   return math.Ceil(
      p.Duration.D.Seconds() * float64(*template.Timescale) / template.Duration,
   )
}

func (a *AdaptationSet) set(period0 *Period) {
   a.period = period0
}

type Initialization string

func replace(s, old, new0 string) string {
   return strings.Replace(s, old, new0, 1)
}

func (p *Period) set(mpd0 *Mpd) {
   p.mpd = mpd0
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

type ListUrl string

func (u ListUrl) Url(r *Representation) (*url.URL, error) {
   if r.BaseUrl != nil {
      return r.BaseUrl.Url.Parse(string(u))
   }
   return url.Parse(string(u))
}

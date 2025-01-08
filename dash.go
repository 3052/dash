package dash

import (
   "encoding/base64"
   "encoding/xml"
   "iter"
   "math"
   "net/url"
   "strconv"
   "strings"
   "time"
)

type Representation struct {
   SegmentTemplate *SegmentTemplate
   SegmentList       *struct {
      Initialization struct {
         SourceUrl string `xml:"sourceURL,attr"`
      }
      SegmentUrl []struct {
         Media string `xml:"media,attr"`
      } `xml:"SegmentURL"`
   }
   Bandwidth         int     `xml:"bandwidth,attr"`
   BaseUrl           *Url    `xml:"BaseURL"`
   Codecs            *string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height            *int64  `xml:"height,attr"`
   Id                string  `xml:"id,attr"`
   MimeType          *string `xml:"mimeType,attr"`
   SegmentBase       *SegmentBase
   Width           *int64 `xml:"width,attr"`
   adaptation_set  *AdaptationSet
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

func (r *Representation) Representation() iter.Seq[Representation] {
   return func(yield func(Representation) bool) {
      for _, p := range r.adaptation_set.period.mpd.Period {
         for _, adapt := range p.AdaptationSet {
            for _, represent := range adapt.Representation {
               if represent.Id == r.Id {
                  if adapt.period == nil {
                     p.set(r.adaptation_set.period.mpd)
                     adapt.set(&p)
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

func (r *Representation) set(adapt *AdaptationSet) {
   r.adaptation_set = adapt
   if v := r.adaptation_set.period.BaseUrl; v != nil {
      if r.BaseUrl == nil {
         r.BaseUrl = &Url{&url.URL{}}
      }
      r.BaseUrl.Url = v.Url.ResolveReference(r.BaseUrl.Url)
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

func (m *Mpd) Unmarshal(data []byte) error {
   return xml.Unmarshal(data, m)
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

func (b *Url) UnmarshalText(data []byte) error {
   if b.Url == nil {
      b.Url = &url.URL{}
      return b.Url.UnmarshalBinary(data)
   }
   var err error
   b.Url, err = b.Url.Parse(string(data))
   if err != nil {
      return err
   }
   return nil
}

type Url struct {
   Url *url.URL
}

type SegmentBase struct {
   Initialization struct {
      Range Range `xml:"range,attr"`
   }
   IndexRange Range `xml:"indexRange,attr"`
}

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

func (m *Mpd) Representation() iter.Seq[Representation] {
   id := map[string]struct{}{}
   return func(yield func(Representation) bool) {
      for _, p := range m.Period {
         for _, adapt := range p.AdaptationSet {
            for _, represent := range adapt.Representation {
               _, ok := id[represent.Id]
               if !ok {
                  if adapt.period == nil {
                     p.set(m)
                     adapt.set(&p)
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

type Mpd struct {
   BaseUrl                   *Url      `xml:"BaseURL"`
   MediaPresentationDuration *Duration `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

package dash

import (
   "encoding/xml"
   "fmt"
   "iter"
   "math"
   "net/url"
   "strings"
   "time"
)

func (r *Representation) String() string {
   var b []byte
   if r.Width != nil {
      b = fmt.Appendln(b, "width =", *r.Width)
   }
   if r.Height != nil {
      b = fmt.Appendln(b, "height =", *r.Height)
   }
   b = fmt.Appendln(b, "bandwidth =", r.Bandwidth)
   if r.Codecs != nil {
      b = fmt.Appendln(b, "codecs =", *r.Codecs)
   }
   b = fmt.Appendln(b, "mimeType =", *r.MimeType)
   if role := r.adaptation_set.Role; role != nil {
      b = fmt.Appendln(b, "role =", role.Value)
   }
   if lang := r.adaptation_set.Lang; lang != "" {
      b = fmt.Appendln(b, "lang =", lang)
   }
   if id := r.adaptation_set.period.Id; id != "" {
      b = fmt.Appendln(b, "period =", id)
   }
   b = fmt.Append(b, "id = ", r.Id)
   return string(b)
}

func (u *Url) UnmarshalText(data []byte) error {
   (*u)[0] = &url.URL{}
   return u[0].UnmarshalBinary(data)
}

type Url [1]*url.URL

func (d *Duration) UnmarshalText(data []byte) error {
   var err error
   (*d)[0], err = time.ParseDuration(strings.ToLower(
      strings.TrimPrefix(string(data), "PT"),
   ))
   if err != nil {
      return err
   }
   return nil
}

type Duration [1]time.Duration

func replace(s, old, new1 string) string {
   return strings.Replace(s, old, new1, 1)
}

// SegmentTemplate, SegmentUrl
// dashif.org/Guidelines-TimingModel#addressing-explicit
func using_time(data string) bool {
   return strings.Contains(data, "$Time$")
}

func (a *AdaptationSet) set(period1 *Period) {
   a.period = period1
}

type ContentProtection struct {
   Pssh        string `xml:"pssh"`
   SchemeIdUri string `xml:"schemeIdUri,attr"`
}

func (m *Mpd) Set(url2 *url.URL) {
   if m.BaseUrl[0] == nil {
      m.BaseUrl[0] = &url.URL{}
   }
   m.BaseUrl[0] = url2.ResolveReference(m.BaseUrl[0])
}

func (m *Mpd) Unmarshal(data []byte) error {
   return xml.Unmarshal(data, m)
}

type Period struct {
   AdaptationSet []AdaptationSet
   BaseUrl       Url       `xml:"BaseURL"`
   Duration      *Duration `xml:"duration,attr"`
   Id            string    `xml:"id,attr"`
   mpd           *Mpd
}

func (p *Period) set(mpd1 *Mpd) {
   p.mpd = mpd1
   if base := p.mpd.BaseUrl[0]; base != nil {
      if p.BaseUrl[0] == nil {
         p.BaseUrl[0] = &url.URL{}
      }
      p.BaseUrl[0] = base.ResolveReference(p.BaseUrl[0])
   }
   if p.Duration == nil {
      p.Duration = &p.mpd.MediaPresentationDuration
   }
}

func (r *Representation) set(adapt *AdaptationSet) {
   r.adaptation_set = adapt
   if base := r.adaptation_set.period.BaseUrl[0]; base != nil {
      if r.BaseUrl[0] == nil {
         r.BaseUrl[0] = &url.URL{}
      }
      r.BaseUrl[0] = base.ResolveReference(r.BaseUrl[0])
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
      r.MimeType = &r.adaptation_set.MimeType
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

func (m *Mpd) Representation() iter.Seq[Representation] {
   return func(yield func(Representation) bool) {
      id := map[string]struct{}{}
      for _, period1 := range m.Period {
         for _, adapt := range period1.AdaptationSet {
            for _, represent := range adapt.Representation {
               _, ok := id[represent.Id]
               if !ok {
                  if adapt.period == nil {
                     period1.set(m)
                     adapt.set(&period1)
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
   BaseUrl                   Url      `xml:"BaseURL"`
   MediaPresentationDuration Duration `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

func (r *Representation) Representation() iter.Seq[Representation] {
   return func(yield func(Representation) bool) {
      for _, period1 := range r.adaptation_set.period.mpd.Period {
         for _, adapt := range period1.AdaptationSet {
            for _, represent := range adapt.Representation {
               if represent.Id == r.Id {
                  if adapt.period == nil {
                     period1.set(r.adaptation_set.period.mpd)
                     adapt.set(&period1)
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

// SegmentTemplate
func (r *Representation) Media(data string, index int) (*url.URL, error) {
   data = replace(data, "$RepresentationID$", r.Id)
   if using_time(data) {
      data = replace(data, "$Time$", fmt.Sprint(index))
   } else {
      data = replace(data, "$Number$", fmt.Sprint(index))
      data = replace(data, "$Number%02d$", fmt.Sprintf("%02d", index))
      data = replace(data, "$Number%03d$", fmt.Sprintf("%03d", index))
      data = replace(data, "$Number%04d$", fmt.Sprintf("%04d", index))
      data = replace(data, "$Number%05d$", fmt.Sprintf("%05d", index))
      data = replace(data, "$Number%06d$", fmt.Sprintf("%06d", index))
      data = replace(data, "$Number%07d$", fmt.Sprintf("%07d", index))
      data = replace(data, "$Number%08d$", fmt.Sprintf("%08d", index))
      data = replace(data, "$Number%09d$", fmt.Sprintf("%09d", index))
   }
   url2, err := url.Parse(data)
   if err != nil {
      return nil, err
   }
   if r.BaseUrl[0] != nil {
      url2 = r.BaseUrl[0].ResolveReference(url2)
   }
   return url2, nil
}

// SegmentTemplate
func (r *Representation) Initialization(data string) (*url.URL, error) {
   data = replace(data, "$RepresentationID$", r.Id)
   url2, err := url.Parse(data)
   if err != nil {
      return nil, err
   }
   if r.BaseUrl[0] != nil {
      url2 = r.BaseUrl[0].ResolveReference(url2)
   }
   return url2, nil
}

// SegmentList
func (r *Representation) List(data string) (*url.URL, error) {
   url2, err := url.Parse(data)
   if err != nil {
      return nil, err
   }
   if r.BaseUrl[0] != nil {
      url2 = r.BaseUrl[0].ResolveReference(url2)
   }
   return url2, nil
}

type AdaptationSet struct {
   ContentProtection []ContentProtection
   Lang              string `xml:"lang,attr"`
   MimeType          string `xml:"mimeType,attr"`
   Representation    []Representation
   Role              *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
   period          *Period
   // pointers for Representation.String
   Codecs *string `xml:"codecs,attr"`
   Height *int    `xml:"height,attr"`
   Width  *int    `xml:"width,attr"`
}

type Representation struct {
   Bandwidth         int     `xml:"bandwidth,attr"`
   BaseUrl           Url     `xml:"BaseURL"`
   Codecs            *string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Id                string  `xml:"id,attr"`
   MimeType          *string `xml:"mimeType,attr"`
   SegmentBase       *struct {
      Initialization struct {
         Range string `xml:"range,attr"`
      }
      IndexRange string `xml:"indexRange,attr"`
   }
   SegmentList *struct {
      Initialization struct {
         SourceUrl string `xml:"sourceURL,attr"`
      }
      SegmentUrl []struct {
         Media string `xml:"media,attr"`
      } `xml:"SegmentURL"`
   }
   SegmentTemplate *SegmentTemplate
   adaptation_set  *AdaptationSet
   Width           *int `xml:"width,attr"`
   Height          *int `xml:"height,attr"`
}

// SegmentTemplate
func (r *Representation) Segment() iter.Seq[int] {
   template := r.SegmentTemplate
   var address int
   if using_time(template.Media) {
      address = template.PresentationTimeOffset
   } else {
      address = *template.StartNumber
   }
   return func(yield func(int) bool) {
      if template.EndNumber >= 1 {
         for address <= template.EndNumber {
            if !yield(address) {
               return
            }
            address++
         }
      } else if template.SegmentTimeline != nil {
         for _, segment := range template.SegmentTimeline.S {
            for range 1 + segment.R {
               if !yield(address) {
                  return
               }
               if using_time(template.Media) {
                  address += segment.D
               } else {
                  address++
               }
            }
         }
      } else {
         for range r.adaptation_set.period.segment_count(template) {
            if !yield(address) {
               return
            }
            address++
         }
      }
   }
}

// dashif.org/Guidelines-TimingModel#addressing-simple-to-explicit
// SegmentCount = Ceil((AsSeconds(Period@duration)) /
// (SegmentTemplate@duration / SegmentTemplate@timescale))
func (p *Period) segment_count(template *SegmentTemplate) int64 {
   // amc.mpd
   // draken.mpd
   // kanopy.mpd
   // max.mpd
   // paramount.mpd
   duration1 := float64(template.Duration) / float64(*template.Timescale)
   return int64(math.Ceil(p.Duration[0].Seconds() / duration1))
}

type SegmentTemplate struct {
   EndNumber              int    `xml:"endNumber,attr"`
   Initialization         string `xml:"initialization,attr"`
   Media                  string `xml:"media,attr"`
   PresentationTimeOffset int    `xml:"presentationTimeOffset,attr"`
   SegmentTimeline        *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R int `xml:"r,attr"` // repeat
      }
   }
   StartNumber *int `xml:"startNumber,attr"`
   Duration    uint `xml:"duration,attr"`
   // This can be any frequency but typically is the media clock frequency of
   // one of the media streams (or a positive integer multiple thereof).
   Timescale *int `xml:"timescale,attr"`
}

func (s *SegmentTemplate) set() {
   // dashif.org/Guidelines-TimingModel#addressing-simple
   if s.StartNumber == nil {
      start := 1
      s.StartNumber = &start
   }
   // dashif.org/Guidelines-TimingModel#timing-sampletimeline
   if s.Timescale == nil {
      scale := 1
      s.Timescale = &scale
   }
}

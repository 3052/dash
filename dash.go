package dash

import (
   "fmt"
   "iter"
   "math"
   "net/url"
   "strings"
   "time"
)

func (d *Duration) UnmarshalText(data []byte) error {
   var err error
   d[0], err = time.ParseDuration(strings.ToLower(
      strings.TrimPrefix(string(data), "PT"),
   ))
   if err != nil {
      return err
   }
   return nil
}

type Mpd struct {
   BaseUrl                   Url      `xml:"BaseURL"`
   MediaPresentationDuration Duration `xml:"mediaPresentationDuration,attr"`
   Period                    []Period
}

type Url [1]*url.URL

type Duration [1]time.Duration

type Period struct {
   BaseUrl       Url       `xml:"BaseURL"`
   Id            string    `xml:"id,attr"`
   Duration      *Duration `xml:"duration,attr"`
   AdaptationSet []AdaptationSet
   mpd           *Mpd
}

type Representation struct {
   Bandwidth         int     `xml:"bandwidth,attr"`
   BaseUrl           Url     `xml:"BaseURL"`
   Codecs            *string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Id                string  `xml:"id,attr"`
   MimeType          *string `xml:"mimeType,attr"`
   Width             *int    `xml:"width,attr"`
   Height            *int    `xml:"height,attr"`
   adaptation_set    *AdaptationSet
   SegmentTemplate   *SegmentTemplate
   SegmentList       *SegmentList
   SegmentBase       *struct {
      Initialization struct {
         Range string `xml:"range,attr"`
      }
      IndexRange string `xml:"indexRange,attr"`
   }
}

func (r *Representation) GetAdaptationSet() *AdaptationSet {
   return r.adaptation_set
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

func (a *AdaptationSet) GetRole() string {
   if a.Role != nil {
      return a.Role.Value
   }
   return ""
}

type ContentProtection struct {
   Pssh        string `xml:"pssh"`
   SchemeIdUri string `xml:"schemeIdUri,attr"`
}

func replace(s, old, newNew string) string {
   return strings.Replace(s, old, newNew, 1)
}

// SegmentTemplate
func (r *Representation) Segment() iter.Seq[int] {
   template := r.SegmentTemplate
   var address int
   if template.Media.time_address() {
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
               if template.Media.time_address() {
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
   if r.SegmentList != nil {
      if r.BaseUrl[0] != nil {
         r.SegmentList.set(r.BaseUrl[0])
      }
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

type SegmentList struct {
   Initialization struct {
      SourceUrl Url `xml:"sourceURL,attr"`
   }
   SegmentUrl []*struct {
      Media Url `xml:"media,attr"`
   } `xml:"SegmentURL"`
}

type SegmentTemplate struct {
   EndNumber              int            `xml:"endNumber,attr"`
   Initialization         Initialization `xml:"initialization,attr"`
   Media                  Media          `xml:"media,attr"`
   PresentationTimeOffset int            `xml:"presentationTimeOffset,attr"`
   SegmentTimeline        *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R int `xml:"r,attr"` // repeat
      }
   }
   StartNumber *int `xml:"startNumber,attr"`
   Duration    int  `xml:"duration,attr"`
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

// SegmentTemplate
// dashif.org/Guidelines-TimingModel#addressing-explicit
// dashif.org/Guidelines-TimingModel#addressing-simple
func (m Media) time_address() bool {
   return strings.Contains(string(m), "$Time$")
}

type Initialization string

type Media string

func (u *Url) UnmarshalText(data []byte) error {
   u[0] = &url.URL{}
   return u[0].UnmarshalBinary(data)
}

func (a *AdaptationSet) set(newPeriod *Period) {
   a.period = newPeriod
}

func (m *Mpd) Set(newUrl *url.URL) {
   if m.BaseUrl[0] == nil {
      m.BaseUrl[0] = &url.URL{}
   }
   m.BaseUrl[0] = newUrl.ResolveReference(m.BaseUrl[0])
}

func (s *SegmentList) set(newUrl *url.URL) {
   s.Initialization.SourceUrl[0] = newUrl.ResolveReference(
      s.Initialization.SourceUrl[0],
   )
   for _, segment := range s.SegmentUrl {
      segment.Media[0] = newUrl.ResolveReference(segment.Media[0])
   }
}

func (i Initialization) Url(represent *Representation) (*url.URL, error) {
   data := replace(string(i), "$RepresentationID$", represent.Id)
   newUrl, err := url.Parse(data)
   if err != nil {
      return nil, err
   }
   if represent.BaseUrl[0] != nil {
      newUrl = represent.BaseUrl[0].ResolveReference(newUrl)
   }
   return newUrl, nil
}

func (m Media) Url(represent *Representation, address int) (*url.URL, error) {
   data := replace(string(m), "$RepresentationID$", represent.Id)
   if m.time_address() {
      data = replace(data, "$Time$", fmt.Sprint(address))
   } else {
      data = replace(data, "$Number$", fmt.Sprint(address))
      data = replace(data, "$Number%02d$", fmt.Sprintf("%02d", address))
      data = replace(data, "$Number%03d$", fmt.Sprintf("%03d", address))
      data = replace(data, "$Number%04d$", fmt.Sprintf("%04d", address))
      data = replace(data, "$Number%05d$", fmt.Sprintf("%05d", address))
      data = replace(data, "$Number%06d$", fmt.Sprintf("%06d", address))
      data = replace(data, "$Number%07d$", fmt.Sprintf("%07d", address))
      data = replace(data, "$Number%08d$", fmt.Sprintf("%08d", address))
      data = replace(data, "$Number%09d$", fmt.Sprintf("%09d", address))
   }
   newUrl, err := url.Parse(data)
   if err != nil {
      return nil, err
   }
   if represent.BaseUrl[0] != nil {
      newUrl = represent.BaseUrl[0].ResolveReference(newUrl)
   }
   return newUrl, nil
}

func (m *Mpd) Representation() iter.Seq[*Representation] {
   return func(yield func(*Representation) bool) {
      id := map[string]struct{}{}
      for _, newPeriod := range m.Period {
         for _, adapt := range newPeriod.AdaptationSet {
            for _, represent := range adapt.Representation {
               _, ok := id[represent.Id]
               if !ok {
                  if adapt.period == nil {
                     newPeriod.set(m)
                     adapt.set(&newPeriod)
                  }
                  represent.set(&adapt)
                  if !yield(&represent) {
                     return
                  }
                  id[represent.Id] = struct{}{}
               }
            }
         }
      }
   }
}

func (r *Representation) Representation() iter.Seq[*Representation] {
   return func(yield func(*Representation) bool) {
      for _, newPeriod := range r.adaptation_set.period.mpd.Period {
         for _, adapt := range newPeriod.AdaptationSet {
            for _, represent := range adapt.Representation {
               if represent.Id == r.Id {
                  if adapt.period == nil {
                     newPeriod.set(r.adaptation_set.period.mpd)
                     adapt.set(&newPeriod)
                  }
                  represent.set(&adapt)
                  if !yield(&represent) {
                     return
                  }
               }
            }
         }
      }
   }
}

// dashif.org/Guidelines-TimingModel#addressing-simple-to-explicit
// SegmentCount = Ceil((AsSeconds(Period@duration)) /
// (SegmentTemplate@duration / SegmentTemplate@timescale))
func (p *Period) segment_count(template *SegmentTemplate) int64 {
   // amc
   // draken
   // kanopy
   // max
   // paramount
   newDuration := float64(template.Duration) / float64(*template.Timescale)
   return int64(math.Ceil(p.Duration[0].Seconds() / newDuration))
}

func (p *Period) set(newMpd *Mpd) {
   p.mpd = newMpd
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

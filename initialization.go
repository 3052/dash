package dash

import (
   "encoding/base64"
   "math"
   "net/url"
)

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
type Url struct {
   Url *url.URL
}

func (b *Url) UnmarshalText(data []byte) error {
   b.Url = &url.URL{}
   return b.Url.UnmarshalBinary(data)
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

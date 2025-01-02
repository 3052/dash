package dash

import "net/url"

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

type Period struct {
   AdaptationSet []AdaptationSet
   BaseUrl       *Url      `xml:"BaseURL"`
   Duration      *Duration `xml:"duration,attr"`
   Id            string    `xml:"id,attr"`
   mpd           *Mpd
}

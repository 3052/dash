package dash

func (p *Period) set() {
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

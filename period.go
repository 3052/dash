package dash

type Period struct {
   AdaptationSet []AdaptationSet
   BaseUrl       *Url   `xml:"BaseURL"`
   Id            string `xml:"id,attr"`
   mpd           *Mpd
}

package dash

import "net/url"

type Initialization struct {
   S string
}

func (i *Initialization) UnmarshalText(data []byte) error {
   return (*Media)(i).UnmarshalText(data)
}

func (i Initialization) Url(r *Representation) (*url.URL, error) {
   return Media(i).Url(r, 0)
}

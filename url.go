package dash

import "net/url"

type Url struct {
   Url url.URL
}

func (b *Url) UnmarshalText(data []byte) error {
   return b.Url.UnmarshalBinary(data)
}

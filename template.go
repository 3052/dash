package dash

import (
   "fmt"
   "net/url"
   "strings"
)

type SegmentTemplate struct {
   Media          Media           `xml:"media,attr"`
   Initialization *Initialization `xml:"initialization,attr"`
   Duration       float64         `xml:"duration,attr"`
   // This can be any frequency but typically is the media clock frequency of
   // one of the media streams (or a positive integer multiple thereof).
   Timescale              *uint64 `xml:"timescale,attr"`
   StartNumber            *int    `xml:"startNumber,attr"`
   PresentationTimeOffset int     `xml:"presentationTimeOffset,attr"`
   SegmentTimeline        *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R int `xml:"r,attr"` // repeat
      }
   }
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

func (m Media) Url(r *Representation, i int) (*url.URL, error) {
   replace(&m.S, "$RepresentationID$", r.Id)
   if m.time() {
      replace(&m.S, "$Time$", fmt.Sprint(i))
   } else {
      replace(&m.S, "$Number$", fmt.Sprint(i))
      replace(&m.S, "$Number%02d$", fmt.Sprintf("%02d", i))
      replace(&m.S, "$Number%03d$", fmt.Sprintf("%03d", i))
      replace(&m.S, "$Number%04d$", fmt.Sprintf("%04d", i))
      replace(&m.S, "$Number%05d$", fmt.Sprintf("%05d", i))
      replace(&m.S, "$Number%06d$", fmt.Sprintf("%06d", i))
      replace(&m.S, "$Number%07d$", fmt.Sprintf("%07d", i))
      replace(&m.S, "$Number%08d$", fmt.Sprintf("%08d", i))
      replace(&m.S, "$Number%09d$", fmt.Sprintf("%09d", i))
   }
   u, err := url.Parse(m.S)
   if err != nil {
      return nil, err
   }
   if r.BaseUrl != nil {
      u = r.BaseUrl.Url.ResolveReference(u)
   }
   return u, nil
}

func (m Media) time() bool {
   return strings.Contains(m.S, "$Time$")
}

type Media struct {
   S string
}

func (m *Media) UnmarshalText(data []byte) error {
   m.S = string(data)
   return nil
}

func (s *SegmentTemplate) set() {
   // dashif.org/Guidelines-TimingModel#addressing-simple
   if s.StartNumber == nil {
      value := 1
      s.StartNumber = &value
   }
   // dashif.org/Guidelines-TimingModel#timing-sampletimeline
   if s.Timescale == nil {
      var value uint64 = 1
      s.Timescale = &value
   }
}

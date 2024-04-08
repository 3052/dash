package dash

import (
   "encoding/xml"
   "math"
   "strconv"
   "strings"
   "time"
)

type AdaptationSet struct {
   Codecs *string `xml:"codecs,attr"`
   Lang *string `xml:"lang,attr"`
   MimeType *string `xml:"mimeType,attr"`
   Representation []*Representation
   Role *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
   period *Period
}

type MPD struct {
   BaseURL string `xml:"BaseURL"`
   MediaPresentationDuration string `xml:"mediaPresentationDuration,attr"`
   Period []*Period
}

func (m *MPD) Unmarshal(data []byte) error {
   err := xml.Unmarshal(data, m)
   if err != nil {
      return err
   }
   for _, period := range m.Period {
      period.mpd = m
      for _, adapt := range period.AdaptationSet {
         adapt.period = period
         for _, represent := range adapt.Representation {
            represent.adaptation_set = adapt
         }
      }
   }
   return nil
}

type Period struct {
   AdaptationSet []*AdaptationSet
   Duration *string `xml:"duration,attr"`
   mpd *MPD
}

func (p Period) Seconds() (float64, error) {
   s := strings.TrimPrefix(p.get_duration(), "PT")
   d, err := time.ParseDuration(strings.ToLower(s))
   if err != nil {
      return 0, err
   }
   return d.Seconds(), nil
}

func (p Period) get_duration() string {
   if v := p.Duration; v != nil {
      return *v
   }
   return p.mpd.MediaPresentationDuration
}

type Range string

// range-end can exceed 32 bits, so we must use 64 bit
func (r Range) End() (uint64, error) {
   _, end, _ := strings.Cut(string(r), "-")
   return strconv.ParseUint(end, 10, 64)
}

// range-start can exceed 32 bits, so we must use 64 bit
func (r Range) Start() (uint64, error) {
   start, _, _ := strings.Cut(string(r), "-")
   return strconv.ParseUint(start, 10, 64)
}

type SegmentTemplate struct {
   Duration *float64 `xml:"duration,attr"`
   Initialization *string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   SegmentTimeline *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R *int `xml:"r,attr"` // repeat
      }
   }
   StartNumber *int `xml:"startNumber,attr"`
   Timescale *float64 `xml:"timescale,attr"`
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#timing-sampletimeline
func (s SegmentTemplate) get_timescale() float64 {
   if v := s.Timescale; v != nil {
      return *v
   }
   return 1
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#addressing-simple-to-explicit
func (s SegmentTemplate) segment_count(seconds float64) float64 {
   seconds /= *s.Duration / s.get_timescale()
   return math.Ceil(seconds)
}

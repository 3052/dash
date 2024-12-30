package dash

import (
   "iter"
   "net/url"
   "strconv"
   "strings"
)

type Initialization func(string) string

func (i *Initialization) UnmarshalText(data []byte) error {
   *i = func(id string) string {
      return strings.Replace(string(data), "$RepresentationID$", id, 1)
   }
   return nil
}

type SegmentTemplate struct {
   Initialization Initialization `xml:"initialization,attr"`
   Duration               int    `xml:"duration,attr"`
   Media                  string `xml:"media,attr"`
   PresentationTimeOffset int    `xml:"presentationTimeOffset,attr"`
   SegmentTimeline        *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R int `xml:"r,attr"` // repeat
      }
   }
   StartNumber *int `xml:"startNumber,attr"`
   Timescale   int  `xml:"timescale,attr"`
}

type Representation struct {
   Bandwidth       int64   `xml:"bandwidth,attr"`
   Codecs          *string `xml:"codecs,attr"`
   Height          *int64  `xml:"height,attr"`
   Id              string  `xml:"id,attr"`
   MimeType        *string `xml:"mimeType,attr"`
   SegmentTemplate *SegmentTemplate
   Width           *int64 `xml:"width,attr"`
   adaptation_set  *AdaptationSet
}

func (r *Representation) set() {
   if r.Codecs == nil {
      r.Codecs = r.adaptation_set.Codecs
   }
   if r.Height == nil {
      r.Height = r.adaptation_set.Height
   }
   if r.MimeType == nil {
      r.MimeType = r.adaptation_set.MimeType
   }
   if r.SegmentTemplate == nil {
      r.SegmentTemplate = r.adaptation_set.SegmentTemplate
   }
   if r.Width == nil {
      r.Width = r.adaptation_set.Width
   }
}

func (m Mpd) representation() iter.Seq[Representation] {
   return func(yield func(Representation) bool) {
      for _, p := range m.Period {
         p.mpd = &m
         for _, adapt := range p.AdaptationSet {
            adapt.period = &p
            for _, represent := range adapt.Representation {
               represent.adaptation_set = &adapt
               represent.set()
               if !yield(represent) {
                  return
               }
            }
         }
      }
   }
}

func (r *Representation) String() string {
   var b []byte
   if r.Width != nil {
      b = append(b, "width = "...)
      b = strconv.AppendInt(b, *r.Width, 10)
   }
   if r.Height != nil {
      if b != nil {
         b = append(b, '\n')
      }
      b = append(b, "height = "...)
      b = strconv.AppendInt(b, *r.Height, 10)
   }
   if b != nil {
      b = append(b, '\n')
   }
   b = append(b, "bandwidth = "...)
   b = strconv.AppendInt(b, r.Bandwidth, 10)
   if r.Codecs != nil {
      b = append(b, "\ncodecs = "...)
      b = append(b, *r.Codecs...)
   }
   b = append(b, "\nmimeType = "...)
   b = append(b, *r.MimeType...)
   if role := r.adaptation_set.Role; role != nil {
      b = append(b, "\nrole = "...)
      b = append(b, role.Value...)
   }
   if lang := r.adaptation_set.Lang; lang != "" {
      b = append(b, "\nlang = "...)
      b = append(b, lang...)
   }
   if id := r.adaptation_set.period.Id; id != "" {
      b = append(b, "\nperiod = "...)
      b = append(b, id...)
   }
   b = append(b, "\nid = "...)
   b = append(b, r.Id...)
   return string(b)
}

type Mpd struct {
   Period []Period
}

type Url struct {
   Url url.URL
}

type AdaptationSet struct {
   Codecs         *string `xml:"codecs,attr"`
   Height         *int64  `xml:"height,attr"`
   Lang           string  `xml:"lang,attr"`
   MimeType       *string `xml:"mimeType,attr"`
   Representation []Representation
   Role           *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
   Width           *int64 `xml:"width,attr"`
   period          *Period
}

func (b *Url) UnmarshalText(data []byte) error {
   return b.Url.UnmarshalBinary(data)
}

type Period struct {
   AdaptationSet []AdaptationSet
   BaseUrl       *Url   `xml:"BaseURL"`
   Id            string `xml:"id,attr"`
   mpd           *Mpd
}

func (r *Representation) seq() iter.Seq[Representation] {
   return func(yield func(Representation) bool) {
      for r2 := range r.adaptation_set.period.mpd.representation() {
         if r2.Id == r.Id {
            if !yield(r2) {
               return
            }
         }
      }
   }
}

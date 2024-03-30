package dash

import (
   "encoding/xml"
   "fmt"
   "strconv"
   "strings"
   "time"
)

func Unmarshal(b []byte) ([]Representation, error) {
   var media mpd
   err := xml.Unmarshal(b, &media)
   if err != nil {
      return nil, err
   }
   var rs []Representation
   for _, p := range media.Period {
      p.mpd = &media
      for _, a := range p.AdaptationSet {
         a.period = &p
         for _, r := range a.Representation {
            r.adaptation_set = &a
            rs = append(rs, r)
         }
      }
   }
   return rs, nil
}

func (r Representation) Ext() (string, bool) {
   switch r.GetMimeType() {
   case "audio/mp4":
      return ".m4a", true
   case "video/mp4":
      return ".m4v", true
   }
   return "", false
}

func (r Representation) GetCodecs() (string, bool) {
   if v := r.Codecs; v != nil {
      return *v, true
   }
   if v := r.adaptation_set.Codecs; v != nil {
      return *v, true
   }
   return "", false
}

func (r Representation) GetMimeType() string {
   if v := r.MimeType; v != nil {
      return *v
   }
   return *r.adaptation_set.MimeType
}

func (r Representation) GetSegmentTemplate() (*SegmentTemplate, bool) {
   if v := r.SegmentTemplate; v != nil {
      return v, true
   }
   if v := r.adaptation_set.SegmentTemplate; v != nil {
      return v, true
   }
   return nil, false
}

func (r Representation) String() string {
   var b []byte
   if v := r.Width; v != nil {
      b = append(b, "width = "...)
      b = strconv.AppendInt(b, *v, 10)
   }
   if v := r.Height; v != nil {
      if b != nil {
         b = append(b, '\n')
      }
      b = append(b, "height = "...)
      b = strconv.AppendInt(b, *v, 10)
   }
   if b != nil {
      b = append(b, '\n')
   }
   b = append(b, "bandwidth = "...)
   b = strconv.AppendInt(b, r.Bandwidth, 10)
   if v, ok := r.GetCodecs(); ok {
      b = append(b, "\ncodecs = "...)
      b = append(b, v...)
   }
   b = append(b, "\ntype = "...)
   b = append(b, r.GetMimeType()...)
   if v := r.adaptation_set.Role; v != nil {
      b = append(b, "\nrole = "...)
      b = append(b, v.Value...)
   }
   if v := r.adaptation_set.Lang; v != nil {
      b = append(b, "\nlang = "...)
      b = append(b, *v...)
   }
   b = append(b, "\nid = "...)
   b = append(b, r.ID...)
   return string(b)
}

type Range string

// range-start and range-end can both exceed 32 bits, so we must use 64 bit
func (r Range) Scan() (uint64, uint64, error) {
   var start, end uint64
   _, err := fmt.Sscanf(string(r), "%v-%v", &start, &end)
   if err != nil {
      return 0, 0, err
   }
   return start, end, nil
}

type Representation struct {
   Bandwidth int64 `xml:"bandwidth,attr"`
   BaseURL *string
   Codecs *string `xml:"codecs,attr"`
   Height *int64 `xml:"height,attr"`
   ID string `xml:"id,attr"`
   MimeType *string `xml:"mimeType,attr"`
   SegmentBase *struct {
      IndexRange Range `xml:"indexRange,attr"`
      Initialization struct {
         Range Range `xml:"range,attr"`
      }
   }
   SegmentTemplate *SegmentTemplate
   Width *int64 `xml:"width,attr"`
   adaptation_set *AdaptationSet
}

func (r Representation) GetAdaptationSet() *AdaptationSet {
   return r.adaptation_set
}

type AdaptationSet struct {
   Codecs *string `xml:"codecs,attr"`
   Lang *string `xml:"lang,attr"`
   MimeType *string `xml:"mimeType,attr"`
   Representation []Representation
   Role *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
   period *Period
}

func (a AdaptationSet) GetPeriod() *Period {
   return a.period
}

type mpd struct {
   MediaPresentationDuration string `xml:"mediaPresentationDuration,attr"`
   Period []Period
}

type Period struct {
   AdaptationSet []AdaptationSet
   Duration *string `xml:"duration,attr"`
   mpd *mpd
}

func (p Period) GetDuration() string {
   if v := p.Duration; v != nil {
      return *v
   }
   return p.mpd.MediaPresentationDuration
}

func (p Period) Seconds() (float64, error) {
   s := strings.TrimPrefix(p.GetDuration(), "PT")
   d, err := time.ParseDuration(strings.ToLower(s))
   if err != nil {
      return 0, err
   }
   return d.Seconds(), nil
}

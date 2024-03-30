package dash

import (
   "encoding/xml"
   "fmt"
   "strconv"
   "strings"
   "time"
)

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#addressing-simple-to-explicit
type mpd struct {
   MediaPresentationDuration string `xml:"mediaPresentationDuration,attr"`
   Period []period
}

type Range string

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

// range-start and range-end can both exceed 32 bits, so we must use 64 bit
func (r Range) Scan() (uint64, uint64, error) {
   var start, end uint64
   _, err := fmt.Sscanf(string(r), "%v-%v", &start, &end)
   if err != nil {
      return 0, 0, err
   }
   return start, end, nil
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

func (m mpd) Seconds() (float64, error) {
   s := strings.TrimPrefix(m.MediaPresentationDuration, "PT")
   duration, err := time.ParseDuration(strings.ToLower(s))
   if err != nil {
      return 0, err
   }
   return duration.Seconds(), nil
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

type period struct {
   AdaptationSet []adaptation_set
   mpd *mpd
}

type adaptation_set struct {
   Codecs *string `xml:"codecs,attr"`
   Lang *string `xml:"lang,attr"`
   MimeType *string `xml:"mimeType,attr"`
   Representation []Representation
   Role *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
   period *period
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
   adaptation_set *adaptation_set
}

func (s SegmentTemplate) GetInitialization(id string) (string, bool) {
   if v := s.Initialization; v != nil {
      return strings.Replace(*v, "$RepresentationID$", id, 1), true
   }
   return "", false
}

type SegmentTemplate struct {
   Initialization *string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   SegmentTimeline *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R *int `xml:"r,attr"` // repeat
      }
   }
   StartNumber *int `xml:"startNumber,attr"`
}

func (s SegmentTemplate) GetMedia(id string) []string {
   timeline := s.SegmentTimeline
   if timeline == nil {
      return nil
   }
   s.Media = strings.Replace(s.Media, "$RepresentationID$", id, 1)
   var number int
   if s.StartNumber != nil {
      number = *s.StartNumber
   }
   var media []string
   for _, segment := range timeline.S {
      var repeat int
      if segment.R != nil {
         repeat = *segment.R
      }
      for repeat >= 0 {
         var medium string
         replace := strconv.Itoa(number)
         if s.StartNumber != nil {
            medium = strings.Replace(s.Media, "$Number$", replace, 1)
            number++
         } else {
            medium = strings.Replace(s.Media, "$Time$", replace, 1)
            number += segment.D
         }
         media = append(media, medium)
         repeat--
      }
   }
   return media
}

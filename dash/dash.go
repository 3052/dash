package dash

import (
   "encoding/xml"
   "fmt"
   "strconv"
   "strings"
   "time"
)

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

func (r Representation) GetCodecs() (string, bool) {
   if v := r.Codecs; v != "" {
      return v, true
   }
   if v := r.adaptation_set.Codecs; v != "" {
      return v, true
   }
   return "", false
}

func (r Representation) GetMimeType() string {
   if v := r.MimeType; v != "" {
      return v
   }
   return r.adaptation_set.MimeType
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

func (r Representation) Initialization() (string, bool) {
   if v, ok := r.GetSegmentTemplate(); ok {
      if v := v.Initialization; v != "" {
         return strings.Replace(v, "$RepresentationID$", r.ID, 1), true
      }
   }
   return "", false
}

// we need the length for progress meter, so cannot use a channel
func (r Representation) Media() []string {
   st, ok := r.GetSegmentTemplate()
   if !ok {
      return nil
   }
   st.Media = strings.Replace(st.Media, "$RepresentationID$", r.ID, 1)
   var number int
   if st.StartNumber != nil {
      number = *st.StartNumber
   }
   var media []string
   for _, segment := range st.SegmentTimeline.S {
      for segment.R >= 0 {
         var medium string
         if st.StartNumber != nil {
            medium = st.replace("$Number$", number)
            number++
         } else {
            medium = st.replace("$Time$", number)
            number += segment.D
         }
         media = append(media, medium)
         segment.R--
      }
   }
   return media
}

func (r Representation) String() string {
   var b []byte
   if v := r.Width; v >= 1 {
      b = append(b, "width = "...)
      b = strconv.AppendInt(b, v, 10)
   }
   if v := r.Height; v >= 1 {
      if b != nil {
         b = append(b, '\n')
      }
      b = append(b, "height = "...)
      b = strconv.AppendInt(b, v, 10)
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
   if v := r.adaptation_set.Lang; v != "" {
      b = append(b, "\nlang = "...)
      b = append(b, v...)
   }
   b = append(b, "\nid = "...)
   b = append(b, r.ID...)
   return string(b)
}

type SegmentTemplate struct {
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   SegmentTimeline *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R int `xml:"r,attr"` // repeat
      }
   }
   StartNumber *int `xml:"startNumber,attr"`
}

func (s SegmentTemplate) replace(old string, number int) string {
   return strings.Replace(s.Media, old, strconv.Itoa(number), 1)
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#addressing-simple-to-explicit
type mpd struct {
   MediaPresentationDuration string `xml:"mediaPresentationDuration,attr"`
   Period []period
}

func (m mpd) seconds() (float64, error) {
   s := strings.TrimPrefix(m.MediaPresentationDuration, "PT")
   duration, err := time.ParseDuration(strings.ToLower(s))
   if err != nil {
      return 0, err
   }
   return duration.Seconds(), nil
}

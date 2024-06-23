package dash

import (
   "encoding/xml"
   "fmt"
   "math"
   "strconv"
   "strings"
   "time"
)

type Period struct {
   AdaptationSet []*AdaptationSet
   Duration      string `xml:"duration,attr"`
   mpd           *Mpd
}

type ContentProtection struct {
   Pssh        string `xml:"pssh"`
   SchemeIdUri string `xml:"schemeIdUri,attr"`
}

// SegmentIndexBox uses:
// unsigned int(32) subsegment_duration;
// but range values can exceed 32 bits
type Range struct {
   Start uint64
   End   uint64
}

type Mpd struct {
   Period                    []*Period
   MediaPresentationDuration string `xml:"mediaPresentationDuration,attr"`
   BaseUrl                   string `xml:"BaseURL"`
}

type SegmentBase struct {
   Initialization struct {
      Range Range `xml:"range,attr"`
   }
   IndexRange Range `xml:"indexRange,attr"`
}

type AdaptationSet struct {
   ContentProtection []ContentProtection
   Representation    []*Representation
   period            *Period
   Codecs            string `xml:"codecs,attr"`
   Height            int64  `xml:"height,attr"`
   Lang              string `xml:"lang,attr"`
   MimeType          string `xml:"mimeType,attr"`
   Width             int64  `xml:"width,attr"`
   Role              *struct {
      Value string `xml:"value,attr"`
   }
   SegmentTemplate *SegmentTemplate
}

type Representation struct {
   Bandwidth         int64  `xml:"bandwidth,attr"`
   BaseUrl           string `xml:"BaseURL"`
   ContentProtection []ContentProtection
   Height            int64  `xml:"height,attr"`
   Id                string `xml:"id,attr"`
   MimeType          string `xml:"mimeType,attr"`
   Width             int64  `xml:"width,attr"`
   adaptation_set    *AdaptationSet
   Codecs            string `xml:"codecs,attr"`
   SegmentBase       *SegmentBase
   SegmentTemplate   *SegmentTemplate
}

type SegmentTemplate struct {
   Duration float64 `xml:"duration,attr"`
   Initialization string `xml:"initialization,attr"`
   Media string `xml:"media,attr"`
   StartNumber int `xml:"startNumber,attr"`
   PresentationTimeOffset int `xml:"presentationTimeOffset,attr"`
   Timescale float64 `xml:"timescale,attr"`
   SegmentTimeline *struct {
      S []struct {
         D int `xml:"d,attr"` // duration
         R int `xml:"r,attr"` // repeat
      }
   }
}

func (s SegmentTemplate) get_initialization() (string, bool) {
   if v := s.Initialization; v != "" {
      return v, true
   }
   return "", false
}

/////////

func (r Representation) get_segment_template() (*SegmentTemplate, bool) {
   if v := r.SegmentTemplate; v != nil {
      return v, true
   }
   if v := r.adaptation_set.SegmentTemplate; v != nil {
      return v, true
   }
   return nil, false
}

func (s SegmentTemplate) start() int {
   if v := s.PresentationTimeOffset; v >= 1 {
      return v
   }
   return s.StartNumber
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#timing-sampletimeline
func (s SegmentTemplate) get_timescale() float64 {
   if v := s.Timescale; v >= 1 {
      return v
   }
   return 1
}

// dashif-documents.azurewebsites.net/Guidelines-TimingModel/master/Guidelines-TimingModel.html#addressing-simple-to-explicit
func (s SegmentTemplate) segment_count(seconds float64) float64 {
   seconds /= s.Duration / s.get_timescale()
   return math.Ceil(seconds)
}

func (s SegmentTemplate) number(value int) string {
   f := strings.Replace(s.Media, "$Number$", "%d", 1)
   f = strings.Replace(f, "$Number%02d$", "%02d", 1)
   f = strings.Replace(f, "$Number%03d$", "%03d", 1)
   f = strings.Replace(f, "$Number%04d$", "%04d", 1)
   f = strings.Replace(f, "$Number%05d$", "%05d", 1)
   f = strings.Replace(f, "$Number%06d$", "%06d", 1)
   f = strings.Replace(f, "$Number%07d$", "%07d", 1)
   f = strings.Replace(f, "$Number%08d$", "%08d", 1)
   f = strings.Replace(f, "$Number%09d$", "%09d", 1)
   return fmt.Sprintf(f, value)
}

func (s SegmentTemplate) time(value int) string {
   f := strings.Replace(s.Media, "$Time$", "%d", 1)
   return fmt.Sprintf(f, value)
}

func (s SegmentTemplate) GetMedia(r *Representation) ([]string, error) {
   s.Media = r.id(s.Media)
   var media []string
   number := s.start()
   if s.SegmentTimeline != nil {
      for _, segment := range s.SegmentTimeline.S {
         var repeat int
         if segment.R >= 1 {
            repeat = segment.R
         }
         for range 1 + repeat {
            var medium string
            if strings.Contains(s.Media, "$Time$") {
               medium = s.time(number)
               number += segment.D
            } else {
               medium = s.number(number)
               number++
            }
            media = append(media, medium)
         }
      }
   } else {
      seconds, err := r.adaptation_set.period.Seconds()
      if err != nil {
         return nil, err
      }
      for range int(s.segment_count(seconds)) {
         media = append(media, s.number(number))
         number++
      }
   }
   return media, nil
}

func (r Representation) id(value string) string {
   return strings.Replace(value, "$RepresentationID$", r.Id, 1)
}

func (r Representation) get_codecs() (string, bool) {
   if v := r.Codecs; v != "" {
      return v, true
   }
   if v := r.adaptation_set.Codecs; v != "" {
      return v, true
   }
   return "", false
}

func (r Representation) get_height() (int64, bool) {
   if v := r.Height; v >= 1 {
      return v, true
   }
   if v := r.adaptation_set.Height; v >= 1 {
      return v, true
   }
   return 0, false
}

func (r Representation) get_width() (int64, bool) {
   if v := r.Width; v >= 1 {
      return v, true
   }
   if v := r.adaptation_set.Width; v >= 1 {
      return v, true
   }
   return 0, false
}

func (r Representation) get_mime_type() string {
   if v := r.MimeType; v != "" {
      return v
   }
   return r.adaptation_set.MimeType
}

func (p Period) get_duration() string {
   if v := p.Duration; v != "" {
      return v
   }
   return p.mpd.MediaPresentationDuration
}

func (r Representation) Ext() (string, bool) {
   switch r.get_mime_type() {
   case "audio/mp4":
      return ".m4a", true
   case "video/mp4":
      return ".m4v", true
   }
   return "", false
}

func (r Range) MarshalText() ([]byte, error) {
   b := strconv.AppendUint(nil, r.Start, 10)
   b = append(b, '-')
   return strconv.AppendUint(b, r.End, 10), nil
}

func (r *Range) UnmarshalText(text []byte) error {
   // the current testdata always has `-`, so lets assume for now
   start, end, _ := strings.Cut(string(text), "-")
   var err error
   r.Start, err = strconv.ParseUint(start, 10, 64)
   if err != nil {
      return err
   }
   r.End, err = strconv.ParseUint(end, 10, 64)
   if err != nil {
      return err
   }
   return nil
}

func (r Representation) protection() []ContentProtection {
   if v := r.ContentProtection; v != nil {
      return v
   }
   return r.adaptation_set.ContentProtection
}

func (r Representation) GetAdaptationSet() *AdaptationSet {
   return r.adaptation_set
}

func (m *Mpd) Unmarshal(data []byte) error {
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

func (a AdaptationSet) GetPeriod() *Period {
   return a.period
}

// filter out ads, for example:
// hulu.com/watch/5add1b6c-04f2-4038-a925-35db3007d662
func (p Period) Seconds() (float64, error) {
   s := strings.TrimPrefix(p.get_duration(), "PT")
   duration, err := time.ParseDuration(strings.ToLower(s))
   if err != nil {
      return 0, err
   }
   return duration.Seconds(), nil
}

func (p Period) GetMpd() *Mpd {
   return p.mpd
}

func (r Representation) Widevine() (string, bool) {
   for _, p := range r.protection() {
      if p.SchemeIdUri == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         if p.Pssh != "" {
            return p.Pssh, true
         }
      }
   }
   return "", false
}

func (r Representation) String() string {
   var b []byte
   if v, ok := r.get_width(); ok {
      b = append(b, "width = "...)
      b = strconv.AppendInt(b, v, 10)
   }
   if v, ok := r.get_height(); ok {
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
   if v, ok := r.get_codecs(); ok {
      b = append(b, "\ncodecs = "...)
      b = append(b, v...)
   }
   b = append(b, "\ntype = "...)
   b = append(b, r.get_mime_type()...)
   if v := r.adaptation_set.Role; v != nil {
      b = append(b, "\nrole = "...)
      b = append(b, v.Value...)
   }
   if v := r.adaptation_set.Lang; v != "" {
      b = append(b, "\nlang = "...)
      b = append(b, v...)
   }
   b = append(b, "\nid = "...)
   b = append(b, r.Id...)
   return string(b)
}

func (r Representation) Initialization() (string, bool) {
   if v, ok := r.get_segment_template(); ok {
      if v, ok := v.get_initialization(); ok {
         return r.id(v), true
      }
   }
   return "", false
}

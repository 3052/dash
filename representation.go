package dash

import (
   "encoding/xml"
   "net/url"
   "strconv"
   "strings"
)

func (r Representation) Media() []string {
   template, ok := r.get_segment_template()
   if !ok {
      return nil
   }
   number := template.start()
   template.Media = r.id(template.Media)
   var media []string
   if template.SegmentTimeline != nil {
      for _, segment := range template.SegmentTimeline.S {
         for range 1 + segment.R {
            var medium string
            if strings.Contains(template.Media, "$Time$") {
               medium = template.time(number)
               number += segment.D
            } else {
               medium = template.number(number)
               number++
            }
            media = append(media, medium)
         }
      }
   } else {
      seconds := r.adaptation_set.period.get_duration().Duration.Seconds()
      for range template.segment_count(seconds) {
         media = append(media, template.number(number))
         number++
      }
   }
   return media
}

func (r Representation) String() string {
   var b []byte
   if v := r.get_width(); v >= 1 {
      b = append(b, "width = "...)
      b = strconv.AppendUint(b, v, 10)
   }
   if v := r.get_height(); v >= 1 {
      if b != nil {
         b = append(b, '\n')
      }
      b = append(b, "height = "...)
      b = strconv.AppendUint(b, v, 10)
   }
   if b != nil {
      b = append(b, '\n')
   }
   b = append(b, "bandwidth = "...)
   b = strconv.AppendUint(b, r.Bandwidth, 10)
   if v := r.get_codecs(); v != "" {
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
   if v := r.adaptation_set.period.Id; v != "" {
      b = append(b, "\nperiod = "...)
      b = append(b, v...)
   }
   b = append(b, "\nid = "...)
   b = append(b, r.Id...)
   return string(b)
}

func (r Representation) GetAdaptationSet() *AdaptationSet {
   return r.adaptation_set
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

func (r Representation) id(value string) string {
   return strings.Replace(value, "$RepresentationID$", r.Id, 1)
}

func (r Representation) get_mime_type() string {
   if r.MimeType != "" {
      return r.MimeType
   }
   return r.adaptation_set.MimeType
}

func (r Representation) get_content_protection() []ContentProtection {
   if len(r.ContentProtection) >= 1 {
      return r.ContentProtection
   }
   return r.adaptation_set.ContentProtection
}

func (r Representation) get_width() uint64 {
   if r.Width >= 1 {
      return r.Width
   }
   return r.adaptation_set.Width
}

func (r Representation) get_height() uint64 {
   if r.Height >= 1 {
      return r.Height
   }
   return r.adaptation_set.Height
}

func (r Representation) get_codecs() string {
   if r.Codecs != "" {
      return r.Codecs
   }
   return r.adaptation_set.Codecs
}

func (r Representation) get_segment_template() (*SegmentTemplate, bool) {
   if r.SegmentTemplate != nil {
      return r.SegmentTemplate, true
   }
   if r.adaptation_set.SegmentTemplate != nil {
      return r.adaptation_set.SegmentTemplate, true
   }
   return nil, false
}

func (r Representation) Widevine() (Pssh, bool) {
   for _, v := range r.get_content_protection() {
      if v.SchemeIdUri == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         if len(v.Pssh) >= 1 {
            return v.Pssh, true
         }
      }
   }
   return nil, false
}

func (r Representation) Initialization() (string, bool) {
   if v, ok := r.get_segment_template(); ok {
      if v := v.Initialization; v != "" {
         return r.id(v), true
      }
   }
   return "", false
}
func (r Representation) GetBaseUrl() (*BaseUrl, bool) {
   var u *url.URL
   if v := r.adaptation_set.period.mpd.BaseUrl; v != nil {
      u = new(url.URL)
      *u = *v.Url
   }
   if v := r.adaptation_set.period.BaseUrl; v != nil {
      if u == nil {
         u = new(url.URL)
      }
      u = u.ResolveReference(v.Url)
   }
   if v := r.BaseUrl; v != nil {
      if u == nil {
         u = new(url.URL)
      }
      u = u.ResolveReference(v.Url)
   }
   if u != nil {
      return &BaseUrl{u}, true
   }
   return nil, false
}

func Unmarshal(text []byte, base *url.URL) ([]Representation, error) {
   var media Mpd
   err := xml.Unmarshal(text, &media)
   if err != nil {
      return nil, err
   }
   if media.BaseUrl == nil {
      if base != nil {
         media.BaseUrl = &BaseUrl{base}
      }
   }
   var reps []Representation
   for _, per := range media.Period {
      per.mpd = &media
      for _, ada := range per.AdaptationSet {
         ada.period = &per
         for _, rep := range ada.Representation {
            rep.adaptation_set = &ada
            reps = append(reps, rep)
         }
      }
   }
   return reps, nil
}

type Representation struct {
   Bandwidth         uint64 `xml:"bandwidth,attr"`
   BaseUrl           *BaseUrl   `xml:"BaseURL"`
   Codecs            string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height            uint64 `xml:"height,attr"`
   Id                string `xml:"id,attr"`
   MimeType          string `xml:"mimeType,attr"`
   SegmentBase       *SegmentBase
   SegmentTemplate   *SegmentTemplate
   Width             uint64 `xml:"width,attr"`
   adaptation_set    *AdaptationSet
}

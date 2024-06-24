package dash

import (
   "strconv"
   "strings"
)

type Representation struct {
   Bandwidth         uint64  `xml:"bandwidth,attr"`
   BaseUrl           string `xml:"BaseURL"`
   Codecs            string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height            uint64  `xml:"height,attr"`
   Id                string `xml:"id,attr"`
   MimeType          string `xml:"mimeType,attr"`
   SegmentBase       *SegmentBase
   SegmentTemplate   *SegmentTemplate
   Width             uint64  `xml:"width,attr"`
   adaptation_set    *AdaptationSet
}

func (r Representation) GetAdaptationSet() *AdaptationSet {
   return r.adaptation_set
}

func (r Representation) get_segment_template() (*SegmentTemplate, bool) {
   return option(r.SegmentTemplate, r.adaptation_set.SegmentTemplate)
}

func (r Representation) GetMedia() []string {
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
      seconds := r.adaptation_set.period.get_duration().D.Seconds()
      for range template.segment_count(seconds) {
         media = append(media, template.number(number))
         number++
      }
   }
   return media
}

func (r Representation) get_mime_type() string {
   if v := r.MimeType; v != "" {
      return v
   }
   return r.adaptation_set.MimeType
}

func (r Representation) String() string {
   var b []byte
   if v, ok := r.get_width(); ok {
      b = append(b, "width = "...)
      b = strconv.AppendUint(b, v, 10)
   }
   if v, ok := r.get_height(); ok {
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
   if v, ok := r.get_codecs(); ok {
      b = append(b, "\ncodecs = "...)
      b = append(b, v...)
   }
   b = append(b, "\ntype = "...)
   b = append(b, r.get_mime_type()...)
   if v, ok := r.adaptation_set.get_role(); ok {
      b = append(b, "\nrole = "...)
      b = append(b, v.Value...)
   }
   if v, ok := r.adaptation_set.get_lang(); ok {
      b = append(b, "\nlang = "...)
      b = append(b, v...)
   }
   b = append(b, "\nid = "...)
   b = append(b, r.Id...)
   return string(b)
}

func (r Representation) id(value string) string {
   return strings.Replace(value, "$RepresentationID$", r.Id, 1)
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

func (r Representation) Widevine() (string, bool) {
   for _, p := range r.get_content_protection() {
      if p.SchemeIdUri == "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed" {
         return p.get_pssh()
      }
   }
   return "", false
}

func (r Representation) get_width() (uint64, bool) {
   return option(r.Width, r.adaptation_set.Width)
}

func (r Representation) get_height() (uint64, bool) {
   return option(r.Height, r.adaptation_set.Height)
}

func (r Representation) get_codecs() (string, bool) {
   return option(r.Codecs, r.adaptation_set.Codecs)
}

func (r Representation) Initialization() (string, bool) {
   if v, ok := r.get_segment_template(); ok {
      if v, ok := v.get_initialization(); ok {
         return r.id(v), true
      }
   }
   return "", false
}

func (r Representation) get_content_protection() []ContentProtection {
   if v := r.ContentProtection; len(v) >= 1 {
      return v
   }
   return r.adaptation_set.ContentProtection
}

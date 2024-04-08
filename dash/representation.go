package dash

import (
   "errors"
   "strconv"
   "strings"
)

func (r Representation) GetInitialization() (string, bool) {
   if v, ok := r.get_segment_template(); ok {
      if v := v.Initialization; v != nil {
         return strings.Replace(*v, "$RepresentationID$", r.ID, 1), true
      }
   }
   return "", false
}

func (r Representation) GetMedia() ([]string, error) {
   st, ok := r.get_segment_template()
   if !ok {
      return nil, errors.New("get_segment_template")
   }
   st.Media = strings.Replace(st.Media, "$RepresentationID$", r.ID, 1)
   var (
      media []string
      number int
   )
   if st.StartNumber != nil {
      number = *st.StartNumber
   }
   if st.SegmentTimeline != nil {
      for _, segment := range st.SegmentTimeline.S {
         var repeat int
         if segment.R != nil {
            repeat = *segment.R
         }
         for range 1 + repeat {
            var medium string
            replace := strconv.Itoa(number)
            if st.StartNumber != nil {
               medium = strings.Replace(st.Media, "$Number$", replace, 1)
               number++
            } else {
               medium = strings.Replace(st.Media, "$Time$", replace, 1)
               number += segment.D
            }
            media = append(media, medium)
         }
      }
   } else {
      seconds, err := r.adaptation_set.period.Seconds()
      if err != nil {
         return nil, err
      }
      for range int(st.segment_count(seconds)) {
         replace := strconv.Itoa(number)
         medium := strings.Replace(st.Media, "$Number$", replace, 1)
         media = append(media, medium)
         number++
      }
   }
   return media, nil
}

func (r Representation) get_segment_template() (*SegmentTemplate, bool) {
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
   adaptation_set *AdaptationSet
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
   if v := r.adaptation_set.Lang; v != nil {
      b = append(b, "\nlang = "...)
      b = append(b, *v...)
   }
   b = append(b, "\nid = "...)
   b = append(b, r.ID...)
   return string(b)
}

func (r Representation) get_codecs() (string, bool) {
   if v := r.Codecs; v != nil {
      return *v, true
   }
   if v := r.adaptation_set.Codecs; v != nil {
      return *v, true
   }
   return "", false
}

func (r Representation) get_mime_type() string {
   if v := r.MimeType; v != nil {
      return *v
   }
   return *r.adaptation_set.MimeType
}

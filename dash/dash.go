package dash

import (
   "fmt"
   "strconv"
   "strings"
)

type Range struct {
   Start uint64
   End uint64
}

type RawRange string

func (r RawRange) Scan() (*Range, error) {
   var v Range
   _, err := fmt.Sscanf(string(r), "%v-%v", &v.Start, &v.End)
   if err != nil {
      return nil, err
   }
   return &v, nil
}

// we need the length for progress meter, so cannot use a channel
func (r Representation) Media() []string {
   st, ok := r.segment_template()
   if !ok {
      return nil
   }
   replace := func(s, old string) string {
      s = strings.Replace(s, "$RepresentationID$", r.ID, 1)
      return strings.Replace(s, old, strconv.Itoa(st.StartNumber), 1)
   }
   var media []string
   for _, segment := range st.SegmentTimeline.S {
      for segment.R >= 0 {
         var medium string
         if strings.Contains(st.Media, "$Time$") {
            medium = replace(st.Media, "$Time$")
            st.StartNumber += segment.D
         } else {
            medium = replace(st.Media, "$Number$")
            st.StartNumber++
         }
         media = append(media, medium)
         segment.R--
      }
   }
   return media
}

func (r Representation) Ext() (string, bool) {
   switch r.mime_type() {
   case "audio/mp4":
      return ".m4a", true
   case "video/mp4":
      return ".m4v", true
   }
   return "", false
}

func (r Representation) Initialization() (string, bool) {
   if st, ok := r.segment_template(); ok {
      if i := st.Initialization; i != "" {
         return strings.Replace(i, "$RepresentationID$", r.ID, 1), true
      }
   }
   return "", false
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
   if v, ok := r.codecs(); ok {
      b = append(b, "\ncodecs = "...)
      b = append(b, v...)
   }
   b = append(b, "\ntype = "...)
   b = append(b, r.mime_type()...)
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

func (r Representation) codecs() (string, bool) {
   if v := r.Codecs; v != "" {
      return v, true
   }
   if v := r.adaptation_set.Codecs; v != "" {
      return v, true
   }
   return "", false
}

func (r Representation) mime_type() string {
   if v := r.MimeType; v != "" {
      return v
   }
   return r.adaptation_set.MimeType
}

func (r Representation) segment_template() (*SegmentTemplate, bool) {
   if v := r.SegmentTemplate; v != nil {
      return v, true
   }
   if v := r.adaptation_set.SegmentTemplate; v != nil {
      return v, true
   }
   return nil, false
}

type SegmentTemplate struct {
   Media string `xml:"media,attr"`
   SegmentTimeline struct {
      S []struct {
         // duration
         D int `xml:"d,attr"`
         // repeat. this may not exist
         R int `xml:"r,attr"`
      }
   }
   StartNumber int `xml:"startNumber,attr"`
   // this may not exist
   Initialization string `xml:"initialization,attr"`
}

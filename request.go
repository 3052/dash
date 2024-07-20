package dash

import (
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

// SegmentIndexBox uses:
// unsigned int(32) subsegment_duration;
// but range values can exceed 32 bits
type Range struct {
   Start uint64
   End   uint64
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

type Representation struct {
   Bandwidth         uint64 `xml:"bandwidth,attr"`
   BaseUrl           *BaseUrl   `xml:"BaseURL"`
   Codecs            string `xml:"codecs,attr"`
   ContentProtection []ContentProtection
   Height            uint64 `xml:"height,attr"`
   Id                string `xml:"id,attr"`
   MimeType          string `xml:"mimeType,attr"`
   SegmentBase       *struct {
      Initialization struct {
         Range Range `xml:"range,attr"`
      }
      IndexRange Range `xml:"indexRange,attr"`
   }
   SegmentTemplate   *SegmentTemplate
   Width             uint64 `xml:"width,attr"`
   adaptation_set    *AdaptationSet
}

func (r Representation) Initialization() (*http.Request, error) {
   if v := r.SegmentBase; v != nil {
      var req http.Request
      req.URL = r.get_base_url()
   }
   if v, ok := r.get_segment_template(); ok {
      if v := v.Initialization; v != nil {
         var media strings.Builder
         var data struct {
            Representation struct {
               Id string
            }
         }
         data.Representation.Id = r.Id
         err := v.Template.Execute(&media, data)
         if err != nil {
            return nil, true, err
         }
         address, err := r.get_base_url().Parse(media.String())
         if err != nil {
            return nil, true, err
         }
         return address, true, nil
      }
      return nil, false, nil
   }
   return nil, false, nil
}

func (r Representation) Media(t SegmentTemplate) ([]*url.URL, error) {
   var media []string
   var data struct {
      Number uint
      Representation struct {
         Id string
      }
      Time uint
   }
   data.Number = t.StartNumber
   data.Time = t.PresentationTimeOffset
   data.Representation.Id = r.Id
   if t.SegmentTimeline != nil {
      for _, segment := range t.SegmentTimeline.S {
         for range 1 + segment.R {
            var medium strings.Builder
            err := t.Media.Template.Execute(&medium, data)
            if err != nil {
               return nil, err
            }
            media = append(media, medium.String())
            data.Number++
            data.Time += segment.D
         }
      }
   } else {
      seconds := r.adaptation_set.period.get_duration().Duration.Seconds()
      for range t.segment_count(seconds) {
         var medium strings.Builder
         err := t.Media.Template.Execute(&medium, data)
         if err != nil {
            return nil, err
         }
         media = append(media, medium.String())
         data.Number++
      }
   }
   return media, nil
}

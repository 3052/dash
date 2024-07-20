package dash

import (
   "net/http"
   "net/url"
   "strings"
)

func (r Representation) media_template() ([]http.Request, error) {
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

func (r Representation) media_base() ([]http.Request, error) {
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

func (r Representation) Media() ([]http.Request, error) {
   if _, ok := r.get_segment_template(); ok {
      return r.media_template()
   }
   return r.media_base()
}

func (r Representation) Initialization() (*http.Request, bool, error) {
   var req http.Request
   if v := r.SegmentBase; v != nil {
      req.URL = r.get_base_url()
      req.Header = http.Header{
         "range": {v.Initialization.Range},
      }
      return &req, true, nil
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
         req.URL, err = r.get_base_url().Parse(media.String())
         if err != nil {
            return nil, true, err
         }
         return &req, true, nil
      }
   }
   return nil, false, nil
}

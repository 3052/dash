package dash

import (
   "net/http"
   "strings"
)

func (r Representation) media_base() ([]http.Request, error) {
   var reqs []string
   var data struct {
      Number uint
      Representation struct {
         Id string
      }
      Time uint
   }
   data.Number = st.StartNumber
   data.Time = st.PresentationTimeOffset
   data.Representation.Id = r.Id
   if st.SegmentTimeline != nil {
      for _, segment := range st.SegmentTimeline.S {
         for range 1 + segment.R {
            var media strings.Builder
            err := st.Media.Template.Execute(&media, data)
            if err != nil {
               return nil, err
            }
            reqs = append(reqs, media.String())
            data.Number++
            data.Time += segment.D
         }
      }
   } else {
      seconds := r.adaptation_set.period.get_duration().Duration.Seconds()
      for range st.segment_count(seconds) {
         var media strings.Builder
         err := st.Media.Template.Execute(&media, data)
         if err != nil {
            return nil, err
         }
         reqs = append(reqs, media.String())
         data.Number++
      }
   }
   return reqs, nil
}

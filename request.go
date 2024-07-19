package dash

import (
   "net/http"
   "strings"
)

func (r Representation) Initialization() (*http.Request, error) {
   // SegmentBase
   // SegmentTemplate
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

func (r Representation) Media() []string {
   // `template` is a pointer, so if we edit `template.Media` it is permanent
   template, ok := r.get_segment_template()
   if !ok {
      return nil
   }
   id := r.id(template.Media)
   number := template.start()
   var media []string
   if template.SegmentTimeline != nil {
      for _, segment := range template.SegmentTimeline.S {
         for range 1 + segment.R {
            var medium string
            if strings.Contains(template.Media, "$Time$") {
               medium = replace_time(id, number)
               number += segment.D
            } else {
               medium = replace_number(id, number)
               number++
            }
            media = append(media, medium)
         }
      }
   } else {
      seconds := r.adaptation_set.period.get_duration().Duration.Seconds()
      for range template.segment_count(seconds) {
         media = append(media, replace_number(id, number))
         number++
      }
   }
   return media
}

func replace_number(format string, a uint) string {
   format = strings.NewReplacer(
      "$Number$", "%d",
      "$Number%02d$", "%02d",
      "$Number%03d$", "%03d",
      "$Number%04d$", "%04d",
      "$Number%05d$", "%05d",
      "$Number%06d$", "%06d",
      "$Number%07d$", "%07d",
      "$Number%08d$", "%08d",
      "$Number%09d$", "%09d",
   ).Replace(format)
   return fmt.Sprintf(format, a)
}

func replace_time(format string, a uint) string {
   format = strings.Replace(format, "$Time$", "%d", 1)
   return fmt.Sprintf(format, a)
}

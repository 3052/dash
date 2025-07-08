func (s *SegmentTemplate) segments(periodVar *Period) []string {
   var segment []string
   for _, number := range s.numbers(periodVar) {
      segment = append(segment, strconv.Itoa(number))
   }
   return segment
}

// with current data this always uses number addressing
func (s *SegmentTemplate) byEndNumber() []int {
   var numbers []int
   number := *s.StartNumber
   for number <= s.EndNumber {
      numbers = append(numbers, number)
      number++
   }
   return numbers
}

func (*SegmentBase) segments() []string {
   return nil
}

func (s *SegmentList) segments() []string {
   var segments []string
   for _, segment := range s.SegmentUrl {
      segments = append(segments, segment.Media)
   }
   return segments
}

func (r *Representation) Segments(
   adapt *AdaptationSet, periodVar *Period,
) []string {
   if r.SegmentBase != nil {
      return r.SegmentBase.segments()
   }
   if r.SegmentList != nil {
      return r.SegmentList.segments()
   }
   if r.SegmentTemplate != nil {
      return r.SegmentTemplate.segments(periodVar)
   }
   return adapt.SegmentTemplate.segments(periodVar)
}

func (s *SegmentTemplate) numbers(periodVar *Period) []int {
   if s.EndNumber >= 1 {
      return s.byEndNumber()
   }
   if s.SegmentTimeline != nil {
      if strings.Contains(s.Media, "$Time$") {
         return s.byTimelineTime()
      }
      return s.byTimelineNumber()
   }
   return s.byPeriod(periodVar)
}

package dash

import "encoding/xml"

// SegmentTimeline provides an explicit timeline for media segments.
type SegmentTimeline struct {
   XMLName  xml.Name `xml:"SegmentTimeline"`
   Segments []*S     `xml:"S"`
}

// timelineSegment is an internal representation of a single, fully-realized
// segment from a timeline, with its start time and duration.
type timelineSegment struct {
   StartTime uint64
   Duration  uint64
}

// GetSegments expands the S elements, accounting for the repeat attribute 'r',
// to return a flat list of every individual segment's start time and duration.
func (st *SegmentTimeline) GetSegments() []timelineSegment {
   if st == nil {
      return nil
   }
   var expanded []timelineSegment
   var currentTime uint64

   for _, s := range st.Segments {
      // If t is specified, the timeline jumps to this new time.
      if s.T > 0 {
         currentTime = s.T
      }

      // The S element itself represents the first occurrence.
      expanded = append(expanded, timelineSegment{StartTime: currentTime, Duration: s.D})
      currentTime += s.D

      // Handle the repeat count.
      for i := 0; i < s.R; i++ {
         expanded = append(expanded, timelineSegment{StartTime: currentTime, Duration: s.D})
         currentTime += s.D
      }
   }
   return expanded
}

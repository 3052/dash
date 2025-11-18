package dash

import (
   "testing"
)

func TestSegmentTimeline_GetSegments(t *testing.T) {
   st := &SegmentTimeline{
      Segments: []*S{
         {D: 100},        // time: 0, duration: 100
         {D: 100, R: 1},  // time: 100, 200. duration 100 (2 segments total)
         {T: 500, D: 50}, // time: 500, duration: 50
         {D: 50},         // time: 550, duration: 50
      },
   }
   segments := st.GetSegments()
   expectedCount := 5
   if len(segments) != expectedCount {
      t.Fatalf("expected %d segments, got %d", expectedCount, len(segments))
   }

   expectedTimes := []uint64{0, 100, 200, 500, 550}
   for i, s := range segments {
      if s.StartTime != expectedTimes[i] {
         t.Errorf("segment %d: expected start time %d, got %d", i, expectedTimes[i], s.StartTime)
      }
   }
}

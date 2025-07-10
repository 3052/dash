package dash

var states = []struct{
   code string
   example string
}{
   {
      code: "Period.duration != nil",
   },
   {
      code: "Period.duration == nil",
   },
   {
      code: "Representation.SegmentBase != nil",
   },
   {
      code: "Representation.SegmentList != nil",
   },
   {
      code: "Representation.SegmentTemplate != nil",
   },
   {
      code: "SegmentTemplate.SegmentTimeline != nil",
   },
   {
      code: "SegmentTemplate.SegmentTimeline == nil",
   },
   {
      code: "SegmentTemplate.duration == 0",
   },
   {
      code: "SegmentTemplate.duration >= 1",
   },
   {
      code: "SegmentTemplate.endNumber == 0",
   },
   {
      code: "SegmentTemplate.endNumber >= 1",
   },
   {
      code: "SegmentTemplate.startNumber != nil",
   },
   {
      code: "SegmentTemplate.startNumber == nil",
   },
   {
      code: "SegmentTemplate.timescale != nil",
   },
   {
      code: "SegmentTemplate.timescale == nil",
   },
   {
      code: "URL.IsAbs(MPD.BaseURL)",
   },
   {
      code: "URL.IsAbs(MPD.BaseURL) == false",
   },
   {
      code: "len(MPD.Period) == 1",
   },
   {
      code: "len(MPD.Period) >= 2",
   },
   {
      code: `strings.Contains(SegmentTemplate.media, "$Number$")`,
   },
   {
      code: `strings.Contains(SegmentTemplate.media, "$Number%0")`,
   },
   {
      code: `strings.Contains(SegmentTemplate.media, "$Time$")`,
   },
}

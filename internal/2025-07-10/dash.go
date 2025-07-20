package dash

var states = []struct {
   state   string
   example []string
}{
   {
      state: "SegmentTemplate.startNumber != nil",
      example: []string{
         "paramount.mpd",
         "pluto.mpd",
      },
   },
   {
      state: "SegmentTemplate.startNumber == nil",
      example: []string{
         "criterion.mpd",
         "molotov.mpd",
         "rakuten.mpd",
         "rtbf.mpd",
      },
   },
   {
      state: "Representation.SegmentList != nil",
      example: []string{
         "criterion.mpd",
      },
   },
   {
      state: "SegmentTemplate.endNumber >= 1",
      example: []string{
         "molotov.mpd",
      },
   },
   {
      state: "len(MPD.Period) >= 2",
      example: []string{
         // technically Max fits here too, but omit since it uses SegmentBase
         "paramount.mpd",
      },
   },
   {
      state: `strings.Contains(SegmentTemplate.media, "$Number%0")`,
      example: []string{
         "pluto.mpd",
      },
   },
   {
      state: "SegmentTemplate.timescale == nil",
      example: []string{
         "rakuten.mpd",
      },
   },
   {
      state: "URL.IsAbs",
      example: []string{
         "rakuten.mpd",
      },
   },
   {
      state: "Representation.SegmentBase != nil",
      example: []string{
         "rakuten.mpd",
      },
   },
   {
      state: `strings.Contains(SegmentTemplate.media, "$Time$")`,
      example: []string{
         "rtbf.mpd",
      },
   },
   {
      state: "Period.duration != nil",
      example: []string{
         "criterion.mpd",
         "molotov.mpd",
         "paramount.mpd",
         "rtbf.mpd",
      },
   },
   {
      state: "URL.IsAbs == false",
      example: []string{
         "criterion.mpd",
         "molotov.mpd",
         "paramount.mpd",
         "pluto.mpd",
         "rtbf.mpd",
      },
   },
   {
      state: "len(MPD.Period) == 1",
      example: []string{
         "criterion.mpd",
         "molotov.mpd",
         "pluto.mpd",
         "rakuten.mpd",
         "rtbf.mpd",
      },
   },
   {
      state: `strings.Contains(SegmentTemplate.media, "$Number$")`,
      example: []string{
         "molotov.mpd",
         "paramount.mpd",
         "pluto.mpd",
      },
   },
   {
      state: "Period.duration == nil",
      example: []string{
         "pluto.mpd",
         "rakuten.mpd",
      },
   },
   {
      state: "Representation.SegmentTemplate != nil",
      example: []string{
         "molotov.mpd",
         "paramount.mpd",
         "pluto.mpd",
         "rakuten.mpd",
         "rtbf.mpd",
      },
   },
   {
      state: "SegmentTemplate.SegmentTimeline != nil",
      example: []string{
         "paramount.mpd",
         "pluto.mpd",
         "rtbf.mpd",
      },
   },
   {
      state: "SegmentTemplate.SegmentTimeline == nil",
      example: []string{
         "criterion.mpd",
         "molotov.mpd",
         "rakuten.mpd",
      },
   },
   {
      state: "SegmentTemplate.duration >= 1",
      example: []string{
         "criterion.mpd",
         "molotov.mpd",
         "paramount.mpd",
      },
   },
   {
      state: "SegmentTemplate.duration == 0",
      example: []string{
         "pluto.mpd",
         "rakuten.mpd",
         "rtbf.mpd",
      },
   },
   {
      state: "SegmentTemplate.endNumber == 0",
      example: []string{
         "criterion.mpd",
         "paramount.mpd",
         "pluto.mpd",
         "rakuten.mpd",
         "rtbf.mpd",
      },
   },
   {
      state: "SegmentTemplate.timescale != nil",
      example: []string{
         "criterion.mpd",
         "molotov.mpd",
         "paramount.mpd",
         "pluto.mpd",
         "rtbf.mpd",
      },
   },
}

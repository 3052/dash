package dash

var States = []struct {
   state   string
   example []string
}{
   {
      state: "len(MPD.Period) >= 2",
      example: []string{
         // technically Max fits here too, but omit since it uses SegmentBase
         "paramount.mpd",
      },
   },
   {
      state: "SegmentTemplate.startNumber == 0",
      example: []string{
         "max.mpd",
      },
   },
   {
      state: "SegmentTemplate.startNumber == nil",
      example: []string{
         "canal.mpd",
         "criterion.mpd",
         "molotov.mpd",
         "rakuten.mpd",
      },
   },
   {
      state: "SegmentTemplate.timescale == nil",
      example: []string{
         "rakuten.mpd",
      },
   },
   {
      state: "SegmentTemplate.timescale != nil",
      example: []string{
         "canal.mpd",
         "criterion.mpd",
         "max.mpd",
         "molotov.mpd",
         "paramount.mpd",
      },
   },
   {
      state: "Period.duration != nil",
      example: []string{
         "canal.mpd",
         "criterion.mpd",
         "max.mpd",
         "molotov.mpd",
         "paramount.mpd",
      },
   },
   {
      state: "Period.duration == nil",
      example: []string{
         "rakuten.mpd",
      },
   },
   {
      state: "Representation.SegmentBase != nil",
      example: []string{
         "max.mpd",
         "rakuten.mpd",
      },
   },
   {
      state: "Representation.SegmentTemplate != nil",
      example: []string{
         "canal.mpd",
         "molotov.mpd",
         "paramount.mpd",
      },
   },
   {
      state: "SegmentTemplate.SegmentTimeline != nil",
      example: []string{
         "canal.mpd",
         "max.mpd",
         "paramount.mpd",
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
         "max.mpd",
         "molotov.mpd",
         "paramount.mpd",
      },
   },
   {
      state: "SegmentTemplate.duration == 0",
      example: []string{
         "canal.mpd",
         "rakuten.mpd",
      },
   },
   {
      state: "SegmentTemplate.endNumber >= 1",
      example: []string{
         "molotov.mpd",
      },
   },
   {
      state: "SegmentTemplate.endNumber == 0",
      example: []string{
         "canal.mpd",
         "criterion.mpd",
         "max.mpd",
         "paramount.mpd",
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
      state: "URL.IsAbs == false",
      example: []string{
         "canal.mpd",
         "criterion.mpd",
         "max.mpd",
         "molotov.mpd",
         "paramount.mpd",
      },
   },
   {
      state: "len(MPD.Period) == 1",
      example: []string{
         "canal.mpd",
         "criterion.mpd",
         "molotov.mpd",
         "rakuten.mpd",
      },
   },
   {
      state: `strings.Contains(SegmentTemplate.media, "$Number$")`,
      example: []string{
         "max.mpd",
         "molotov.mpd",
         "paramount.mpd",
      },
   },
   {
      state: `strings.Contains(SegmentTemplate.media, "$Number%0")`,
      example: []string{
         "max.mpd",
      },
   },
   {
      state: `strings.Contains(SegmentTemplate.media, "$Time$")`,
      example: []string{
         "canal.mpd",
      },
   },
   {
      state: "Representation.SegmentList != nil",
      example: []string{
         "criterion.mpd",
      },
   },
}

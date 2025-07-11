package dash

var states = []struct {
   state   string
   example []string
}{
   {
      state: "Representation.SegmentBase != nil",
      example: []string{
         "max.mpd",
         "rakuten.mpd",
      },
   },
   {
      state: "Representation.SegmentList != nil",
      example: []string{
         "criterion.mpd",
      },
   },
   {
      state: "Representation.SegmentTemplate != nil",
      example: []string{
         "canal.mpd",
         "max.mpd",
         "molotov.mpd",
      },
   },
   {
      state: "SegmentTemplate.SegmentTimeline != nil",
      example: []string{
         "canal.mpd",
         "max.mpd",
      },
   },
   {
      state: "SegmentTemplate.SegmentTimeline == nil (endNumber or SegmentCount)",
      example: []string{
         "criterion.mpd",
         "molotov.mpd",
         "rakuten.mpd",
      },
   },
   {
      state: "SegmentTemplate.startNumber != nil",
      example: []string{
         "max.mpd",
      },
   },
   {
      state: "SegmentTemplate.startNumber == nil (startNumber = 1)",
      example: []string{
         "canal.mpd",
         "criterion.mpd",
         "molotov.mpd",
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
      },
   },
   {
      state: "SegmentTemplate.timescale == nil (timescale = 1)",
      example: []string{
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
      state: "SegmentTemplate.endNumber == 0 (SegmentTimeline or SegmentCount)",
      example: []string{
         "canal.mpd",
         "criterion.mpd",
         "max.mpd",
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
      },
   },
   {
      state: "len(MPD.Period) >= 2",
      example: []string{
         "max.mpd",
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
      state: `Period.duration != ""`,
      example: []string{
         "canal.mpd",
         "criterion.mpd",
         "max.mpd",
         "molotov.mpd",
      },
   },
   {
      state: `Period.duration == "" (Period.duration = MPD.mediaPresentationDuration)`,
      example: []string{
         "rakuten.mpd",
      },
   },
   {
      state: "SegmentTemplate.duration >= 1",
      example: []string{
         "criterion.mpd",
         "max.mpd",
         "molotov.mpd",
      },
   },
   {
      state: "SegmentTemplate.duration == 0 (endNumber or SegmentTimeline)",
      example: []string{
         "canal.mpd",
         "rakuten.mpd",
      },
   },
}

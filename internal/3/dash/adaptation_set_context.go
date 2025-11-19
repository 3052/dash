package dash

// AdaptationSetContext wraps an AdaptationSet with its parent PeriodContext.
type AdaptationSetContext struct {
   AdaptationSet *AdaptationSet
   Context       PeriodContext
}

package dash

// AdaptationSetScope wraps an AdaptationSet with its parent PeriodScope.
type AdaptationSetScope struct {
   AdaptationSet *AdaptationSet
   Scope         PeriodScope
}

package dash

// AdaptationSetNode wraps an AdaptationSet with its parent PeriodNode.
type AdaptationSetNode struct {
   AdaptationSet *AdaptationSet
   Node          PeriodNode
}

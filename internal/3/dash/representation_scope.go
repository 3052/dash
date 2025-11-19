package dash

// RepresentationScope wraps a Representation with its parent AdaptationSetScope.
type RepresentationScope struct {
   Representation *Representation
   Scope          AdaptationSetScope
}

// GetSegmentTemplateScope returns the SegmentTemplateScope for this representation.
// It looks for the SegmentTemplate in the Representation first, then falls back to the AdaptationSet.
// Returns nil if no SegmentTemplate is found.
func (rs RepresentationScope) GetSegmentTemplateScope() *SegmentTemplateScope {
   var st *SegmentTemplate
   if rs.Representation.SegmentTemplate != nil {
      st = rs.Representation.SegmentTemplate
   } else if rs.Scope.AdaptationSet.SegmentTemplate != nil {
      st = rs.Scope.AdaptationSet.SegmentTemplate
   }

   if st == nil {
      return nil
   }

   return &SegmentTemplateScope{
      SegmentTemplate: st,
      Scope:           rs,
   }
}

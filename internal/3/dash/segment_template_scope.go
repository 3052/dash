package dash

// SegmentTemplateScope wraps a SegmentTemplate with its parent RepresentationScope.
type SegmentTemplateScope struct {
   SegmentTemplate *SegmentTemplate
   Scope           RepresentationScope
}

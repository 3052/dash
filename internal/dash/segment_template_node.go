package dash

// SegmentTemplateNode wraps a SegmentTemplate with its parent RepresentationNode.
type SegmentTemplateNode struct {
   SegmentTemplate *SegmentTemplate
   Node            RepresentationNode
}

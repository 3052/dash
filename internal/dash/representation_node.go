package dash

// RepresentationNode wraps a Representation with its parent AdaptationSetNode.
type RepresentationNode struct {
   Representation *Representation
   Node           AdaptationSetNode
}

// GetSegmentTemplateNode returns the SegmentTemplateNode for this representation.
// It looks for the SegmentTemplate in the Representation first, then falls back to the AdaptationSet.
// Returns nil if no SegmentTemplate is found.
func (rn RepresentationNode) GetSegmentTemplateNode() *SegmentTemplateNode {
   var st *SegmentTemplate
   if rn.Representation.SegmentTemplate != nil {
      st = rn.Representation.SegmentTemplate
   } else if rn.Node.AdaptationSet.SegmentTemplate != nil {
      st = rn.Node.AdaptationSet.SegmentTemplate
   }

   if st == nil {
      return nil
   }

   return &SegmentTemplateNode{
      SegmentTemplate: st,
      Node:            rn,
   }
}

package dash

// PeriodNode wraps Period with a parent pointer.
type PeriodNode struct {
   Parent *MPD
   Value  *Period
}

// AdaptationSetNode wraps AdaptationSet with a parent pointer.
type AdaptationSetNode struct {
   Parent *Period
   Value  *AdaptationSet
}

// RepresentationNode wraps Representation with a parent pointer.
type RepresentationNode struct {
   Parent *AdaptationSet
   Value  *Representation
}

// SegmentTemplateNode wraps SegmentTemplate with a parent pointer.
type SegmentTemplateNode struct {
   // Parent can be AdaptationSet or Representation
   Parent interface{}
   Value  *SegmentTemplate
}

// SegmentListNode wraps SegmentList with a parent pointer.
type SegmentListNode struct {
   Parent *Representation
   Value  *SegmentList
}

// InitializationNode wraps Initialization with a parent pointer.
type InitializationNode struct {
   // Parent can be SegmentBase or SegmentList
   Parent interface{}
   Value  *Initialization
}

// SegmentURLNode wraps SegmentURL with a parent pointer.
type SegmentURLNode struct {
   Parent *SegmentList
   Value  *SegmentURL
}

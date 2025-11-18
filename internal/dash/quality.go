package dash

import "errors"

// RepresentationContext holds the parent elements that provide the necessary
// context for a Representation to be fully resolved, especially for generating
// segment URLs.
type RepresentationContext struct {
   Period        *Period
   AdaptationSet *AdaptationSet
}

// Quality represents a single logical quality level, identified by a unique
// Representation ID. It contains the core Representation data once, and a slice
// of all the different contexts in which this quality level appears (e.g.,
// once in a main content Period, and again in an ad break Period).
type Quality struct {
   // The embedded Representation provides direct access to all its fields
   // like ID, Bandwidth, Codecs, etc. This data is shared across all contexts.
   *Representation

   // Contexts is a slice of all the Period/AdaptationSet pairs where this
   // Representation ID was found.
   Contexts []*RepresentationContext
}

// ListMediaSegmentURLs generates the list of media segment URLs for this quality
// level, but specifically for the given context. The context must be one of the
// contexts available in the `Contexts` slice.
func (q *Quality) ListMediaSegmentURLs(ctx *RepresentationContext) ([]string, error) {
   if q == nil || q.Representation == nil {
      return nil, errors.New("quality or representation is nil")
   }
   if ctx == nil || ctx.Period == nil || ctx.AdaptationSet == nil {
      return nil, errors.New("context and its elements cannot be nil")
   }

   return q.Representation.ListMediaSegmentURLs(ctx.Period, ctx.AdaptationSet.SegmentTemplate)
}

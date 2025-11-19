package dash

// RepresentationContext wraps a Representation with its parent AdaptationSetContext.
type RepresentationContext struct {
   Representation *Representation
   Context        AdaptationSetContext
}

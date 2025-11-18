package dash

// Quality represents a single quality option (a Representation) enriched
// with inherited context from its parent AdaptationSet.
type Quality struct {
   *Representation
   Lang        string
   ContentType string
}

package dash

import (
   "errors"
   "net/url"
)

// RepresentationContext holds the parent elements that provide the necessary
// context for a Representation to be fully resolved, especially for generating
// segment URLs.
type RepresentationContext struct {
   Period        *Period
   AdaptationSet *AdaptationSet
}

// Quality represents a single logical quality level, identified by a unique
// Representation ID. It contains the core Representation data once, and a slice
// of all the different contexts in which this quality level appears.
type Quality struct {
   // The embedded Representation provides direct access to all its fields
   // like ID, Bandwidth, Codecs, etc. This data is shared across all contexts.
   *Representation

   // Contexts is a slice of all the Period/AdaptationSet pairs where this
   // Representation ID was found.
   Contexts []*RepresentationContext

   // parentMPD holds a reference to the top-level MPD object, making this
   // Quality object fully self-contained for URL resolution.
   parentMPD *MPD
}

// effectiveBaseURL calculates the final, resolved Base URL for a specific context
// by applying the hierarchy: MPD -> Period -> Representation.
func (q *Quality) effectiveBaseURL(ctx *RepresentationContext) (*url.URL, error) {
   if q == nil || q.parentMPD == nil {
      return nil, errors.New("quality or its parent MPD is nil")
   }

   // 1. Start with the MPD-level BaseURL.
   base, err := url.Parse(q.parentMPD.BaseURL)
   if err != nil {
      return nil, err
   }

   // 2. Resolve the Period-level BaseURL on top of it.
   if ctx.Period != nil && ctx.Period.BaseURL != "" {
      periodBase, err := url.Parse(ctx.Period.BaseURL)
      if err != nil {
         return nil, err
      }
      base = base.ResolveReference(periodBase)
   }

   // 3. Resolve the Representation-level BaseURL on top of that.
   if q.Representation != nil && q.Representation.BaseURL != "" {
      repBase, err := url.Parse(q.Representation.BaseURL)
      if err != nil {
         return nil, err
      }
      base = base.ResolveReference(repBase)
   }

   return base, nil
}

// AbsoluteInitializationURL returns the fully resolved URL for the initialization segment
// for a given context.
func (q *Quality) AbsoluteInitializationURL(ctx *RepresentationContext) (string, error) {
   if q == nil || q.Representation == nil {
      return "", errors.New("quality or representation is nil")
   }

   tpl := q.Representation.SegmentTemplate
   if tpl == nil && ctx.AdaptationSet != nil {
      tpl = ctx.AdaptationSet.SegmentTemplate
   }
   if tpl == nil || tpl.Initialization == "" {
      return "", errors.New("no initialization segment specified")
   }

   relativeURL := q.Representation.ResolveURL(tpl.Initialization)

   base, err := q.effectiveBaseURL(ctx)
   if err != nil {
      return "", err
   }

   relative, err := url.Parse(relativeURL)
   if err != nil {
      return "", err
   }

   return base.ResolveReference(relative).String(), nil
}

// AbsoluteMediaSegmentURLs generates the list of fully resolved, absolute URLs
// for all media segments for a given context.
func (q *Quality) AbsoluteMediaSegmentURLs(ctx *RepresentationContext) ([]string, error) {
   relativeURLs, err := q.ListMediaSegmentURLs(ctx)
   if err != nil {
      return nil, err
   }

   base, err := q.effectiveBaseURL(ctx)
   if err != nil {
      return nil, err
   }

   absoluteURLs := make([]string, len(relativeURLs))
   for i, relURL := range relativeURLs {
      relative, err := url.Parse(relURL)
      if err != nil {
         return nil, err
      }
      absoluteURLs[i] = base.ResolveReference(relative).String()
   }

   return absoluteURLs, nil
}

// ListMediaSegmentURLs generates the list of media segment URLs for this specific
// Quality object, using the context (Period, AdaptationSet) it was created with.
func (q *Quality) ListMediaSegmentURLs(ctx *RepresentationContext) ([]string, error) {
   if q == nil || q.Representation == nil {
      return nil, errors.New("quality or representation is nil")
   }

   var asTpl *SegmentTemplate
   if ctx.AdaptationSet != nil {
      asTpl = ctx.AdaptationSet.SegmentTemplate
   }

   return q.Representation.ListMediaSegmentURLs(ctx.Period, asTpl)
}

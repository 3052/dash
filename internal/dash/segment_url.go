package dash

// SegmentURL defines a specific media segment source.
type SegmentURL struct {
   Media string `xml:"media,attr,omitempty"`

   // Navigation
   Parent *SegmentList `xml:"-"`
}

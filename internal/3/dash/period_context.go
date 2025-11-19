package dash

// PeriodContext wraps a Period with a pointer to its parent MPD.
type PeriodContext struct {
   Period *Period
   MPD    *MPD
}

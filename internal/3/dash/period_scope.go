package dash

// PeriodScope wraps a Period with a pointer to its parent MPD.
type PeriodScope struct {
   Period *Period
   MPD    *MPD
}

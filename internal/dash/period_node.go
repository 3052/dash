package dash

// PeriodNode wraps a Period with a pointer to its parent MPD.
type PeriodNode struct {
   Period *Period
   MPD    *MPD
}

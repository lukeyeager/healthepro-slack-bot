package week

import "time"

// DisplayMonday returns the Monday of the week to display.
// On Saturday or Sunday the upcoming week is returned, so the page
// always shows a full Mon–Fri block rather than a mostly-past one.
func DisplayMonday(now time.Time) time.Time {
	var daysToMonday int
	switch now.Weekday() {
	case time.Saturday:
		daysToMonday = 2
	case time.Sunday:
		daysToMonday = 1
	default: // Monday=0 offset, Tue=−1, … Fri=−4
		daysToMonday = -(int(now.Weekday()) - int(time.Monday))
	}
	d := now.AddDate(0, 0, daysToMonday)
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

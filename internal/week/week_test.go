package week

import (
	"testing"
	"time"
)

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 12, 0, 0, 0, time.UTC)
}

func TestDisplayMonday(t *testing.T) {
	tests := []struct {
		name       string
		input      time.Time
		wantMonday time.Time
	}{
		{"monday returns self", date(2026, 4, 20), date(2026, 4, 20)},
		{"tuesday returns prev monday", date(2026, 4, 21), date(2026, 4, 20)},
		{"friday returns prev monday", date(2026, 4, 24), date(2026, 4, 20)},
		{"saturday returns next monday", date(2026, 4, 25), date(2026, 4, 27)},
		{"sunday returns next monday", date(2026, 4, 26), date(2026, 4, 27)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DisplayMonday(tt.input)
			wantStr := tt.wantMonday.Format("2006-01-02")
			gotStr := got.Format("2006-01-02")
			if gotStr != wantStr {
				t.Errorf("DisplayMonday(%s) = %s, want %s",
					tt.input.Format("2006-01-02"), gotStr, wantStr)
			}
		})
	}
}

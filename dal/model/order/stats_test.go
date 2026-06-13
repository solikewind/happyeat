package order

import (
	"testing"
	"time"

	entschema "github.com/solikewind/happyeat/dal/model/ent/schema"
)

func TestParseStatsDateRange(t *testing.T) {
	today := dateOnlyCST(time.Now())
	r, err := ParseStatsDateRange("", "")
	if err != nil {
		t.Fatal(err)
	}
	if !r.Start.Equal(today) || !r.End.Equal(today) {
		t.Fatalf("default range want today %v, got %v-%v", today, r.Start, r.End)
	}

	r, err = ParseStatsDateRange("2026-06-11", "2026-06-13")
	if err != nil {
		t.Fatal(err)
	}
	wantStart := time.Date(2026, 6, 11, 0, 0, 0, 0, entschema.CST)
	wantEnd := time.Date(2026, 6, 13, 0, 0, 0, 0, entschema.CST)
	if !r.Start.Equal(wantStart) || !r.End.Equal(wantEnd) {
		t.Fatalf("range mismatch: %v-%v", r.Start, r.End)
	}
	if !r.EndExclusive.Equal(wantEnd.AddDate(0, 0, 1)) {
		t.Fatalf("end exclusive = %v", r.EndExclusive)
	}
}

func TestParseStatsDateRangeSwap(t *testing.T) {
	r, err := ParseStatsDateRange("2026-06-13", "2026-06-11")
	if err != nil {
		t.Fatal(err)
	}
	if r.Start.After(r.End) {
		t.Fatal("start should not be after end")
	}
}

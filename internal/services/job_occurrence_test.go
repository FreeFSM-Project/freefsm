package services

import (
	"testing"
	"time"
)

func TestResetSubtasksMarksAllIncomplete(t *testing.T) {
	t.Parallel()

	subtasks := resetSubtasks([]JobSubtask{
		{Title: "Inspect", Completed: true, SortOrder: 1},
		{Title: "Report", Completed: false, SortOrder: 2},
	})

	for _, subtask := range subtasks {
		if subtask.Completed {
			t.Fatalf("subtask %q remained completed", subtask.Title)
		}
	}
}

func TestShiftVisitsMovesVisitDatesByOccurrenceDateDelta(t *testing.T) {
	t.Parallel()

	loc := time.FixedZone("test", -6*60*60)
	sourceStart := time.Date(2026, 1, 10, 9, 0, 0, 0, loc)
	nextStart := time.Date(2026, 1, 17, 14, 0, 0, 0, loc)
	visits := shiftVisits([]JobVisit{
		{Date: "2026-01-10", StartTime: "09:00", EndTime: "10:00"},
		{Date: "2026-01-12", StartTime: "11:00", EndTime: "12:00"},
		{Date: "not-a-date", StartTime: "13:00", EndTime: "14:00"},
	}, &sourceStart, nextStart)

	if got, want := visits[0].Date, "2026-01-17"; got != want {
		t.Fatalf("first visit date = %q, want %q", got, want)
	}
	if got, want := visits[1].Date, "2026-01-19"; got != want {
		t.Fatalf("second visit date = %q, want %q", got, want)
	}
	if got, want := visits[2].Date, "not-a-date"; got != want {
		t.Fatalf("invalid visit date = %q, want %q", got, want)
	}
}

package main

import (
	"log"
	"testing"
	"time"
)

func TestReversedTimePeriods(t *testing.T) {
	from := time.Date(2022, 2, 23, 11, 0, 0, 0, time.UTC)
	to := time.Date(2022, 2, 24, 4, 0, 0, 0, time.UTC)
	step := 4 * time.Hour
	expected := []TimePeriod{
		{
			To:   time.Date(2022, 2, 24, 4, 0, 0, 0, time.UTC),
			From: time.Date(2022, 2, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			To:   time.Date(2022, 2, 24, 0, 0, 0, 0, time.UTC),
			From: time.Date(2022, 2, 23, 20, 0, 0, 0, time.UTC),
		},
		{
			To:   time.Date(2022, 2, 23, 20, 0, 0, 0, time.UTC),
			From: time.Date(2022, 2, 23, 16, 0, 0, 0, time.UTC),
		},
		{
			To:   time.Date(2022, 2, 23, 16, 0, 0, 0, time.UTC),
			From: time.Date(2022, 2, 23, 12, 0, 0, 0, time.UTC),
		},
		{
			To:   time.Date(2022, 2, 23, 12, 0, 0, 0, time.UTC),
			From: time.Date(2022, 2, 23, 11, 0, 0, 0, time.UTC),
		},
	}
	periods := ReversedTimePeriods(from, to, step)
	assertReversedTimePeriods(t, expected, periods)

	// what to is after from
	from = time.Date(2022, 2, 24, 4, 0, 0, 0, time.UTC)
	to = time.Date(2022, 2, 23, 11, 0, 0, 0, time.UTC)
	step = 4 * time.Hour
	expected = []TimePeriod{}
	periods = ReversedTimePeriods(from, to, step)
	assertReversedTimePeriods(t, expected, periods)
}

func assertReversedTimePeriods(t *testing.T, expected, actual []TimePeriod) {
	for _, period := range actual {
		log.Printf("%v", period)
	}
	if len(expected) != len(actual) {
		t.Fatalf("expected %d periods, got %d", len(expected), len(actual))
	}

	for i, period := range expected {
		if period.From != actual[i].From || period.To != actual[i].To {
			t.Fatalf("expected %v, got %v", period, actual[i])
		}
	}
}

package main

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(0, 0, 0, 10, 32, 0, 0, time.UTC),
			want: "10:32",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "PM",
			tm:   time.Date(0, 0, 0, 13, 48, 0, 0, time.UTC),
			want: "13:48",
		},
		{
			name: "CET",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "09:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // <== Run each sub-test in parallel
			hd := humanDate(tt.tm)
			t.Logf("testing indexing %q for %q", tt.name, tt.want)
			assert.Equal(t, hd, tt.want)
		})
	}
}

package decor

import (
	"testing"
	"time"
)

func TestChooseTimeProducer(t *testing.T) {
	cases := []struct {
		name      string
		style     TimeStyle
		remaining time.Duration
		expected  string
	}{
		{
			name:      "HHMMSS under an hour",
			style:     ET_STYLE_HHMMSS,
			remaining: 5*time.Minute + 6*time.Second,
			expected:  "00:05:06",
		},
		{
			name:      "HHMMSS within a day",
			style:     ET_STYLE_HHMMSS,
			remaining: 10*time.Hour + 5*time.Minute + 6*time.Second,
			expected:  "10:05:06",
		},
		{
			name:      "HHMMSS hours do not wrap at 60",
			style:     ET_STYLE_HHMMSS,
			remaining: 70*time.Hour + 5*time.Minute + 6*time.Second,
			expected:  "70:05:06",
		},
		{
			name:      "HHMM hours do not wrap at 60",
			style:     ET_STYLE_HHMM,
			remaining: 130*time.Hour + 5*time.Minute,
			expected:  "130:05",
		},
		{
			name:      "MMSS keeps MM:SS below one hour",
			style:     ET_STYLE_MMSS,
			remaining: 5*time.Minute + 6*time.Second,
			expected:  "05:06",
		},
		{
			name:      "MMSS switches to HH:MM:SS past 60 hours",
			style:     ET_STYLE_MMSS,
			remaining: 70*time.Hour + 5*time.Minute + 6*time.Second,
			expected:  "70:05:06",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := chooseTimeProducer(test.style)(test.remaining)
			if got != test.expected {
				t.Errorf("Want: %q, Got: %q\n", test.expected, got)
			}
		})
	}
}

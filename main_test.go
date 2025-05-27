package durationiso8601

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	type args struct {
		t time.Time
		s string
	}

	baseTime, _ := time.Parse(time.RFC3339, "2023-01-01T12:00:00Z")
	leapYearTime, _ := time.Parse(time.RFC3339, "2020-02-01T12:00:00Z") // Leap year for testing

	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name: "empty string",
			args: args{
				t: baseTime,
				s: "",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "invalid format",
			args: args{
				t: baseTime,
				s: "P1X",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "simple seconds",
			args: args{
				t: baseTime,
				s: "PT30S",
			},
			want:    30 * time.Second,
			wantErr: false,
		},
		{
			name: "simple minutes",
			args: args{
				t: baseTime,
				s: "PT5M",
			},
			want:    5 * time.Minute,
			wantErr: false,
		},
		{
			name: "simple hours",
			args: args{
				t: baseTime,
				s: "PT2H",
			},
			want:    2 * time.Hour,
			wantErr: false,
		},
		{
			name: "simple days",
			args: args{
				t: baseTime,
				s: "P3D",
			},
			want:    3 * 24 * time.Hour,
			wantErr: false,
		},
		{
			name: "complex duration",
			args: args{
				t: baseTime,
				s: "P1DT2H30M45S",
			},
			want:    (24+2)*time.Hour + 30*time.Minute + 45*time.Second,
			wantErr: false,
		},
		{
			name: "months basic",
			args: args{
				t: baseTime,
				s: "P2M",
			},
			want:    31*24*time.Hour + 28*24*time.Hour, // Jan(31) + Feb(28) in non-leap year
			wantErr: false,
		},
		{
			name: "months in leap year",
			args: args{
				t: leapYearTime,
				s: "P1M",
			},
			want:    29 * 24 * time.Hour, // Feb has 29 days in leap year
			wantErr: false,
		},
		{
			name: "years basic",
			args: args{
				t: baseTime,
				s: "P1Y",
			},
			want:    365 * 24 * time.Hour, // Non-leap year
			wantErr: false,
		},
		{
			name: "years starting in leap year",
			args: args{
				t: leapYearTime,
				s: "P1Y",
			},
			want:    366 * 24 * time.Hour, // Leap year
			wantErr: false,
		},
		{
			name: "complex with all units",
			args: args{
				t: baseTime,
				s: "P1Y2M3DT4H5M6S",
			},
			want:    baseTime.AddDate(1, 2, 3).Add(4*time.Hour + 5*time.Minute + 6*time.Second).Sub(baseTime),
			wantErr: false,
		},
		{
			name: "negative duration",
			args: args{
				t: baseTime,
				s: "-P1DT2H",
			},
			want:    -((24 + 2) * time.Hour),
			wantErr: false,
		},
		{
			name: "decimal values",
			args: args{
				t: baseTime,
				s: "PT1.5H",
			},
			want:    90 * time.Minute,
			wantErr: false,
		},
		{
			name: "week format",
			args: args{
				t: baseTime,
				s: "P2W",
			},
			want:    14 * 24 * time.Hour,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.args.t, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

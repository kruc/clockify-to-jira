package main

import (
	"testing"
)

func Test_dosko(t *testing.T) {
	type args struct {
		timeSpentSeconds int
		stachurskyMode   int
	}
	tests := []struct {
		name  string
		args  args
		want1 int
		want2 string
		want3 string
	}{
		{
			name: "Dosko 1m - org 16m",
			args: args{
				timeSpentSeconds: 960, // 16min
				stachurskyMode:   1,
			},
			want1: 960,
			want2: "16m0s",
			want3: "16m0s",
		},
		{
			name: "Dosko 15m - org 16m",
			args: args{
				timeSpentSeconds: 960, // 16min
				stachurskyMode:   15,
			},
			want1: 900,
			want2: "16m0s",
			want3: "15m0s",
		},
		{
			name: "Dosko 15m - org 2m",
			args: args{
				timeSpentSeconds: 120, // 2min
				stachurskyMode:   15,
			},
			want1: 900,
			want2: "2m0s",
			want3: "15m0s",
		},
		{
			name: "Dosko 15m - org 22m29s",
			args: args{
				timeSpentSeconds: 1349, // 22min49s
				stachurskyMode:   15,
			},
			want1: 900,
			want2: "22m29s",
			want3: "15m0s",
		},
		{
			name: "Dosko 15m - org 22m30s",
			args: args{
				timeSpentSeconds: 1350, // 22min30s
				stachurskyMode:   15,
			},
			want1: 1800,
			want2: "22m30s",
			want3: "30m0s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2, got3 := dosko(tt.args.timeSpentSeconds, tt.args.stachurskyMode)
			if got1 != tt.want1 {
				t.Errorf("dosko() got = %d, want %d", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("dosko() got = %s, want %s", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("dosko() got = %s, want %s", got3, tt.want3)
			}
		})
	}
}

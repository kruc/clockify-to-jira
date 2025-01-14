package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_adjustClockifyDate(t *testing.T) {
	type args struct {
		clockifyDate time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "Dirty hack - some api issues - to fix in code",
			args: args{
				clockifyDate: time.Unix(1616683240, 0),
			},
			want: time.Unix(1616683240, 1000000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := adjustClockifyDate(tt.args.clockifyDate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("adjustClockifyDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseIssueID(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Parse issue id with []",
			args: args{"[ID-123] Some description"},
			want: "ID-123",
		},
		{
			name: "Parse issue id without []",
			args: args{"ID-123 Some description"},
			want: "ID-123",
		},
		{
			name: "Parse issue id with :",
			args: args{"ID-123: Some description"},
			want: "ID-123",
		},
		{
			name: "Parse issue id with : and []",
			args: args{"[ID-123]: Some description"},
			want: "ID-123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseIssueID(tt.args.value); got != tt.want {
				t.Errorf("parseIssueID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseIssueComment(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Parse comment v1",
			args: args{"[ID-123]: Some description 1"},
			want: "Some description 1",
		},
		{
			name: "Parse comment v2",
			args: args{"ID-123: Some description 2"},
			want: "Some description 2",
		},
		{
			name: "Parse comment v3",
			args: args{"ID-123 Some description 3"},
			want: "Some description 3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseIssueComment(tt.args.value); got != tt.want {
				t.Errorf("parseIssueComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTimeDiff(t *testing.T) {
	type args struct {
		start time.Time
		stop  time.Time
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Normal diff",
			args: args{
				start: time.Now(),
				stop:  time.Now().Add(10 * time.Second),
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTimeDiff(tt.args.start, tt.args.stop); got != tt.want {
				t.Errorf("getTimeDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

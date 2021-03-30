package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/kruc/clockify-api/gctag"
)

func Test_removeTag(t *testing.T) {
	type args struct {
		tagsList    []string
		tagToRemove string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Remove exists tag",
			args: args{
				tagsList:    []string{"tag1", "tag2", "tag3"},
				tagToRemove: "tag1",
			},
			want: []string{"tag2", "tag3"},
		},
		{
			name: "Remove non exists tag",
			args: args{
				tagsList:    []string{"tag1", "tag2", "tag3"},
				tagToRemove: "tagX",
			},
			want: []string{"tag1", "tag2", "tag3"},
		},
		{
			name: "Remove empty tag",
			args: args{
				tagsList:    []string{"tag1", "tag2", "tag3"},
				tagToRemove: "",
			},
			want: []string{"tag1", "tag2", "tag3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeTag(tt.args.tagsList, tt.args.tagToRemove); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dosko(t *testing.T) {
	type args struct {
		timeSpentSeconds int
		stachurskyMode   int
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 doskoDebugInfo
	}{
		{
			name: "Dosko 1m - org 16m",
			args: args{
				timeSpentSeconds: 960, // 16min
				stachurskyMode:   1,
			},
			want: 960,
			want1: doskoDebugInfo{
				originalTime: "16m0s",
				doskoTime:    "16m0s",
			},
		},
		{
			name: "Dosko 15m - org 16m",
			args: args{
				timeSpentSeconds: 960, // 16min
				stachurskyMode:   15,
			},
			want: 900,
			want1: doskoDebugInfo{
				originalTime: "16m0s",
				doskoTime:    "15m0s",
			},
		},
		{
			name: "Dosko 15m - org 2m",
			args: args{
				timeSpentSeconds: 120, // 2min
				stachurskyMode:   15,
			},
			want: 900,
			want1: doskoDebugInfo{
				originalTime: "2m0s",
				doskoTime:    "15m0s",
			},
		},
		{
			name: "Dosko 15m - org 22m29s",
			args: args{
				timeSpentSeconds: 1349, // 22min49s
				stachurskyMode:   15,
			},
			want: 900,
			want1: doskoDebugInfo{
				originalTime: "22m29s",
				doskoTime:    "15m0s",
			},
		},
		{
			name: "Dosko 15m - org 22m30s",
			args: args{
				timeSpentSeconds: 1350, // 22min30s
				stachurskyMode:   15,
			},
			want: 1800,
			want1: doskoDebugInfo{
				originalTime: "22m30s",
				doskoTime:    "30m0s",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := dosko(tt.args.timeSpentSeconds, tt.args.stachurskyMode)
			if got != tt.want {
				t.Errorf("dosko() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("dosko() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_displayTagsName(t *testing.T) {
	type args struct {
		tags []gctag.Tag
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Regular",
			args: args{
				tags: []gctag.Tag{
					{
						Name: "tag1",
					},
					{
						Name: "tag23",
					},
				},
			},
			want: []string{"tag1", "tag23"},
		},
		{
			name: "Empty",
			args: args{
				tags: []gctag.Tag{},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := displayTagsName(tt.args.tags); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("displayTagsName() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTimeDiff(tt.args.start, tt.args.stop); got != tt.want {
				t.Errorf("getTimeDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

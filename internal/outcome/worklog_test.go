package outcome

import (
	"testing"
	"time"

	"github.com/kruc/clockify-to-jira/internal/assert"
)

func TestGetWorklog(t *testing.T) {

	dateTime, _ := time.Parse(timeFormat, "2024-09-16 06:00:00")

	data := WorklogData{
		Description: "Time entry description",
		Workspace:   "Workspace",
		Client:      "Client",
		Project:     "Project",
		Date:        dateTime,
		TimeSpent: DoskoDetails{
			OriginalTime: "8h7m0s",
			RoundedTime:  "8h0m0s",
			Precision:    15,
		},
		Comment: "Comment",
		Tags:    []string{"Tag1", "Tag2"},
	}

	got := data.GetSummary()

	want := `Worklog: Time entry description
---------
Workspace: Workspace
Client: Client
Project: Project
Date: 2024-09-16 06:00:00
Time spent: 8h0m0s (clockify: 8h7m0s stachurskyMode: 15m)
Comment: Comment
Tags: [Tag1 Tag2]
---------
`
	assert.Strings(t, got, want)
}

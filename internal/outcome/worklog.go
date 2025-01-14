package outcome

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

const (
	worklogTemplate = `Worklog: {{.Description}}
---------
Workspace: {{.Workspace}}
Client: {{.Client}}
Project: {{.Project}}
Date: {{.Date}}
Time spent: {{.TimeSpent}}
Comment: {{.Comment}}
Tags: {{.Tags}}
---------
`
)

type DoskoDetails struct {
	OriginalTime string
	RoundedTime  string
	Precision    int
}

func (dd *DoskoDetails) toString() string {
	return fmt.Sprintf("%+v (clockify: %+v stachurskyMode: %+vm)", dd.RoundedTime, dd.OriginalTime, dd.Precision)
}

type WorklogData struct {
	Description string
	Workspace   string
	Client      string
	Project     string
	Date        time.Time
	TimeSpent   DoskoDetails
	Comment     string
	Tags        []string
}

type Worklog struct {
	Description string
	Workspace   string
	Client      string
	Project     string
	Date        string
	TimeSpent   string
	Comment     string
	Tags        []string
}

func (w *WorklogData) GetSummary() string {
	worklog := w.prepareWorklogData()

	var output bytes.Buffer

	t := template.Must(template.New("worklog").Parse(template.HTMLEscapeString(worklogTemplate)))

	t.Execute(&output, worklog)

	return output.String()
}

func (w *WorklogData) prepareWorklogData() Worklog {

	return Worklog{
		Description: w.Description,
		Workspace:   w.Workspace,
		Client:      w.Client,
		Project:     w.Project,
		Date:        w.Date.Format(timeFormat),
		TimeSpent:   w.TimeSpent.toString(),
		Comment:     w.Comment,
		Tags:        w.Tags,
	}
}

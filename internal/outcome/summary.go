package outcome

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

const (
	timeFormat                            = "2006-01-02 15:04:05"
	ErrSummaryParseTotalTimeDurationError = OutcomeErr("Cannot parse given string")
	summaryTemplate                       = `Workspace: {{.Workspace}}
-------
SUMMARY
-------
Time entries range: {{.Start}} - {{.End}}
Number of time entries: {{.TimeEntriesNumber}}
Total time: {{.TotalTime}}
Total dosko: {{.TotalDoskoTime}} (t={{.Dosko}}m)
---------
`
)

type OutcomeErr string

func (e OutcomeErr) Error() string {
	return string(e)
}

type SummaryData struct {
	Workspace      string
	Start          time.Time
	End            time.Time
	entriesCount   int
	totalTime      int
	totalDoskoTime int
	doskoFactor    int
}

type Summary struct {
	Workspace         string
	Start             string
	End               string
	TimeEntriesNumber int
	TotalTime         string
	TotalDoskoTime    string
	Dosko             int
}

func (d *SummaryData) IncreaseTimeEntryCount() {
	d.entriesCount++
}

func (d *SummaryData) AddTimeEntryDuration(timeEntryDuration int) {
	d.totalTime += timeEntryDuration
}

func (d *SummaryData) AddDoskoTimeEntryDuration(timeEntryDuration int) {
	d.totalDoskoTime += timeEntryDuration
}

func (d *SummaryData) AddDoskoFactor(doskoFactor int) {
	d.doskoFactor = doskoFactor
}

func (d *SummaryData) getTotalTime(totalTime int) (string, error) {
	parsedDuration, err := time.ParseDuration(fmt.Sprintf("%ds", totalTime))

	if err != nil {
		return "", ErrSummaryParseTotalTimeDurationError
	}

	return parsedDuration.String(), nil
}

func (d *SummaryData) GetSummary() (string, error) {

	summary, err := d.prepareSummaryData()

	if err != nil {
		return "", err
	}

	var output bytes.Buffer

	t := template.Must(template.New("summary").Parse(summaryTemplate))

	t.Execute(&output, summary)

	// Determine how to test it
	// err = t.Execute(&output, summary)

	// if err != nil {
	// 	return "", err
	// }

	return output.String(), nil
}

func (d *SummaryData) prepareSummaryData() (Summary, error) {

	totalTime, err := d.getTotalTime(d.totalTime)

	if err != nil {
		return Summary{}, err
	}

	totalDoskoTime, err := d.getTotalTime(d.totalDoskoTime)

	if err != nil {
		return Summary{}, err
	}

	summary := Summary{
		Workspace:         d.Workspace,
		Start:             d.Start.Format(timeFormat),
		End:               d.End.Format(timeFormat),
		TimeEntriesNumber: d.entriesCount,
		TotalTime:         totalTime,
		TotalDoskoTime:    totalDoskoTime,
		Dosko:             d.doskoFactor,
	}

	return summary, nil
}

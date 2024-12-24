package outcome

import (
	"testing"
	"time"

	"github.com/kruc/clockify-to-jira/internal/assert"
)

func TestSummaryData(t *testing.T) {

	data := SummaryData{
		entriesCount:   12,
		totalTime:      100,
		totalDoskoTime: 20,
		doskoFactor:    5,
	}

	t.Run("Increase timeentry count", func(t *testing.T) {

		data.IncreaseTimeEntryCount()

		assert.Ints(t, data.entriesCount, 13)
	})

	t.Run("Add timeentry duration", func(t *testing.T) {

		data.AddTimeEntryDuration(120)

		assert.Ints(t, data.totalTime, 220)
	})

	t.Run("Add dosko timeentry duration", func(t *testing.T) {

		data.AddDoskoTimeEntryDuration(20)

		assert.Ints(t, data.totalDoskoTime, 40)
	})

	t.Run("Add dosko factor", func(t *testing.T) {

		data.AddDoskoFactor(30)

		assert.Ints(t, data.doskoFactor, 30)
	})

	t.Run("Get total time in proper format", func(t *testing.T) {

		s, _ := data.getTotalTime(20)

		assert.Strings(t, s, "20s")

		ms, _ := data.getTotalTime(250)

		assert.Strings(t, ms, "4m10s")

		hms, _ := data.getTotalTime(3600)

		assert.Strings(t, hms, "1h0m0s")
	})

	t.Run("Get total time return error on invalid data", func(t *testing.T) {

		_, err := data.getTotalTime(12312312322)

		assert.Errors(t, err, ErrSummaryParseTotalTimeDurationError)
	})
}

func TestGetSummary(t *testing.T) {

	t.Run("Get templated summary", func(t *testing.T) {

		data := SummaryData{
			Workspace:      "WorkspaceKey",
			Start:          time.Date(2024, time.April, 11, 21, 34, 01, 0, time.UTC),
			End:            time.Date(2024, time.May, 11, 21, 34, 01, 0, time.UTC),
			entriesCount:   12,
			totalTime:      100,
			totalDoskoTime: 200,
			doskoFactor:    5,
		}

		got, _ := data.GetSummary()

		want := `Workspace: WorkspaceKey
-------
SUMMARY
-------
Time entries range: 2024-04-11 21:34:01 - 2024-05-11 21:34:01
Number of time entries: 12
Total time: 1m40s
Total dosko: 3m20s (t=5m)
---------
`
		assert.Strings(t, got, want)
	})

	t.Run("Get errors on invalid totalTime input data", func(t *testing.T) {
		data := SummaryData{
			totalTime: 10009283729293,
		}

		_, err := data.GetSummary()

		assert.Errors(t, err, ErrSummaryParseTotalTimeDurationError)
	})

	t.Run("Get errors on invalid totalDoskoTime input data", func(t *testing.T) {
		data := SummaryData{
			totalDoskoTime: 999283729293,
		}

		_, err := data.GetSummary()

		assert.Errors(t, err, ErrSummaryParseTotalTimeDurationError)
	})
}

func TestError(t *testing.T) {
	t.Run("ErrNotFound", func(t *testing.T) {
		got := OutcomeErr("Error message").Error()
		want := "Error message"

		assert.Strings(t, got, want)
	})
}

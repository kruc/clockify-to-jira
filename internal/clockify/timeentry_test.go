package clockify

import (
	"testing"
	"time"

	"github.com/kruc/clockify-to-jira/internal/assert"
)

func TestGetTimeEntries(t *testing.T) {

	t.Run("Get TimeEntries", func(t *testing.T) {

		fakeClient := &fakeClient{}
		fakeClient.getMeSuccessResponse()
		fakeClient.logRangeSuccessResponse()

		initClient = func(string) (clockifyApiClient, error) {
			return fakeClient, nil
		}

		apiClient, _ := NewClient("token")

		end := time.Now()
		start := time.Now().AddDate(0, 0, -5)

		timeEntries, err := apiClient.GetTimeEntriesFromGivenPeriod(start, end, "ws1")

		assert.Ints(t, len(timeEntries), 2)
		assert.Errors(t, err, nil)
		assert.Strings(t, timeEntries[0].Description, "XYZ-123 Timentry description")
		assert.Strings(t, timeEntries[0].ID, "id1")
		assert.Strings(t, timeEntries[0].ClientName, "clientName1")
		assert.Strings(t, timeEntries[0].ProjectID, "projectID1")
		assert.Strings(t, timeEntries[0].ProjectName, "projectName1")
		assert.Strings(t, timeEntries[0].Start.String(), "2025-01-08 10:30:00 +0000 UTC")
		assert.Strings(t, timeEntries[0].End.String(), "2025-01-08 17:00:00 +0000 UTC")
		assert.Strings(t, timeEntries[0].Duration, "PT7H30M")
		assert.Strings(t, timeEntries[0].Tags["tagName1"].ID, "tagId1")
		assert.Strings(t, timeEntries[0].Tags["tagName1"].Name, "tagName1")
		assert.Strings(t, timeEntries[0].Tags["tagName1"].WorkspaceID, "ws1")
		assert.Strings(t, timeEntries[0].Tags["tagName2"].ID, "tagId2")
		assert.Strings(t, timeEntries[0].Tags["tagName2"].Name, "tagName2")
		assert.Strings(t, timeEntries[0].Tags["tagName2"].WorkspaceID, "ws1")

		assert.Strings(t, timeEntries[1].ID, "id2")
		assert.Strings(t, timeEntries[1].Description, "ABC-123 Timentry description")
		assert.Strings(t, timeEntries[1].ClientName, "clientName2")
		assert.Strings(t, timeEntries[1].ProjectID, "projectID2")
		assert.Strings(t, timeEntries[1].ProjectName, "projectName2")
		assert.Strings(t, timeEntries[1].Start.String(), "2025-01-09 08:00:00 +0000 UTC")
		assert.Strings(t, timeEntries[1].End.String(), "2025-01-09 16:15:00 +0000 UTC")
		assert.Strings(t, timeEntries[1].Duration, "PT8H15M")
		assert.Strings(t, timeEntries[1].Tags["tagName3"].ID, "tagId3")
		assert.Strings(t, timeEntries[1].Tags["tagName3"].Name, "tagName3")
		assert.Strings(t, timeEntries[1].Tags["tagName3"].WorkspaceID, "ws1")
		assert.Strings(t, timeEntries[1].Tags["tagName4"].ID, "tagId4")
		assert.Strings(t, timeEntries[1].Tags["tagName4"].Name, "tagName4")
		assert.Strings(t, timeEntries[1].Tags["tagName4"].WorkspaceID, "ws1")
	})

	t.Run("Return error on fetching logged in user", func(t *testing.T) {

		fakeClient := &fakeClient{}
		fakeClient.getMeErrorResponse()

		initClient = func(string) (clockifyApiClient, error) {
			return fakeClient, nil
		}

		apiClient, _ := NewClient("token")

		end := time.Now()
		start := time.Now().AddDate(0, 0, -5)

		_, err := apiClient.GetTimeEntriesFromGivenPeriod(start, end, "ws1")

		assert.Errors(t, err, ErrClockifyFailToFetchLoggedInUserData)
	})

	t.Run("Return error on fetching time entries", func(t *testing.T) {

		fakeClient := &fakeClient{}
		fakeClient.getMeSuccessResponse()
		fakeClient.logRangeErrorResponse()

		initClient = func(string) (clockifyApiClient, error) {
			return fakeClient, nil
		}

		apiClient, _ := NewClient("token")

		end := time.Now()
		start := time.Now().AddDate(0, 0, -5)

		_, err := apiClient.GetTimeEntriesFromGivenPeriod(start, end, "ws1")

		assert.Errors(t, err, ErrClockifyFailToFetchTimeEntries)
	})
}

func TestUpdateTimeEntry(t *testing.T) {

	t.Run("Update TimeEntry", func(t *testing.T) {
		fakeClient := &fakeClient{}
		fakeClient.updateTimeEntrySuccessResponse()

		initClient = func(string) (clockifyApiClient, error) {
			return fakeClient, nil
		}

		apiClient, _ := NewClient("token")
		workspaceID := "ws1"

		end := time.Date(1986, time.January, 5, 10, 46, 28, 0, &time.Location{})
		start := time.Date(1986, time.January, 1, 10, 46, 28, 0, &time.Location{})

		timeEntry := TimeEntry{
			ID:          "timeEntryID",
			Description: "description",
			ProjectID:   "projectID",
			Start:       start,
			End:         &end,
			Tags: map[string]Tag{
				"tagName1": {ID: "tagId1"},
				"tagName2": {ID: "tagId2"},
			},
		}

		timeEntry, err := apiClient.UpdateTimeEntry(workspaceID, timeEntry)

		assert.Errors(t, err, nil)
		assert.Strings(t, timeEntry.ID, "timeEntryID")
		assert.Strings(t, timeEntry.Description, "description")
		assert.Strings(t, timeEntry.ProjectID, "projectID")
		assert.Strings(t, timeEntry.Start.String(), start.String())
		assert.Strings(t, timeEntry.End.String(), end.String())
		assert.Strings(t, timeEntry.Tags["tagName1"].ID, "tagId1")
		assert.Strings(t, timeEntry.Tags["tagName2"].ID, "tagId2")
	})

	t.Run("Get error on update timeEntry", func(t *testing.T) {
		fakeClient := &fakeClient{}
		fakeClient.updateTimeEntryErrorResponse()

		initClient = func(string) (clockifyApiClient, error) {
			return fakeClient, nil
		}

		apiClient, _ := NewClient("token")

		workspaceID := "ws1"
		timeEntry := TimeEntry{}

		_, err := apiClient.UpdateTimeEntry(workspaceID, timeEntry)

		assert.Errors(t, err, ErrClockifyTimeEntryUpdateFailed)
	})

	t.Run("Get error on incorrect tags after timeEntry update", func(t *testing.T) {
		fakeClient := &fakeClient{}
		fakeClient.updateTimeEntrySuccessResponse()

		initClient = func(string) (clockifyApiClient, error) {
			return fakeClient, nil
		}

		apiClient, _ := NewClient("token")
		workspaceID := "ws1"

		end := time.Date(1986, time.January, 5, 10, 46, 28, 0, &time.Location{})
		start := time.Date(1986, time.January, 1, 10, 46, 28, 0, &time.Location{})

		timeEntry := TimeEntry{
			ID:          "timeEntryID",
			Description: "description",
			ProjectID:   "projectID",
			Start:       start,
			End:         &end,
			Tags: map[string]Tag{
				"tagName1": {ID: "tagId1"},
				"tagName3": {ID: "tagId3"},
			},
		}

		_, err := apiClient.UpdateTimeEntry(workspaceID, timeEntry)

		assert.Errors(t, err, ErrClockifyTimeEntryTagsIncorrect)
	})

	t.Run("Get error on incorrect number of tags after timeEntry update", func(t *testing.T) {
		fakeClient := &fakeClient{}
		fakeClient.updateTimeEntrySuccessResponse()

		initClient = func(string) (clockifyApiClient, error) {
			return fakeClient, nil
		}

		apiClient, _ := NewClient("token")
		workspaceID := "ws1"

		end := time.Date(1986, time.January, 5, 10, 46, 28, 0, &time.Location{})
		start := time.Date(1986, time.January, 1, 10, 46, 28, 0, &time.Location{})

		timeEntry := TimeEntry{
			ID:          "timeEntryID",
			Description: "description",
			ProjectID:   "projectID",
			Start:       start,
			End:         &end,
			Tags: map[string]Tag{
				"tagName1": {ID: "tagId1"},
			},
		}

		_, err := apiClient.UpdateTimeEntry(workspaceID, timeEntry)

		assert.Errors(t, err, ErrClockifyInaccurateNumberOfTags)
	})
}

func TestGetLogRangeParameters(t *testing.T) {

	fakeClient := &fakeClient{}
	fakeClient.getMeSuccessResponse()

	initClient = func(string) (clockifyApiClient, error) {
		return fakeClient, nil
	}

	apiClient, _ := NewClient("token")

	format := "2006-01-02 15:04:05"

	workspaceID := "ws1"

	now := time.Date(1986, time.January, 5, 10, 46, 28, 0, &time.Location{})

	expectedStart := time.Date(1986, time.January, 1, 10, 46, 28, 0, &time.Location{})

	parameters, err := apiClient.getLongRangeParameters(expectedStart, now, workspaceID)
	assert.Errors(t, err, nil)
	assert.Strings(t, parameters.Workspace, workspaceID)
	assert.Strings(t, parameters.UserID, "userId")
	assert.Strings(t, parameters.FirstDate.Format(format), expectedStart.Format(format))
	assert.Strings(t, parameters.LastDate.Format(format), now.Format(format))
}

func TestTagsManagement(t *testing.T) {

	t.Run("Get tags ID list", func(t *testing.T) {

		timeEntry := TimeEntry{
			Tags: map[string]Tag{
				"tagName1": {ID: "tagId1"},
				"tagName2": {ID: "tagId2"},
			},
		}

		idsList := timeEntry.GetTagIDsList()

		assert.StringSlices(t, idsList, []string{"tagId1", "tagId2"})
	})

	t.Run("Get tags name list", func(t *testing.T) {

		timeEntry := TimeEntry{
			Tags: map[string]Tag{
				"tagName1": {ID: "tagId1", Name: "tagName1"},
				"tagName2": {ID: "tagId2", Name: "tagName2"},
			},
		}

		idsList := timeEntry.GetTagNamesList()

		assert.StringSlices(t, idsList, []string{"tagName1", "tagName2"})
	})

	t.Run("Add tag to timeentry", func(t *testing.T) {

		timeEntry := TimeEntry{
			Tags: map[string]Tag{},
		}

		tag := Tag{
			ID:   "tagId",
			Name: "tagName",
		}
		tags := timeEntry.Tags

		assert.Ints(t, len(tags), 0)

		timeEntry.AddTag(tag)

		assert.Strings(t, tags["tagName"].ID, "tagId")
		assert.Ints(t, len(tags), 1)
	})

	t.Run("Remove tag from timeentry", func(t *testing.T) {

		tag1 := Tag{
			ID:   "tagId1",
			Name: "tagName1",
		}

		tag2 := Tag{
			ID:   "tagId2",
			Name: "tagName2",
		}

		timeEntry := TimeEntry{
			Tags: map[string]Tag{tag1.Name: tag1, tag2.Name: tag2},
		}

		tags := timeEntry.Tags

		assert.Ints(t, len(tags), 2)
		assert.Strings(t, tags["tagName1"].ID, "tagId1")
		assert.Strings(t, tags["tagName2"].ID, "tagId2")

		timeEntry.RemoveTag("tagName1")

		assert.Ints(t, len(tags), 1)
		assert.Strings(t, tags["tagName2"].ID, "tagId2")
	})

	t.Run("Check if timeentry has tags - names", func(t *testing.T) {

		tag1 := Tag{
			ID:   "tagId1",
			Name: "tagName1",
		}

		tag2 := Tag{
			ID:   "tagId2",
			Name: "tagName2",
		}

		tag3 := Tag{
			ID:   "tagId3",
			Name: "tagName3",
		}

		timeEntry := TimeEntry{
			Tags: map[string]Tag{
				tag1.Name: tag1,
				tag2.Name: tag2,
				tag3.Name: tag3,
			},
		}

		isTagged := timeEntry.IsTaggedWith("tagName1")

		assert.Bools(t, isTagged, true)

		isTagged = timeEntry.IsTaggedWith("tagName1", "tagName2")

		assert.Bools(t, isTagged, true)

		isTagged = timeEntry.IsTaggedWith("randomTag")

		assert.Bools(t, isTagged, false)

		isTagged = timeEntry.IsTaggedWith("tagName1", "randomTag")

		assert.Bools(t, isTagged, false)
	})

	t.Run("Check if timeentry has tags - ids", func(t *testing.T) {
		tag1 := Tag{
			ID:   "tagId1",
			Name: "tagName1",
		}

		tag2 := Tag{
			ID:   "tagId2",
			Name: "tagName2",
		}

		timeEntry := TimeEntry{
			Tags: map[string]Tag{
				tag1.Name: tag1,
				tag2.Name: tag2,
			},
		}

		hasTags := timeEntry.hasTagsWithIds([]string{"tagId1", "tagId2"})

		assert.Bools(t, hasTags, true)

		hasTags = timeEntry.hasTagsWithIds([]string{"tagId1", "randomTag"})

		assert.Bools(t, hasTags, false)

	})
}

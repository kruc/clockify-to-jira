package clockify

import (
	"errors"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
)

type fakeClient struct {
	getMeResponse           func() (dto.User, error)
	logRangeResponse        func() ([]dto.TimeEntry, error)
	getTagsResponse         func() ([]dto.Tag, error)
	updateTimeEntryResponse func() (dto.TimeEntryImpl, error)
}

func (f *fakeClient) getTagsSuccessResponse() {
	f.getTagsResponse = func() ([]dto.Tag, error) {
		tags := []dto.Tag{
			{ID: "tagId1", Name: "tag1", WorkspaceID: "worskapce_id"},
			{ID: "tagId2", Name: "tag2", WorkspaceID: "worskapce_id"},
		}

		return tags, nil
	}
}

func (f *fakeClient) getTagsErrorResponse() {
	f.getTagsResponse = func() ([]dto.Tag, error) {
		return nil, errors.New("random-error")
	}
}

func (f *fakeClient) GetTags(_ api.GetTagsParam) ([]dto.Tag, error) {
	return f.getTagsResponse()
}

func (f *fakeClient) logRangeSuccessResponse() {
	f.logRangeResponse = func() ([]dto.TimeEntry, error) {
		timeEntryStart1 := time.Date(2025, time.January, 8, 10, 30, 0, 0, time.UTC)
		timeEntryEnd1 := time.Date(2025, time.January, 8, 17, 0, 0, 0, time.UTC)

		timeEntryStart2 := time.Date(2025, time.January, 9, 8, 00, 0, 0, time.UTC)
		timeEntryEnd2 := time.Date(2025, time.January, 9, 16, 15, 0, 0, time.UTC)

		timeEntry1 := dto.TimeEntry{
			ID:          "id1",
			WorkspaceID: "ws1",
			Description: "XYZ-123 Timentry description",
			TimeInterval: dto.TimeInterval{
				Start:    timeEntryStart1,
				End:      &timeEntryEnd1,
				Duration: "PT7H30M",
			},
			Project: &dto.Project{
				ID:         "projectID1",
				Name:       "projectName1",
				ClientName: "clientName1",
			},
			Tags: []dto.Tag{
				{ID: "tagId1", Name: "tagName1", WorkspaceID: "ws1"},
				{ID: "tagId2", Name: "tagName2", WorkspaceID: "ws1"},
			},
		}
		timeEntry2 := dto.TimeEntry{
			ID:          "id2",
			WorkspaceID: "ws1",
			Description: "ABC-123 Timentry description",
			TimeInterval: dto.TimeInterval{
				Start:    timeEntryStart2,
				End:      &timeEntryEnd2,
				Duration: "PT8H15M",
			},
			Project: &dto.Project{
				ID:         "projectID2",
				Name:       "projectName2",
				ClientName: "clientName2",
			},
			Tags: []dto.Tag{
				{ID: "tagId3", Name: "tagName3", WorkspaceID: "ws1"},
				{ID: "tagId4", Name: "tagName4", WorkspaceID: "ws1"},
				{ID: "tagId5", Name: "tagName5", WorkspaceID: "ws1"},
			},
		}
		return []dto.TimeEntry{timeEntry1, timeEntry2}, nil
	}
}

func (f *fakeClient) logRangeErrorResponse() {
	f.logRangeResponse = func() ([]dto.TimeEntry, error) {
		return nil, errors.New("random-error")
	}
}

func (f *fakeClient) LogRange(_ api.LogRangeParam) ([]dto.TimeEntry, error) {
	return f.logRangeResponse()
}

func (f *fakeClient) getMeSuccessResponse() {
	f.getMeResponse = func() (dto.User, error) {
		return dto.User{ID: "userId"}, nil
	}
}

func (f *fakeClient) getMeErrorResponse() {
	f.getMeResponse = func() (dto.User, error) {
		return dto.User{}, errors.New("random error")
	}
}

func (f *fakeClient) GetMe() (dto.User, error) {
	return f.getMeResponse()
}

func (f *fakeClient) updateTimeEntrySuccessResponse() {
	f.updateTimeEntryResponse = func() (dto.TimeEntryImpl, error) {
		end := time.Date(1986, time.January, 5, 10, 46, 28, 0, &time.Location{})
		start := time.Date(1986, time.January, 1, 10, 46, 28, 0, &time.Location{})

		timeEntry := dto.TimeEntryImpl{
			ID:           "timeEntryID",
			Description:  "description",
			ProjectID:    "projectID",
			TimeInterval: dto.TimeInterval{Start: start, End: &end},
			TagIDs:       []string{"tagId1", "tagId2"},
		}

		return timeEntry, nil
	}
}

func (f *fakeClient) updateTimeEntryErrorResponse() {
	f.updateTimeEntryResponse = func() (dto.TimeEntryImpl, error) {
		return dto.TimeEntryImpl{}, errors.New("random error")
	}
}

func (f *fakeClient) UpdateTimeEntry(api.UpdateTimeEntryParam) (dto.TimeEntryImpl, error) {
	return f.updateTimeEntryResponse()
}

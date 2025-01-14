package clockify

import (
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
)

const (
	ErrClockifyClientInitError             = ClockifyErr("Clockify client init error - check your token")
	ErrClockifyFailToFetchLoggedInUserData = ClockifyErr("Cannot get logged in user data")
	ErrClockifyFailToFetchTimeEntries      = ClockifyErr("Cannot fetch timeentries")
	ErrClockifyFailToFetchWorkspaceTags    = ClockifyErr("Cannot fetch workspace tags")
	ErrClockifyTimeEntryUpdateFailed       = ClockifyErr("Cannot update time entry")
	ErrClockifyTimeEntryTagsIncorrect      = ClockifyErr("Incorrect tags after timentry update")
	ErrClockifyInaccurateNumberOfTags      = ClockifyErr("Inaccurate number of tags after timeentry update")
)

type ClockifyErr string

func (e ClockifyErr) Error() string {
	return string(e)
}

type clockifyApiClient interface {
	LogRange(_ api.LogRangeParam) ([]dto.TimeEntry, error)
	GetTags(api.GetTagsParam) ([]dto.Tag, error)
	GetMe() (dto.User, error)
	UpdateTimeEntry(api.UpdateTimeEntryParam) (dto.TimeEntryImpl, error)
}

type ApiClient struct {
	client clockifyApiClient
}

var initClient = func(token string) (clockifyApiClient, error) {

	return api.NewClient(token)
}

func NewClient(token string) (*ApiClient, error) {

	clockifyClient, err := initClient(token)

	if err != nil {
		return nil, ErrClockifyClientInitError
	}

	clockifyApiClient := &ApiClient{
		client: clockifyClient,
	}

	return clockifyApiClient, nil
}

func (c *ApiClient) GetWorkspaceTags(worskapceId string) (map[string]Tag, error) {

	params := api.GetTagsParam{Workspace: worskapceId}
	tags, err := c.client.GetTags(params)

	if err != nil {
		return nil, ErrClockifyFailToFetchWorkspaceTags
	}

	tagsMap := make(map[string]Tag, len(tags))

	for _, tag := range tags {
		tagsMap[tag.Name] = Tag(tag)
	}

	return tagsMap, nil
}

func (c *ApiClient) GetTimeEntriesFromGivenPeriod(start, end time.Time, workspaceId string) ([]TimeEntry, error) {

	logRangeParam, err := c.getLongRangeParameters(start, end, workspaceId)

	if err != nil {
		return nil, ErrClockifyFailToFetchLoggedInUserData
	}

	timeEntries, err := c.client.LogRange(logRangeParam)

	if err != nil {
		return nil, ErrClockifyFailToFetchTimeEntries
	}

	result := mapTimeEntries(timeEntries)

	return result, nil
}

func (c *ApiClient) getLongRangeParameters(start, end time.Time, workspaceID string) (api.LogRangeParam, error) {

	userId, err := c.client.GetMe()

	if err != nil {
		return api.LogRangeParam{}, err
	}

	parameters := api.LogRangeParam{
		UserID:    userId.ID,
		Workspace: workspaceID,
		FirstDate: start,
		LastDate:  end,
	}

	return parameters, nil
}

func (c *ApiClient) UpdateTimeEntry(workspaceId string, timeEntry TimeEntry) (TimeEntry, error) {

	updateParams := api.UpdateTimeEntryParam{
		Workspace:   workspaceId,
		TimeEntryID: timeEntry.ID,
		Description: timeEntry.Description,
		ProjectID:   timeEntry.ProjectID,
		Start:       timeEntry.Start,
		End:         timeEntry.End,
		TagIDs:      timeEntry.GetTagIDsList(),
	}

	updatedTimeEntry, err := c.client.UpdateTimeEntry(updateParams)

	if err != nil {
		return TimeEntry{}, ErrClockifyTimeEntryUpdateFailed
	}

	if len(timeEntry.Tags) != len(updatedTimeEntry.TagIDs) {
		return TimeEntry{}, ErrClockifyInaccurateNumberOfTags
	}

	if !timeEntry.hasTagsWithIds(updatedTimeEntry.TagIDs) {
		return TimeEntry{}, ErrClockifyTimeEntryTagsIncorrect
	}

	return TimeEntry{
		ID:          updatedTimeEntry.ID,
		Description: updatedTimeEntry.Description,
		ProjectID:   updatedTimeEntry.ProjectID,
		Start:       updatedTimeEntry.TimeInterval.Start,
		End:         updatedTimeEntry.TimeInterval.End,
		Tags:        timeEntry.Tags,
	}, nil
}

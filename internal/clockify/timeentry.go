package clockify

import (
	"slices"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

type Tag dto.Tag

type TimeEntry struct {
	ID          string
	Description string
	ClientName  string
	ProjectID   string
	ProjectName string
	Start       time.Time
	End         *time.Time
	Duration    string
	Tags        map[string]Tag
}

func (te *TimeEntry) GetTagIDsList() []string {

	idsList := []string{}

	for _, tag := range te.Tags {
		idsList = append(idsList, tag.ID)
	}

	return idsList
}

func (te *TimeEntry) GetTagNamesList() []string {

	namesList := []string{}

	for _, tag := range te.Tags {
		namesList = append(namesList, tag.Name)
	}

	return namesList
}

func (te *TimeEntry) AddTag(tag Tag) {

	te.Tags[tag.Name] = tag
}

func (te *TimeEntry) RemoveTag(tagName string) {

	delete(te.Tags, tagName)
}

func (te *TimeEntry) IsTaggedWith(tagNames ...string) bool {

	for _, tagName := range tagNames {
		_, ok := te.Tags[tagName]

		if !ok {
			return false
		}
	}

	return true
}

func (te *TimeEntry) hasTagsWithIds(idsList []string) bool {

	for _, tag := range te.Tags {
		ok := slices.Contains(idsList, tag.ID)

		if !ok {
			return false
		}
	}

	return true
}

func mapTimeEntries(timeEntries []dto.TimeEntry) []TimeEntry {

	result := make([]TimeEntry, len(timeEntries))

	for key, timeEntry := range timeEntries {
		result[key].ID = timeEntry.ID
		result[key].Description = timeEntry.Description
		result[key].ClientName = timeEntry.Project.ClientName
		result[key].ProjectID = timeEntry.Project.ID
		result[key].ProjectName = timeEntry.Project.Name
		result[key].Start = timeEntry.TimeInterval.Start
		result[key].End = timeEntry.TimeInterval.End
		result[key].Duration = timeEntry.TimeInterval.Duration
		result[key].Tags = parseTags(timeEntry.Tags)
	}

	return result
}

func parseTags(tags []dto.Tag) map[string]Tag {

	result := map[string]Tag{}

	for _, tag := range tags {
		result[tag.Name] = Tag(tag)
	}

	return result
}

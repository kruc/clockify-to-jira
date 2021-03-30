package main

import (
	"fmt"
	s "strings"
	"time"

	"github.com/kruc/clockify-api/gctag"
)

func dosko(timeSpentSeconds, stachurskyMode int) (int, doskoDebugInfo) {

	d, err := time.ParseDuration(fmt.Sprintf("%vs", timeSpentSeconds))
	if err != nil {
		panic(err)
	}

	stachurskyFactor := time.Duration(stachurskyMode) * time.Minute
	roundedValue := d.Round(stachurskyFactor)

	if int(roundedValue.Seconds()) == 0 {
		roundedValue = stachurskyFactor
	}

	doskoDebugInfo := doskoDebugInfo{
		originalTime: d.String(),
		doskoTime:    roundedValue.String(),
	}

	return int(roundedValue.Seconds()), doskoDebugInfo
}

func removeTag(tagsList []string, tagToRemove string) []string {
	for i := 0; i < len(tagsList); i++ {
		if tagsList[i] == tagToRemove {
			tagsList = append(tagsList[:i], tagsList[i+1:]...)
			i-- // form the remove item index to start iterate next item
		}
	}
	return tagsList
}

func displayTagsName(tags []gctag.Tag) []string {
	tagsName := []string{}
	for _, tag := range tags {
		tagsName = append(tagsName, tag.Name)
	}

	return tagsName
}

func adjustClockifyDate(clockifyDate time.Time) time.Time {
	clockifyDate = clockifyDate.Add(time.Millisecond * 1)

	return clockifyDate
}

func parseIssueID(value string) string {
	fields := s.Fields(value)

	return trimBrackets(fields[0])
}

func trimBrackets(issueID string) string {
	trimmedissueID := s.TrimPrefix(issueID, "[")
	trimmedissueID = s.TrimSuffix(trimmedissueID, ":")
	trimmedissueID = s.TrimSuffix(trimmedissueID, "]")

	return trimmedissueID
}

func parseIssueComment(value string) string {
	fields := s.Fields(value)

	return s.Join(fields[1:], " ")
}

func getTimeDiff(start, stop time.Time) int {
	return int(stop.Sub(start).Seconds())
}

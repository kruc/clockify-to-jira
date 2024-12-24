package main

import (
	"fmt"
	s "strings"
	"time"
)

func dosko(timeSpentSeconds, stachurskyMode int) (int, string, string) {

	d, err := time.ParseDuration(fmt.Sprintf("%vs", timeSpentSeconds))
	if err != nil {
		panic(err)
	}

	stachurskyFactor := time.Duration(stachurskyMode) * time.Minute
	roundedValue := d.Round(stachurskyFactor)

	if int(roundedValue.Seconds()) == 0 {
		roundedValue = stachurskyFactor
	}

	return int(roundedValue.Seconds()), d.String(), roundedValue.String()
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

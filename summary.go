package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type summary struct {
	start          time.Time
	end            time.Time
	entriesCount   int
	totalTime      int
	totalDoskoTime int
	doskoFactor    int
}

func (s *summary) increaseTimeEntryCount() {
	s.entriesCount++
}

func (s *summary) addTimeEntryDuration(timeEntryDuration int) {
	s.totalTime += timeEntryDuration
}

func (s *summary) addDoskoTimeEntryDuration(timeEntryDuration int) {
	s.totalDoskoTime += timeEntryDuration
}

func (s *summary) addDoskoFactor(doskoFactor int) {
	s.doskoFactor = doskoFactor
}

func (s *summary) getTotalTime() string {
	parsedDuration, err := time.ParseDuration(fmt.Sprintf("%ds", s.totalTime))

	if err != nil {
		panic(err)
	}

	return parsedDuration.String()
}

func (s *summary) getTotalDoskoTime() string {
	parsedDuration, err := time.ParseDuration(fmt.Sprintf("%ds", s.totalDoskoTime))

	if err != nil {
		panic(err)
	}

	return parsedDuration.String()
}

func (s *summary) show() {
	timeFormat := "2006-01-02 15:04:05"
	log.Infof("\n")
	log.Infof("---------")
	log.Infof("Summary:")
	log.Infof("---------")
	log.Infof("Time entries range: %v - %v\n", s.start.Format(timeFormat), s.end.Format(timeFormat))
	log.Infof("Number of time entries: %+v\n", s.entriesCount)
	log.Infof("Total time: %+v\n", s.getTotalTime())
	log.Infof("Total dosko: %+v (t=%vm)\n", s.getTotalDoskoTime(), s.doskoFactor)
}

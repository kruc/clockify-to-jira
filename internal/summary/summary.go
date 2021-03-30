package summary

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type Details struct {
	Start          time.Time
	End            time.Time
	entriesCount   int
	totalTime      int
	totalDoskoTime int
	doskoFactor    int
	TimeFormat     string
}

func (s *Details) IncreaseTimeEntryCount() {
	s.entriesCount++
}

func (s *Details) AddTimeEntryDuration(timeEntryDuration int) {
	s.totalTime += timeEntryDuration
}

func (s *Details) AddDoskoTimeEntryDuration(timeEntryDuration int) {
	s.totalDoskoTime += timeEntryDuration
}

func (s *Details) AddDoskoFactor(doskoFactor int) {
	s.doskoFactor = doskoFactor
}

func (s *Details) getTotalTime() string {
	parsedDuration, err := time.ParseDuration(fmt.Sprintf("%ds", s.totalTime))

	if err != nil {
		panic(err)
	}

	return parsedDuration.String()
}

func (s *Details) getTotalDoskoTime() string {
	parsedDuration, err := time.ParseDuration(fmt.Sprintf("%ds", s.totalDoskoTime))

	if err != nil {
		panic(err)
	}

	return parsedDuration.String()
}

func (s *Details) Show() {
	log.Infof("\n")
	log.Infof("---------")
	log.Infof("Summary:")
	log.Infof("---------")
	log.Infof("Time entries range: %v - %v\n", s.Start.Format(s.TimeFormat), s.End.Format(s.TimeFormat))
	log.Infof("Number of time entries: %+v\n", s.entriesCount)
	log.Infof("Total time: %+v\n", s.getTotalTime())
	log.Infof("Total dosko: %+v (t=%vm)\n", s.getTotalDoskoTime(), s.doskoFactor)
}

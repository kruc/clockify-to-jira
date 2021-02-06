package main

import (
	"fmt"
	"os"
	s "strings"
	"time"

	"github.com/andygrunwald/go-jira"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	clockifyapi "github.com/kruc/clockify-api"
	"github.com/kruc/clockify-api/gctimeentry"
)

type clockifyData struct {
	client           string
	project          string
	issueID          string
	issueComment     string
	started          time.Time
	timeSpentSeconds int
}

type doskoDebugInfo struct {
	originalTime string
	doskoTime    string
}

type clockifyTags map[string]string

var (
	globalConfig   globalConfigType
	config         = "config"
	configPath     string
	logFile        *os.File
	debug          bool
	applyMode      bool
	clientSelector string
	version        bool
	// BuildVersion info
	BuildVersion string
	// BuildDate info
	BuildDate string
	// GitCommit info
	GitCommit string
)

const timeFormat = "2006-01-02 15:04:05"

func init() {

	if !checkConfiguration() {
		os.Exit(1)
	}

	globalConfig = parseGlobalConfig()

	flag.BoolVar(&applyMode, "apply", false, "Update jira tasks workloads")
	flag.IntVarP(&globalConfig.period, "period", "p", globalConfig.period, "Migrate time entries from last given days")
	flag.StringVarP(&globalConfig.logFormat, "format", "f", globalConfig.logFormat, "Log format (text|json)")
	flag.StringVarP(&globalConfig.logOutput, "output", "o", globalConfig.logOutput, "Log output (stdout|filename)")
	flag.StringVarP(&globalConfig.workspaceID, "workspace", "w", globalConfig.workspaceID, "Clockify workspace id")
	flag.IntVarP(&globalConfig.defaultClient.stachurskyMode, "tryb-niepokorny", "t", globalConfig.defaultClient.stachurskyMode, "Rounding up the value of logged time up (minutes)")
	flag.BoolVarP(&version, "version", "v", false, "Display version")
	flag.StringVarP(&clientSelector, "client", "c", "", "Migrate time entries from given client")
	flag.Parse()

	// Prepare logger
	configureLogger()
}

func main() {
	defer logFile.Close()
	if version {
		displayVersion()
		return
	}

	clockifyClient, err := clockifyapi.NewClient(viper.GetString("clockify_token"))
	timeEntryClient := clockifyClient.TimeEntryClient
	userClient := clockifyClient.UserClient
	tagClient := clockifyClient.TagClient

	currentUser, err := userClient.GetCurrentlyLoggedInUser()
	if err != nil {
		panic(err)
	}

	tags, err := tagClient.GetTags(globalConfig.workspaceID)
	if err != nil {
		panic(err)
	}

	clockifyTags := tags.ToMap()

	end := time.Now()
	start := end.Add(time.Hour * 24 * time.Duration(globalConfig.period) * -1)

	queryParameters := gctimeentry.QueryParameters{
		Start:    start,
		End:      end,
		Hydrated: true,
		PageSize: 150,
	}

	timeEntries, err := timeEntryClient.GetRange(queryParameters, globalConfig.workspaceID, currentUser.ID)

	if err != nil {
		log.Error(err)
		return
	}

	if clientSelector != "" {
		mapedTimeEntries := timeEntries.ToMap()
		timeEntries = mapedTimeEntries[s.ToLower(clientSelector)]
	}

	summary := summary{start: start, end: end}
	for _, timeEntry := range timeEntries {
		if timeEntry.IsTagged(globalConfig.jiraMigrationSuccessTag) || timeEntry.IsTagged(globalConfig.jiraMigrationSkipTag) {
			continue
		}
		log.Infof("\n")
		log.Info(fmt.Sprintf("Worklog: %v", timeEntry.Description))

		if timeEntry.ProjectID == "" {
			log.WithFields(log.Fields{
				"timeEntry": timeEntry,
				"reason":    "Probably time entry is not assign to project in Clockify",
				"solution":  "Edit time entry in clockify and assign it to project",
			}).Error(err)
			continue
		}

		clientConfigPath := fmt.Sprintf("client.%v", s.ToLower(timeEntry.Project.ClientName))

		if !viper.IsSet(clientConfigPath) {
			generateClientConfigTemplate(clientConfigPath)
			continue
		}

		if !viper.GetBool(fmt.Sprintf("%v.%v", clientConfigPath, "enabled")) {
			log.Warnf("Don't forget to enable client (set %v.enabled = true)", clientConfigPath)
			continue
		}

		clientConfig := parseClientConfig(clientConfigPath, globalConfig)
		timeDiff := getTimeDiff(timeEntry.TimeInterval.Start, timeEntry.TimeInterval.End)
		timeSpentSeconds, doskoDebugInfo := dosko(timeDiff, clientConfig.stachurskyMode)

		summary.increaseTimeEntryCount()
		summary.addTimeEntryDuration(timeDiff)
		summary.addDoskoTimeEntryDuration(timeSpentSeconds)
		summary.addDoskoFactor(clientConfig.stachurskyMode)

		clockifyData := clockifyData{
			client:           s.ToLower(timeEntry.Project.ClientName),
			project:          s.ToLower(timeEntry.Project.Name),
			issueID:          parseIssueID(timeEntry.Description),
			issueComment:     parseIssueComment(timeEntry.Description),
			started:          adjustClockifyDate(timeEntry.TimeInterval.Start),
			timeSpentSeconds: timeSpentSeconds,
		}

		// JIRA PART
		tp := jira.BasicAuthTransport{
			Username: clientConfig.jiraUsername,
			Password: clientConfig.jiraPassword,
		}

		jiraClient, _ := jira.NewClient(tp.Client(), clientConfig.jiraHost)

		tt := jira.Time(clockifyData.started)
		worklogRecord := jira.WorklogRecord{
			Comment:          clockifyData.issueComment,
			TimeSpentSeconds: clockifyData.timeSpentSeconds,
			Started:          &tt,
		}
		log.Infof("-------")
		log.Infof("Client: %+v\n", clockifyData.client)
		log.Infof("Project: %+v\n", clockifyData.project)
		log.Infof("Date: %+v\n", clockifyData.started.Format(timeFormat))
		log.Infof("Time spent: %+v (clockify: %+v stachurskyMode: %+vm)\n", doskoDebugInfo.doskoTime, doskoDebugInfo.originalTime, clientConfig.stachurskyMode)
		log.Infof("Comment: %+v\n", worklogRecord.Comment)

		if applyMode == true {

			jwr, jr, err := jiraClient.Issue.AddWorklogRecord(clockifyData.issueID, &worklogRecord)

			if err != nil {
				log.WithFields(log.Fields{
					"worklogRecord": jwr,
					"response":      jr,
				}).Error(err)

				timeEntry.AddTag(clockifyTags[globalConfig.jiraMigrationFailedTag].ID)
				log.Info(fmt.Sprintf("Add %v tag", globalConfig.jiraMigrationFailedTag))
			} else {
				log.Info(fmt.Sprintf("Jira workload added"))
				timeEntry.RemoveTag(clockifyTags[globalConfig.jiraMigrationFailedTag].ID)
				timeEntry.AddTag(clockifyTags[globalConfig.jiraMigrationSuccessTag].ID)
				log.Info(fmt.Sprintf("Add %v tag", globalConfig.jiraMigrationSuccessTag))
			}
			// TODO: Create timeEntry update struct
			timeEntry.Start = timeEntry.TimeInterval.Start
			timeEntry.End = timeEntry.TimeInterval.End
			te, err := timeEntryClient.Update(globalConfig.workspaceID, timeEntry.ID, &timeEntry)

			if err != nil {
				log.WithFields(log.Fields{
					"timeEntry": te,
				}).Error(err)
			}
			issueURL := fmt.Sprintf("%v/browse/%v?focusedWorklogId=%s", clientConfig.jiraHost, clockifyData.issueID, "123")
			log.Infof("Issue url: %v\n", issueURL)
			log.Info(fmt.Sprintf("Finish processing %v: %v", timeEntry.ID, timeEntry.Description))
		}

	}
	summary.show()
}

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

func displayVersion() {
	fmt.Printf("BuildVersion: %s\tBuildDate: %s\tGitCommit: %s\n", BuildVersion, BuildDate, GitCommit)
}

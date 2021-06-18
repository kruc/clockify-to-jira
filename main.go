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
	"github.com/kruc/clockify-to-jira/internal/config"
	"github.com/kruc/clockify-to-jira/internal/summary"
	"github.com/kruc/clockify-to-jira/internal/version"
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
	logFile        *os.File
	debugMode      bool
	applyMode      bool
	clientSelector string
	versionFlag    bool
)

const timeFormat = "2006-01-02 15:04:05"

func init() {

	flag.BoolVar(&applyMode, "apply", false, "Update jira tasks workloads")
	flag.BoolVar(&debugMode, "debug", false, "Include already logged time entries")
	flag.StringVarP(&clientSelector, "client", "c", "", "Migrate time entries from given client")

	flag.Parse()

	if applyMode && debugMode {
		log.Warning("Debug and apply flag cannot be set together! Ignoring debug")
		debugMode = false
	}
}

func main() {

	defer logFile.Close()
	if version.VersionFlag {
		version.DisplayVersion()
		return
	}

	globalConfig := config.GetGlobalConfig()

	configureLogger(globalConfig.LogFormat, globalConfig.LogOutput)

	clockifyClient, err := clockifyapi.NewClient(viper.GetString("clockify_token"))

	if err != nil {
		panic(err.Error())
	}
	timeEntryClient := clockifyClient.TimeEntryClient
	userClient := clockifyClient.UserClient
	tagClient := clockifyClient.TagClient

	currentUser, err := userClient.GetCurrentlyLoggedInUser()
	if err != nil {
		panic(err)
	}

	tags, err := tagClient.GetTags(globalConfig.WorkspaceID)
	if err != nil {
		panic(err)
	}

	clockifyTags := tags.ToMap()

	end := time.Now()
	start := end.Add(time.Hour * 24 * time.Duration(globalConfig.Period) * -1)

	queryParameters := gctimeentry.QueryParameters{
		Start:    start,
		End:      end,
		Hydrated: true,
		PageSize: 150,
	}

	timeEntries, err := timeEntryClient.GetRange(queryParameters, globalConfig.WorkspaceID, currentUser.ID)

	if err != nil {
		log.Error(err)
		return
	}

	if clientSelector != "" {
		mapedTimeEntries := timeEntries.ToMap()
		timeEntries = mapedTimeEntries[s.ToLower(clientSelector)]
	}
	summary := summary.Details{Start: start, End: end, TimeFormat: timeFormat}

	for _, timeEntry := range timeEntries {
		if (timeEntry.IsTagged(globalConfig.JiraMigrationSuccessTag) ||
			timeEntry.IsTagged(globalConfig.JiraMigrationSkipTag) ||
			timeEntry.TimeInterval.Duration == "") &&
			!debugMode {
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
			config.GenerateClientConfigTemplate(clientConfigPath)
			continue
		}

		if !viper.GetBool(fmt.Sprintf("%v.%v", clientConfigPath, "enabled")) {
			log.Warnf("Don't forget to enable client (set %v.enabled = true)", clientConfigPath)
			continue
		}

		clientConfig := config.ParseClientConfig(clientConfigPath, globalConfig)
		timeDiff := getTimeDiff(timeEntry.TimeInterval.Start, timeEntry.TimeInterval.End)
		timeSpentSeconds, doskoDebugInfo := dosko(timeDiff, clientConfig.StachurskyMode)

		summary.IncreaseTimeEntryCount()
		summary.AddTimeEntryDuration(timeDiff)
		summary.AddDoskoTimeEntryDuration(timeSpentSeconds)
		summary.AddDoskoFactor(clientConfig.StachurskyMode)

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
			Username: clientConfig.JiraUsername,
			Password: clientConfig.JiraPassword,
		}

		jiraClient, _ := jira.NewClient(tp.Client(), clientConfig.JiraHost)

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
		log.Infof("Time spent: %+v (clockify: %+v stachurskyMode: %+vm)\n", doskoDebugInfo.doskoTime, doskoDebugInfo.originalTime, clientConfig.StachurskyMode)
		log.Infof("Comment: %+v\n", worklogRecord.Comment)
		log.Infof("Tags: %v", displayTagsName(timeEntry.Tags))

		if applyMode == true {

			jwr, jr, err := jiraClient.Issue.AddWorklogRecord(clockifyData.issueID, &worklogRecord)

			if err != nil {
				log.WithFields(log.Fields{
					"worklogRecord": jwr,
					"response":      jr,
				}).Error(err)

				timeEntry.AddTag(clockifyTags[globalConfig.JiraMigrationFailedTag].ID)
				log.Info(fmt.Sprintf("Add %v tag", globalConfig.JiraMigrationFailedTag))
			} else {
				log.Info(fmt.Sprintf("Jira workload added"))
				timeEntry.RemoveTag(clockifyTags[globalConfig.JiraMigrationFailedTag].ID)
				timeEntry.AddTag(clockifyTags[globalConfig.JiraMigrationSuccessTag].ID)
				log.Info(fmt.Sprintf("Add %v tag", globalConfig.JiraMigrationSuccessTag))
			}
			// TODO: Create timeEntry update struct
			timeEntry.Start = timeEntry.TimeInterval.Start
			timeEntry.End = timeEntry.TimeInterval.End
			te, err := timeEntryClient.Update(globalConfig.WorkspaceID, timeEntry.ID, &timeEntry)

			if err != nil {
				log.WithFields(log.Fields{
					"timeEntry": te,
				}).Error(err)
			}
			issueURL := fmt.Sprintf("%v/browse/%v?focusedWorklogId=%s", clientConfig.JiraHost, clockifyData.issueID, "123")
			log.Infof("Issue url: %v\n", issueURL)
			log.Info(fmt.Sprintf("Finish processing %v: %v", timeEntry.ID, timeEntry.Description))
		}

	}
	summary.Show()
}

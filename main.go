package main

import (
	"fmt"
	"os"
	"slices"
	s "strings"
	"time"

	"github.com/andygrunwald/go-jira"

	"github.com/kruc/clockify-to-jira/internal/clockify"
	"github.com/kruc/clockify-to-jira/internal/config"
	"github.com/kruc/clockify-to-jira/internal/flag"
	"github.com/kruc/clockify-to-jira/internal/logger"
	"github.com/kruc/clockify-to-jira/internal/outcome"
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

func main() {

	log := logger.InitializeLogger()
	flag, err := flag.InitializeFlags(os.Args)

	if err != nil {
		log.Error("Ops, something went wrong during flag initialization!",
			"error", err,
			"args", os.Args,
		)
		return
	}

	if flag.Help {
		flag.PrintDefaults()
		return
	}

	if flag.Version {
		version.ShowBuildDetails(os.Stdout)
		return
	}

	config, err := config.LoadFromYamlFile(flag.ConfigFilePath)

	if err != nil {
		log.Error("Ops, something went wrong while loading the configuration!",
			"error", err)
		return
	}

	if flag.Period != 0 {
		config.OverwritePeriodSetting(flag.Period)
	}

	if flag.Precision != 0 {
		config.OverwritePrecisionSetting(flag.Precision)
	}

	workspaces, err := config.FindWorkspaces(flag.Workspaces)

	if err != nil {
		log.Error("Ops, something went wrong during workspace listing!",
			"error", err)
		return
	}

	clockifyClient, err := clockify.NewClient(config.Global.ClockifyToken)

	if err != nil {
		log.Error("Ops, something went wrong during clockify client initialization!",
			"error", err)
	}

	ch := make(chan string)

	for workspaceKey, workspace := range workspaces {

		go func(chan string) {

			clockifyTags, err := clockifyClient.GetWorkspaceTags(workspace.WorkspaceId)

			if err != nil {
				log.Error("Ops, something went wrong during tags fetching!",
					"error", err)
			}

			now := time.Now()
			start, end := config.GetTimeInterval(&now)

			timeEntries, err := clockifyClient.GetTimeEntriesFromGivenPeriod(start, end, workspace.WorkspaceId)

			if err != nil {
				log.Error("Ops, something went wrong during time entries fetching!",
					"error", err)
				return
			}

			summaryData := outcome.SummaryData{Start: start, End: end, Workspace: workspaceKey}

			slices.Reverse(timeEntries)

			for _, timeEntry := range timeEntries {

				if (timeEntry.IsTaggedWith(workspace.JiraMigrationSuccessTag) ||
					timeEntry.IsTaggedWith(workspace.JiraMigrationSkipTag) ||
					timeEntry.Duration == "") &&
					!flag.Debug {

					continue
				}

				if timeEntry.ProjectID == "" {
					log.Error("Ops, project not assign to time entry!",
						"solution", "Edit time entry in clockify and assign it to project",
						"timeEntry", timeEntry.Description,
					)
					continue
				}

				clientConfigId := s.ToLower(timeEntry.ClientName)

				if len(flag.Clients) != 0 && !slices.Contains(flag.Clients, clientConfigId) {
					continue
				}

				clientConfig, err := workspace.GetClient(clientConfigId)

				if err != nil {
					log.Error("Ops, something went wrong during get client!",
						"error", err)
					continue
				}

				if !clientConfig.Enabled {
					log.Warn("Don't forget to enable client",
						"solution", fmt.Sprintf("set workspaces.%s.clients.%s.enabled to true", workspaceKey, clientConfigId),
					)
					continue
				}

				timeDiff := getTimeDiff(timeEntry.Start, *timeEntry.End)
				timeSpentSeconds, originalTime, roundedTime := dosko(timeDiff, clientConfig.StachurskyMode)

				summaryData.IncreaseTimeEntryCount()
				summaryData.AddTimeEntryDuration(timeDiff)
				summaryData.AddDoskoTimeEntryDuration(timeSpentSeconds)
				summaryData.AddDoskoFactor(clientConfig.StachurskyMode)

				// JIRA PART
				clockifyData := clockifyData{
					client:           s.ToLower(timeEntry.ClientName),
					project:          s.ToLower(timeEntry.ProjectName),
					issueID:          parseIssueID(timeEntry.Description),
					issueComment:     parseIssueComment(timeEntry.Description),
					started:          adjustClockifyDate(timeEntry.Start),
					timeSpentSeconds: timeSpentSeconds,
				}

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

				if flag.Apply {

					jwr, jr, err := jiraClient.Issue.AddWorklogRecord(clockifyData.issueID, &worklogRecord)

					if err != nil {
						log.Error("Ops, something went wrong during worklog record adding!",
							"error", err,
							"worklogRecord", jwr,
							"response", jr,
						)

						timeEntry.AddTag(clockifyTags[workspace.JiraMigrationFailedTag])
						log.Info(fmt.Sprintf("Add %v tag", workspace.JiraMigrationFailedTag))
					} else {
						log.Info("Jira workload added")
						timeEntry.RemoveTag(workspace.JiraMigrationFailedTag)
						timeEntry.AddTag(clockifyTags[workspace.JiraMigrationSuccessTag])
						log.Info(fmt.Sprintf("Add %v tag", workspace.JiraMigrationSuccessTag))
					}

					te, err := clockifyClient.UpdateTimeEntry(workspace.WorkspaceId, timeEntry)

					if err != nil {
						log.Error("Ops, something went wrong during time entry updating",
							"error", err,
							"timeEntry", te,
						)
					}

					issueURL := fmt.Sprintf("%v/browse/%v?focusedWorklogId=%s", clientConfig.JiraHost, clockifyData.issueID, "123")
					log.Info("Finish timentry processing",
						"Id", timeEntry.ID,
						"Description", timeEntry.Description,
						"IssueUrl", issueURL)
				}

				worklogData := outcome.WorklogData{
					Description: timeEntry.Description,
					Workspace:   workspaceKey,
					Client:      clockifyData.client,
					Project:     clockifyData.project,
					Date:        clockifyData.started,
					Comment:     worklogRecord.Comment,
					Tags:        timeEntry.GetTagNamesList(),
					TimeSpent: outcome.DoskoDetails{
						OriginalTime: originalTime,
						RoundedTime:  roundedTime,
						Precision:    clientConfig.StachurskyMode,
					},
				}

				worklog := worklogData.GetSummary()

				log.Info(worklog)
			}
			summary, err := summaryData.GetSummary()
			if err != nil {
				log.Error("Ops, something went wrong during fetching summary!",
					"error", err)
			}
			ch <- summary
		}(ch)
	}

	for i := 0; i < len(workspaces); i++ {
		log.Info(<-ch)
	}
}

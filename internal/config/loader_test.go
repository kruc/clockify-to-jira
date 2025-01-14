package config

import (
	"testing"

	"github.com/kruc/clockify-to-jira/internal/assert"
)

func TestLoadConfigFromYamlFile(t *testing.T) {

	t.Run("Load from yaml file", func(t *testing.T) {
		config, err := LoadFromYamlFile("./examples/config.yaml")
		assert.Errors(t, err, nil)

		global := config.Global
		assert.Strings(t, global.ClockifyToken, "clockify-token")
		assert.Ints(t, global.Period, 1)

		defaultClient := config.DefaultClient
		assert.Strings(t, defaultClient.JiraClientUser, "firstname.lastname")
		assert.Strings(t, defaultClient.JiraHost, "https://jira.atlassian.net")
		assert.Strings(t, defaultClient.JiraPassword, "jira-password")
		assert.Strings(t, defaultClient.JiraUsername, "firstname.lastname@domain.io")
		assert.Ints(t, defaultClient.StachurskyMode, 15)

		defaultWorkspace := config.DefaultWorkspace
		assert.Strings(t, defaultWorkspace.JiraMigrationFailedTag, "jira-migration-failed")
		assert.Strings(t, defaultWorkspace.JiraMigrationSkipTag, "jira-migration-skip")
		assert.Strings(t, defaultWorkspace.JiraMigrationSuccessTag, "logged")

		workspaces := config.Workspaces
		assert.Ints(t, len(workspaces), 2)

		ws1 := workspaces["ws_1"]
		assert.Strings(t, ws1.WorkspaceId, "ws-1")
		assert.Strings(t, ws1.JiraMigrationFailedTag, "jira-migration-failed")
		assert.Strings(t, ws1.JiraMigrationSkipTag, "jira-migration-skip")
		assert.Strings(t, ws1.JiraMigrationSuccessTag, "logged")
		assert.Ints(t, len(ws1.Clients), 2)

		ws1Client1 := ws1.Clients["client_1"]
		assert.Bools(t, ws1Client1.Enabled, true)
		assert.Strings(t, ws1Client1.JiraClientUser, "username")
		assert.Strings(t, ws1Client1.JiraHost, "https://domain.atlassian.net")
		assert.Strings(t, ws1Client1.JiraPassword, "jirapassword-client-1")
		assert.Strings(t, ws1Client1.JiraUsername, "username@domain.com")
		assert.Ints(t, ws1Client1.StachurskyMode, 30)

		ws1Client2 := ws1.Clients["client_2"]
		assert.Bools(t, ws1Client2.Enabled, false)
		assert.Strings(t, ws1Client2.JiraClientUser, "firstname.lastname")
		assert.Strings(t, ws1Client2.JiraHost, "https://jira.atlassian.net")
		assert.Strings(t, ws1Client2.JiraPassword, "jira-password")
		assert.Strings(t, ws1Client2.JiraUsername, "firstname.lastname@domain.io")
		assert.Ints(t, ws1Client2.StachurskyMode, 15)

		ws2 := workspaces["ws_2"]
		assert.Ints(t, len(ws2.Clients), 1)
		assert.Strings(t, ws2.WorkspaceId, "ws-2")
		assert.Strings(t, ws2.JiraMigrationFailedTag, "failed")
		assert.Strings(t, ws2.JiraMigrationSkipTag, "skipped")
		assert.Strings(t, ws2.JiraMigrationSuccessTag, "logged")

		ws2Client3 := ws2.Clients["client_3"]
		assert.Bools(t, ws2Client3.Enabled, false)
		assert.Strings(t, ws2Client3.JiraClientUser, "firstname.lastname")
		assert.Strings(t, ws2Client3.JiraHost, "https://jira.atlassian.net")
		assert.Strings(t, ws2Client3.JiraPassword, "jirapassword-client-3")
		assert.Strings(t, ws2Client3.JiraUsername, "username3@domain.com")
		assert.Ints(t, ws2Client3.StachurskyMode, 15)
	})

	t.Run("Throw error if file not found", func(t *testing.T) {
		_, err := LoadFromYamlFile("./invalid-path")

		assert.Errors(t, err, ErrLoaderConfigFileNotFound)
	})

	t.Run("Throw error on invalid config", func(t *testing.T) {
		_, err := LoadFromYamlFile("./examples/invalid-config.yaml")

		assert.Errors(t, err, ErrLoaderInvalidConfiguration)
	})
}

func TestLoadConfigDataFromInlineData(t *testing.T) {

	t.Run("Load from inline data", func(t *testing.T) {

		validInlineConfig := []byte(`global:
  clockify_token: inline-clockify-token
  period: 9
`)

		configFromBytes, _ := Load(validInlineConfig)

		global := configFromBytes.Global
		assert.Strings(t, global.ClockifyToken, "inline-clockify-token")
		assert.Ints(t, global.Period, 9)
	})

	t.Run("Throw error on invalid config", func(t *testing.T) {

		invalidInlineConfig := []byte(`
global:
  clockify_token: inline-clockify-token
  log_format: inline-text
   log_output: inline-stdout #invalid line
  period: 9
`)

		_, err := Load(invalidInlineConfig)

		assert.Errors(t, err, ErrLoaderInvalidConfiguration)
	})
}

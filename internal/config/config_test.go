package config

import (
	"testing"
	"time"

	"github.com/kruc/clockify-to-jira/internal/assert"
)

const (
	clockifyToken = "clockifyToken"
	period        = 1

	workspaceId1 = "workspaceId1"
	workspaceId2 = "workspaceId2"
	workspaceId3 = "workspaceId3"
)

var (
	partiallyOverwrittenWorkspace = Workspace{
		WorkspaceId: workspaceId2,
	}

	fullyOverwrittenWorkspace = Workspace{
		WorkspaceId:             workspaceId3,
		JiraMigrationFailedTag:  jiraMigrationFailedTag,
		JiraMigrationSkipTag:    jiraMigrationSkipTag,
		JiraMigrationSuccessTag: jiraMigrationSuccessTag,
	}

	config = Config{
		Global: Global{
			ClockifyToken: clockifyToken,
			Period:        period,
		},

		DefaultClient: Client{
			JiraClientUser: jiraClientUser,
			JiraPassword:   jiraPassword,
			JiraUsername:   jiraUsername,
			JiraHost:       jiraHost,
			StachurskyMode: stachurskyMode,
			Enabled:        enabled,
		},

		DefaultWorkspace: defaultWorkspace,

		Workspaces: Workspaces{
			workspaceId1: &defaultWorkspace,
			workspaceId2: &partiallyOverwrittenWorkspace,
			workspaceId3: &fullyOverwrittenWorkspace,
		},
	}
)

func TestReadWorkspacesConfig(t *testing.T) {

	t.Run("Full default workspaces configuration", func(t *testing.T) {
		got := config.combineWithDefaultConfig()

		assert.Strings(t, got[workspaceId1].WorkspaceId, workspaceIdDefault)
		assert.Strings(t, got[workspaceId1].JiraMigrationFailedTag, jiraMigrationFailedTagDefault)
		assert.Strings(t, got[workspaceId1].JiraMigrationSkipTag, jiraMigrationSkipTagDefault)
		assert.Strings(t, got[workspaceId1].JiraMigrationSuccessTag, jiraMigrationSuccessTagDefault)
	})

	t.Run("Override workspaceId configuration", func(t *testing.T) {
		got := config.combineWithDefaultConfig()

		assert.Strings(t, got[workspaceId2].WorkspaceId, workspaceId2)
		assert.Strings(t, got[workspaceId2].JiraMigrationFailedTag, jiraMigrationFailedTagDefault)
		assert.Strings(t, got[workspaceId2].JiraMigrationSkipTag, jiraMigrationSkipTagDefault)
		assert.Strings(t, got[workspaceId2].JiraMigrationSuccessTag, jiraMigrationSuccessTagDefault)
	})

	t.Run("Override full workspace configuration", func(t *testing.T) {
		got := config.combineWithDefaultConfig()

		assert.Strings(t, got[workspaceId3].WorkspaceId, workspaceId3)
		assert.Strings(t, got[workspaceId3].JiraMigrationFailedTag, jiraMigrationFailedTag)
		assert.Strings(t, got[workspaceId3].JiraMigrationSkipTag, jiraMigrationSkipTag)
		assert.Strings(t, got[workspaceId3].JiraMigrationSuccessTag, jiraMigrationSuccessTag)
	})
}

func TestGetWorkspace(t *testing.T) {
	t.Run("Get workspace configuration", func(t *testing.T) {
		got, err := config.GetWorkspace(workspaceId1)

		assert.Errors(t, err, nil)
		assert.Strings(t, got.WorkspaceId, workspaceIdDefault)
		assert.Strings(t, got.JiraMigrationFailedTag, jiraMigrationFailedTagDefault)
		assert.Strings(t, got.JiraMigrationSkipTag, jiraMigrationSkipTagDefault)
		assert.Strings(t, got.JiraMigrationSuccessTag, jiraMigrationSuccessTagDefault)
	})

	t.Run("Get error when workspace not found", func(t *testing.T) {
		_, err := config.GetWorkspace("missing-workspace")

		assert.Errors(t, err, ErrWorkspaceNotFound)
	})
}

func TestGetWorkspaces(t *testing.T) {
	t.Run("Get workspaces", func(t *testing.T) {
		workspaceSelector := []string{}
		workspaces, err := config.FindWorkspaces(workspaceSelector)

		assert.Errors(t, err, nil)
		assert.Ints(t, len(workspaces), 3)

		assert.Strings(t, workspaces[workspaceId1].WorkspaceId, workspaceIdDefault)
		assert.Strings(t, workspaces[workspaceId2].WorkspaceId, workspaceId2)
		assert.Strings(t, workspaces[workspaceId3].WorkspaceId, workspaceId3)
	})

	t.Run("Get workspaces matching selector", func(t *testing.T) {

		workspaceSelector := []string{workspaceId1, workspaceId3}
		workspaces, err := config.FindWorkspaces(workspaceSelector)

		assert.Errors(t, err, nil)
		assert.Ints(t, len(workspaces), len(workspaceSelector))
		assert.Strings(t, workspaces[workspaceId1].WorkspaceId, workspaceIdDefault)
		assert.Strings(t, workspaces[workspaceId3].WorkspaceId, workspaceId3)
	})

	t.Run("Get error when no matching workspaces found", func(t *testing.T) {

		workspaceSelector := []string{"no-matching-id"}
		_, err := config.FindWorkspaces(workspaceSelector)

		assert.Errors(t, err, ErrWorkspacesNotMatchingSelector)
	})

	t.Run("Get error when workspaces not configured", func(t *testing.T) {
		config := Config{
			Workspaces: Workspaces{},
		}
		workspaceSelector := []string{}
		_, err := config.FindWorkspaces(workspaceSelector)

		assert.Errors(t, err, ErrWorkspacesNotConfigured)
	})
}

func TestConfigOverwrites(t *testing.T) {
	t.Run("Overwrite period config", func(t *testing.T) {

		config.OverwritePeriodSetting(78)

		assert.Ints(t, config.Global.Period, 78)
	})

	t.Run("Overwrite precision config", func(t *testing.T) {

		config := Config{
			DefaultClient: Client{
				StachurskyMode: 10,
			},
			Workspaces: Workspaces{
				workspaceId1: &Workspace{
					Clients: Clients{
						clientId1: &Client{
							StachurskyMode: 20,
						},
						clientId2: &Client{
							StachurskyMode: 30,
						},
					},
				},
				workspaceId2: &Workspace{
					Clients: Clients{
						clientId3: &Client{
							StachurskyMode: 40,
						},
					},
				},
			},
		}

		config.OverwritePrecisionSetting(50)

		assert.Ints(t, config.DefaultClient.StachurskyMode, 50)
		assert.Ints(t, config.Workspaces[workspaceId1].Clients[clientId1].StachurskyMode, 50)
		assert.Ints(t, config.Workspaces[workspaceId1].Clients[clientId2].StachurskyMode, 50)
		assert.Ints(t, config.Workspaces[workspaceId2].Clients[clientId3].StachurskyMode, 50)
	})
}

func TestGetTimeInterval(t *testing.T) {

	now := time.Now()

	start, end := config.GetTimeInterval(&now)

	expectedStart := now.AddDate(0, 0, -config.Global.Period)

	assert.Strings(t, start.String(), expectedStart.String())
	assert.Strings(t, end.String(), now.String())
}

func TestError(t *testing.T) {
	t.Run("ErrNotFound", func(t *testing.T) {
		got := ConfigErr("Error message").Error()
		want := "Error message"

		assert.Strings(t, got, want)
	})
}

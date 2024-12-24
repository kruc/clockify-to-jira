package config

import (
	"testing"

	"github.com/kruc/clockify-to-jira/internal/assert"
)

const (
	workspaceId             = "workspaceId"
	jiraMigrationFailedTag  = "jiraMigrationFailedTag"
	jiraMigrationSkipTag    = "jiraMigrationSkipTag"
	jiraMigrationSuccessTag = "jiraMigrationSuccessTag"

	workspaceIdDefault             = "workspaceIdDefault"
	jiraMigrationFailedTagDefault  = "jiraMigrationFailedTagDefault"
	jiraMigrationSkipTagDefault    = "jiraMigrationSkipTagDefault"
	jiraMigrationSuccessTagDefault = "jiraMigrationSuccessTagDefault"
)

var (
	defaultWorkspace = Workspace{
		WorkspaceId:             workspaceIdDefault,
		JiraMigrationFailedTag:  jiraMigrationFailedTagDefault,
		JiraMigrationSkipTag:    jiraMigrationSkipTagDefault,
		JiraMigrationSuccessTag: jiraMigrationSuccessTagDefault,
		Clients:                 nil,
	}
)

func TestGetWorkspaceConfig(t *testing.T) {

	t.Run("Get workspace config", func(t *testing.T) {

		workspace := Workspace{
			WorkspaceId: workspaceId,
			Clients: Clients{
				clientId1: &Client{},
				clientId2: client,
			},
		}

		finalWorkspace := workspace.combineWithDefaultConfig(defaultWorkspace, defaultClient)

		assert.Strings(t, finalWorkspace.WorkspaceId, workspaceId)
		assert.Strings(t, finalWorkspace.JiraMigrationFailedTag, jiraMigrationFailedTagDefault)
		assert.Strings(t, finalWorkspace.JiraMigrationSkipTag, jiraMigrationSkipTagDefault)
		assert.Strings(t, finalWorkspace.JiraMigrationSuccessTag, jiraMigrationSuccessTagDefault)
		assert.Ints(t, len(finalWorkspace.Clients), 2)
		assert.Strings(t, finalWorkspace.Clients[clientId1].JiraClientUser, jiraClientUserDefault)
		assert.Strings(t, finalWorkspace.Clients[clientId1].JiraUsername, jiraUsernameDefault)
		assert.Strings(t, finalWorkspace.Clients[clientId1].JiraPassword, jiraPasswordDefault)
		assert.Strings(t, finalWorkspace.Clients[clientId1].JiraHost, jiraHostDefault)

		assert.Strings(t, finalWorkspace.Clients[clientId2].JiraClientUser, jiraClientUser)
		assert.Strings(t, finalWorkspace.Clients[clientId2].JiraUsername, jiraUsername)
		assert.Strings(t, finalWorkspace.Clients[clientId2].JiraPassword, jiraPassword)
		assert.Strings(t, finalWorkspace.Clients[clientId2].JiraHost, jiraHost)
	})

	t.Run("Partially combine with default workspace config", func(t *testing.T) {

		workspace := Workspace{
			WorkspaceId:             workspaceId,
			JiraMigrationSuccessTag: jiraMigrationSuccessTag,
			Clients: Clients{
				clientId1: &Client{
					JiraClientUser: jiraClientUser,
					JiraPassword:   jiraPassword,
				},
			},
		}

		finalWorkspace := workspace.combineWithDefaultConfig(defaultWorkspace, defaultClient)

		assert.Strings(t, finalWorkspace.WorkspaceId, workspaceId)
		assert.Strings(t, finalWorkspace.JiraMigrationFailedTag, jiraMigrationFailedTagDefault)
		assert.Strings(t, finalWorkspace.JiraMigrationSkipTag, jiraMigrationSkipTagDefault)
		assert.Strings(t, finalWorkspace.JiraMigrationSuccessTag, jiraMigrationSuccessTag)
		assert.Ints(t, len(finalWorkspace.Clients), 1)

		assert.Strings(t, finalWorkspace.Clients[clientId1].JiraClientUser, jiraClientUser)
		assert.Strings(t, finalWorkspace.Clients[clientId1].JiraPassword, jiraPassword)
		assert.Strings(t, finalWorkspace.Clients[clientId1].JiraHost, jiraHostDefault)
		assert.Strings(t, finalWorkspace.Clients[clientId1].JiraUsername, jiraUsernameDefault)
	})

	t.Run("Override default workspace config", func(t *testing.T) {

		workspace := Workspace{
			WorkspaceId:             workspaceId,
			JiraMigrationFailedTag:  jiraMigrationFailedTag,
			JiraMigrationSkipTag:    jiraMigrationSkipTag,
			JiraMigrationSuccessTag: jiraMigrationSuccessTag,

			Clients: Clients{
				clientId1: client,
				clientId2: &Client{},
			},
		}

		finalWorkspace := workspace.combineWithDefaultConfig(defaultWorkspace, defaultClient)

		assert.Strings(t, finalWorkspace.WorkspaceId, workspaceId)
		assert.Strings(t, finalWorkspace.JiraMigrationFailedTag, jiraMigrationFailedTag)
		assert.Strings(t, finalWorkspace.JiraMigrationSkipTag, jiraMigrationSkipTag)
		assert.Strings(t, finalWorkspace.JiraMigrationSuccessTag, jiraMigrationSuccessTag)
		assert.Ints(t, len(finalWorkspace.Clients), 2)

		assert.Strings(t, finalWorkspace.Clients[clientId1].JiraClientUser, jiraClientUser)
		assert.Strings(t, finalWorkspace.Clients[clientId2].JiraClientUser, jiraClientUserDefault)

		assert.Strings(t, finalWorkspace.Clients[clientId1].JiraClientUser, jiraClientUser)
		assert.Strings(t, finalWorkspace.Clients[clientId1].JiraPassword, jiraPassword)
		assert.Strings(t, finalWorkspace.Clients[clientId1].JiraHost, jiraHost)
		assert.Strings(t, finalWorkspace.Clients[clientId1].JiraUsername, jiraUsername)
	})
}

func TestGetClient(t *testing.T) {

	t.Run("Get client config", func(t *testing.T) {

		workspace := Workspace{
			Clients: Clients{
				clientId1: client,
			},
		}

		client, err := workspace.GetClient(clientId1)

		assert.Errors(t, err, nil)
		assert.Strings(t, client.JiraClientUser, client.JiraClientUser)
	})

	t.Run("Throw error if client not configured", func(t *testing.T) {

		workspace := Workspace{
			WorkspaceId:             workspaceId,
			JiraMigrationFailedTag:  jiraMigrationFailedTag,
			JiraMigrationSkipTag:    jiraMigrationSkipTag,
			JiraMigrationSuccessTag: jiraMigrationSuccessTag,

			Clients: nil,
		}

		_, err := workspace.GetClient("invalid-client-id")

		assert.Errors(t, err, ErrClientNotFound)
	})
}

func TestOverwriteWorkspaceClientsPrecisionConfig(t *testing.T) {
	workspace := Workspace{
		Clients: Clients{
			clientId1: client,
		},
	}

	workspace.overwritePrecisionSetting(70)

	assert.Ints(t, client.StachurskyMode, 70)
}

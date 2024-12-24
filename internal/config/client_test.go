package config

import (
	"testing"

	"github.com/kruc/clockify-to-jira/internal/assert"
)

const (
	clientId1 = "clientId1"
	clientId2 = "clientId2"
	clientId3 = "clientId3"

	jiraClientUserDefault = "jiraClientUserDefault"
	jiraPasswordDefault   = "jiraPasswordDefault"
	jiraUsernameDefault   = "jiraUsernameDefault"
	jiraHostDefault       = "jiraHostDefault"
	stachurskyModeDefault = 30
	enabledDefault        = false

	jiraClientUser = "jiraClientUser"
	jiraPassword   = "jiraPassword"
	jiraUsername   = "jiraUsername"
	jiraHost       = "jiraHost"
	stachurskyMode = 15
	enabled        = true
)

var (
	defaultClient = Client{
		JiraClientUser: jiraClientUserDefault,
		JiraHost:       jiraHostDefault,
		JiraPassword:   jiraPasswordDefault,
		JiraUsername:   jiraUsernameDefault,
		StachurskyMode: stachurskyModeDefault,
		Enabled:        enabledDefault,
	}

	client = &Client{
		JiraClientUser: jiraClientUser,
		JiraHost:       jiraHost,
		JiraPassword:   jiraPassword,
		JiraUsername:   jiraUsername,
		StachurskyMode: stachurskyMode,
	}
)

func TestGetClientConfig(t *testing.T) {
	clientId1 := client
	clientId2 := Client{
		JiraClientUser: jiraClientUser,
		JiraHost:       jiraHost,
	}
	clientId3 := Client{
		JiraClientUser: jiraClientUser,
		JiraHost:       jiraHost,
		JiraPassword:   jiraPassword,
		JiraUsername:   jiraUsername,
		StachurskyMode: stachurskyMode,
		Enabled:        enabled,
	}

	t.Run("Get client config", func(t *testing.T) {
		client := clientId1.combineWithDefaultConfig(defaultClient)

		assert.Strings(t, client.JiraClientUser, jiraClientUser)
		assert.Strings(t, client.JiraPassword, jiraPassword)
		assert.Strings(t, client.JiraUsername, jiraUsername)
		assert.Strings(t, client.JiraHost, jiraHost)
		assert.Ints(t, client.StachurskyMode, stachurskyMode)
		assert.Bools(t, client.Enabled, false)
	})

	t.Run("Partially combine with default client config", func(t *testing.T) {

		finalClient := clientId2.combineWithDefaultConfig(defaultClient)

		assert.Strings(t, finalClient.JiraClientUser, jiraClientUser)
		assert.Strings(t, finalClient.JiraPassword, jiraPasswordDefault)
		assert.Strings(t, finalClient.JiraUsername, jiraUsernameDefault)
		assert.Strings(t, finalClient.JiraHost, jiraHost)
		assert.Ints(t, finalClient.StachurskyMode, stachurskyModeDefault)
		assert.Bools(t, finalClient.Enabled, enabledDefault)
	})

	t.Run("Override default client config", func(t *testing.T) {
		finalClient := clientId3.combineWithDefaultConfig(defaultClient)

		assert.Strings(t, finalClient.JiraClientUser, jiraClientUser)
		assert.Strings(t, finalClient.JiraPassword, jiraPassword)
		assert.Strings(t, finalClient.JiraUsername, jiraUsername)
		assert.Strings(t, finalClient.JiraHost, jiraHost)
		assert.Ints(t, finalClient.StachurskyMode, stachurskyMode)
		assert.Bools(t, finalClient.Enabled, enabled)
	})
}

func TestOverwriteClientPrecisionConfig(t *testing.T) {
	client := Client{
		StachurskyMode: 10,
	}

	client.overwritePrecisionSetting(70)

	assert.Ints(t, client.StachurskyMode, 70)
}

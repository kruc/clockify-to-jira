package main

import (
	"strings"
)

type globalConfigKey struct {
	path           string
	name           string
	defaultValue   interface{}
	requiresChange bool
}

func (g globalConfigKey) fullName() string {
	if g.path != "" {
		return strings.Join([]string{g.path, g.name}, ".")
	}

	return g.name
}

func jiraHostConfig(path string) globalConfigKey {
	return globalConfigKey{path, "jira_host", "https://headstart.atlassian.net", false}
}
func jiraPasswordConfig(path string) globalConfigKey {
	return globalConfigKey{path, "jira_password", "(visit https://id.atlassian.com/manage/api-tokens)", true}
}
func jiraUsernameConfig(path string) globalConfigKey {
	return globalConfigKey{path, "jira_username", "firstname.lastname@domain.com", true}
}
func jiraClientUserConfig(path string) globalConfigKey {
	return globalConfigKey{path, "jira_client_user", "firstname.lastname", true}
}
func stachurskyModeConfig(path string) globalConfigKey {
	return globalConfigKey{path, "stachursky_mode", 1, false}
}
func logFormatConfig() globalConfigKey {
	return globalConfigKey{"", "log_format", "text", false}
}
func logOutputConfig() globalConfigKey {
	return globalConfigKey{"", "log_output", "stdout", false}
}
func jiraMigrationSuccessTagConfig() globalConfigKey {
	return globalConfigKey{"", "jira_migration_success_tag", "logged", false}
}
func jiraMigrationFailedTagConfig() globalConfigKey {
	return globalConfigKey{"", "jira_migration_failed_tag", "jira-migration-failed", false}
}
func jiraMigrationSkipTagConfig() globalConfigKey {
	return globalConfigKey{"", "jira_migration_skip_tag", "jira-migration-skip", false}
}
func periodConfig() globalConfigKey {
	return globalConfigKey{"", "period", 1, false}
}
func clockifyTokenConfig() globalConfigKey {
	return globalConfigKey{"", "clockify_token", "(visit https://clockify.me/user/settings)", true}
}
func workspaceIDConfig() globalConfigKey {
	return globalConfigKey{"", "workspace_id", "(visit https://clockify.me/workspaces click settings and get id from url)", true}
}

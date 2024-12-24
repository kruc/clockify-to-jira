package config

import (
	"bytes"
	"os"
	"testing"

	"github.com/kruc/clockify-to-jira/internal/assert"
)

const (
	tmpComfigBasePath = "./tmp/"
)

func TestGenerateDefaultConfig(t *testing.T) {

	t.Run("Generate default config template", func(t *testing.T) {
		buffer := bytes.Buffer{}

		generateClientConfigTemplate(&buffer)

		want := `global:
  clockify_token: (visit https://app.clockify.me/user/preferences#advanced)
  period: 7
default_client:
  jira_client_user: firstname.lastname
  jira_host: https://headstart.atlassian.net
  jira_username: firstname.lastname@domain.com
  jira_password: (visit https://id.atlassian.com/manage/api-tokens)
  stachursky_mode: 15
  enabled: false
default_workspace:
  workspace_id: (visit https://app.clockify.me/workspaces -> settings -> id from url)
  jira_migration_failed_tag: jira-migration-failed
  jira_migration_skip_tag: jira-migration-skip
  jira_migration_success_tag: logged
  clients: {}
workspaces: {}
`

		assert.Strings(t, buffer.String(), want)
	})
}

func TestSaveConfigToYamlFile(t *testing.T) {
	t.Run("Write default configuration to yaml file", func(t *testing.T) {
		configFilePath := tmpComfigBasePath + "config.yaml"

		defer os.RemoveAll(tmpComfigBasePath)

		err := InitializeConfig(configFilePath)

		assert.Errors(t, err, nil)
		assert.FileExists(t, configFilePath)
	})

	t.Run("Write default configuration to yaml file in subdirectory", func(t *testing.T) {
		configFilePath := tmpComfigBasePath + "subdir1/subdir2/config.yaml"

		defer os.RemoveAll(tmpComfigBasePath)

		err := InitializeConfig(configFilePath)

		assert.Errors(t, err, nil)
		assert.FileExists(t, configFilePath)
	})

	t.Run("Throws error on empty file path", func(t *testing.T) {
		configFilePath := ""
		err := InitializeConfig(configFilePath)

		assert.Errors(t, err, ErrGeneratorConfigFileError)
	})

	t.Run("Throws error if file already exists", func(t *testing.T) {
		configFilePath := tmpComfigBasePath + "existsing-config.yaml"
		defer os.RemoveAll(tmpComfigBasePath)

		err := InitializeConfig(configFilePath)

		assert.Errors(t, err, nil)

		err = InitializeConfig(configFilePath)

		assert.Errors(t, err, ErrGeneratorConfigFileAlreadyExists)
	})
}

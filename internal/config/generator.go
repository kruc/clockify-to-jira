package config

import (
	"errors"
	"io"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

const (
	ErrGeneratorConfigFileError         = ConfigErr("Error during creating configuration file")
	ErrGeneratorConfigFileAlreadyExists = ConfigErr("Configuration file already exists")
)

var (
	configTemplate = Config{
		Global: Global{
			ClockifyToken: "(visit https://app.clockify.me/user/preferences#advanced)",
			Period:        7,
		},
		DefaultClient: Client{
			JiraClientUser: "firstname.lastname",
			JiraHost:       "https://headstart.atlassian.net",
			JiraUsername:   "firstname.lastname@domain.com",
			JiraPassword:   "(visit https://id.atlassian.com/manage/api-tokens)",
			StachurskyMode: 15,
			Enabled:        false,
		},
		DefaultWorkspace: Workspace{
			WorkspaceId:             "(visit https://app.clockify.me/workspaces -> settings -> id from url)",
			JiraMigrationFailedTag:  "jira-migration-failed",
			JiraMigrationSkipTag:    "jira-migration-skip",
			JiraMigrationSuccessTag: "logged",
		},
		Workspaces: Workspaces{},
	}
)

func InitializeConfig(filePath string) error {

	fileExists := checkFileExists(filePath)

	if fileExists {
		return ErrGeneratorConfigFileAlreadyExists
	}

	createConfigDirectoryStructure(filePath)

	configFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0755)

	if err != nil {
		return ErrGeneratorConfigFileError
	}

	defer configFile.Close()

	generateClientConfigTemplate(configFile)

	return nil
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)

	return !errors.Is(error, os.ErrNotExist)
}

func generateClientConfigTemplate(writer io.Writer) {
	yamlData, _ := yaml.Marshal(configTemplate)

	writer.Write(yamlData)
}

func createConfigDirectoryStructure(filePath string) {

	dirPath := path.Dir(filePath)
	os.MkdirAll(dirPath, 0755)
}

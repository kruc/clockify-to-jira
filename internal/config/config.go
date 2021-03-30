package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ClientConfig avaiable fields
type ClientConfig struct {
	JiraUsername   string
	JiraPassword   string
	JiraClientUser string
	JiraHost       string
	StachurskyMode int
}

// GlobalConfigType struct
type GlobalConfigType struct {
	DefaultClient           ClientConfig
	Period                  int
	LogFormat               string
	LogOutput               string
	JiraMigrationSuccessTag string
	JiraMigrationFailedTag  string
	JiraMigrationSkipTag    string
	WorkspaceID             string
}

var (
	config       = "config"
	configPath   string
	globalConfig GlobalConfigType
)

func init() {

	if !checkConfiguration() {
		os.Exit(1)
	}

	globalConfig = parseGlobalConfig()

	flag.IntVarP(&globalConfig.Period, "period", "p", globalConfig.Period, "Migrate time entries from last given days")
	flag.StringVarP(&globalConfig.LogFormat, "format", "f", globalConfig.LogFormat, "Log format (text|json)")
	flag.StringVarP(&globalConfig.LogOutput, "output", "o", globalConfig.LogOutput, "Log output (stdout|filename)")
	flag.StringVarP(&globalConfig.WorkspaceID, "workspace", "w", globalConfig.WorkspaceID, "Clockify workspace id")
	flag.IntVarP(&globalConfig.DefaultClient.StachurskyMode, "tryb-niepokorny", "t", globalConfig.DefaultClient.StachurskyMode, "Rounding up the value of logged time up (minutes)")
}

// GetGlobalConfig provides global configuration
func GetGlobalConfig() GlobalConfigType {
	return globalConfig
}

// checkConfiguration validate configuration
func checkConfiguration() bool {

	configPath = fmt.Sprintf("%v/.clockify-to-jira", os.Getenv("HOME"))
	os.MkdirAll(configPath, 0755)
	os.OpenFile(fmt.Sprintf("%v/%v.yaml", configPath, config), os.O_CREATE|os.O_RDWR, 0666)

	viper.SetConfigName(config)
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	globalConfigKeys := []globalConfigKey{
		jiraHostConfig("default_client"),
		jiraPasswordConfig("default_client"),
		jiraUsernameConfig("default_client"),
		jiraClientUserConfig("default_client"),
		stachurskyModeConfig("default_client"),
		logFormatConfig(),
		logOutputConfig(),
		jiraMigrationSuccessTagConfig(),
		jiraMigrationFailedTagConfig(),
		jiraMigrationSkipTagConfig(),
		periodConfig(),
		clockifyTokenConfig(),
		workspaceIDConfig(),
	}

	configReady := true
	log.Info("Checking configuration...")
	for _, configKey := range globalConfigKeys {
		if !viper.IsSet(configKey.fullName()) {
			log.Warnf("Missing key: %v - config updated!\n", configKey.name)
			viper.Set(fmt.Sprintf("%v", configKey.fullName()), configKey.defaultValue)
			configReady = false
		} else if configKey.requiresChange && viper.GetString(configKey.fullName()) == configKey.defaultValue {

			log.Errorf("Invalid value for %v - change it", configKey.fullName())
			configReady = false
		}
	}

	err = viper.WriteConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"configPath": configPath,
		}).Error(err)
		return false
	}

	log.Infof("Customize configuration in file: %v/config.yaml\n", configPath)

	return configReady
}

// GenerateClientConfigTemplate generate clients template
func GenerateClientConfigTemplate(configPath string) {
	fmt.Printf("Generating config template for %v...\n", configPath)
	viper.Set(jiraUsernameConfig(configPath).fullName(), "FILL_IT OR REMOVE TO USE DEFAULT_CLIENT")
	viper.Set(jiraPasswordConfig(configPath).fullName(), "FILL_IT OR REMOVE TO USE DEFAULT_CLIENT")
	viper.Set(jiraClientUserConfig(configPath).fullName(), "FILL_IT OR REMOVE TO USE DEFAULT_CLIENT")
	viper.Set(jiraHostConfig(configPath).fullName(), "FILL_IT OR REMOVE TO USE DEFAULT_CLIENT")
	viper.Set(stachurskyModeConfig(configPath).fullName(), "FILL_IT OR REMOVE TO USE DEFAULT_CLIENT")
	viper.Set(fmt.Sprintf("%v.%v", configPath, "enabled"), false)
	err := viper.WriteConfig()

	if err != nil {
		log.WithFields(log.Fields{
			"configPath": configPath,
		}).Error(err)
		return
	}
	log.WithFields(log.Fields{
		"configPath": configPath,
	}).Info("Client config template created!\n")
}

// parseGlobalConfig parsing config from confguration file
func parseGlobalConfig() GlobalConfigType {
	clientDefaultConfigPath := "default_client"

	globalConfig := GlobalConfigType{
		DefaultClient:           ParseClientConfig(clientDefaultConfigPath, globalConfig),
		Period:                  viper.GetInt(periodConfig().name),
		LogFormat:               viper.GetString(logFormatConfig().name),
		LogOutput:               viper.GetString(logOutputConfig().name),
		JiraMigrationSuccessTag: viper.GetString(jiraMigrationSuccessTagConfig().name),
		JiraMigrationFailedTag:  viper.GetString(jiraMigrationFailedTagConfig().name),
		JiraMigrationSkipTag:    viper.GetString(jiraMigrationSkipTagConfig().name),
		WorkspaceID:             viper.GetString(workspaceIDConfig().name),
	}

	return globalConfig
}

// ParseClientConfig parse client config
func ParseClientConfig(clientConfigPath string, globalConfig GlobalConfigType) ClientConfig {

	clientConfig := ClientConfig{
		JiraUsername:   getString(jiraUsernameConfig(clientConfigPath).fullName(), globalConfig.DefaultClient.JiraUsername),
		JiraPassword:   getString(jiraPasswordConfig(clientConfigPath).fullName(), globalConfig.DefaultClient.JiraPassword),
		JiraClientUser: getString(jiraClientUserConfig(clientConfigPath).fullName(), globalConfig.DefaultClient.JiraClientUser),
		JiraHost:       getString(jiraHostConfig(clientConfigPath).fullName(), globalConfig.DefaultClient.JiraHost),
		StachurskyMode: getInt(stachurskyModeConfig(clientConfigPath).fullName(), globalConfig.DefaultClient.StachurskyMode),
	}

	if flag.CommandLine.Changed("tryb-niepokorny") {
		clientConfig.StachurskyMode = globalConfig.DefaultClient.StachurskyMode
	}

	return clientConfig
}

func getString(key, defaultValue string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}

	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if viper.IsSet(key) {
		return viper.GetInt(key)
	}

	return defaultValue
}

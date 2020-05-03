package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type clientConfig struct {
	jiraUsername   string
	jiraPassword   string
	jiraClientUser string
	jiraHost       string
	stachurskyMode int
}

type globalConfigType struct {
	defaultClient           clientConfig
	period                  int
	logFormat               string
	logOutput               string
	jiraMigrationSuccessTag string
	jiraMigrationFailedTag  string
	jiraMigrationSkipTag    string
	workspaceID             string
}

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

func generateClientConfigTemplate(configPath string) {
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

func parseGlobalConfig() globalConfigType {
	clientDefaultConfigPath := "default_client"

	globalConfig := globalConfigType{
		defaultClient:           parseClientConfig(clientDefaultConfigPath, globalConfig),
		period:                  viper.GetInt(periodConfig().name),
		logFormat:               viper.GetString(logFormatConfig().name),
		logOutput:               viper.GetString(logOutputConfig().name),
		jiraMigrationSuccessTag: viper.GetString(jiraMigrationSuccessTagConfig().name),
		jiraMigrationFailedTag:  viper.GetString(jiraMigrationFailedTagConfig().name),
		jiraMigrationSkipTag:    viper.GetString(jiraMigrationSkipTagConfig().name),
		workspaceID:             viper.GetString(workspaceIDConfig().name),
	}

	return globalConfig
}

func parseClientConfig(clientConfigPath string, globalConfig globalConfigType) clientConfig {

	clientConfig := clientConfig{
		jiraUsername:   getString(jiraUsernameConfig(clientConfigPath).fullName(), globalConfig.defaultClient.jiraUsername),
		jiraPassword:   getString(jiraPasswordConfig(clientConfigPath).fullName(), globalConfig.defaultClient.jiraPassword),
		jiraClientUser: getString(jiraClientUserConfig(clientConfigPath).fullName(), globalConfig.defaultClient.jiraClientUser),
		jiraHost:       getString(jiraHostConfig(clientConfigPath).fullName(), globalConfig.defaultClient.jiraHost),
		stachurskyMode: getInt(stachurskyModeConfig(clientConfigPath).fullName(), globalConfig.defaultClient.stachurskyMode),
	}

	if flag.CommandLine.Changed("tryb-niepokorny") {
		clientConfig.stachurskyMode = globalConfig.defaultClient.stachurskyMode
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

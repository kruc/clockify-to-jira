package config

import (
	"slices"
	"time"
)

const (
	ErrWorkspaceNotFound             = ConfigErr("Cannot find workspace configuration")
	ErrWorkspacesNotConfigured       = ConfigErr("Cannot find workspaces in configuration")
	ErrWorkspacesNotMatchingSelector = ConfigErr("Cannot find workspaces matching provided selector")
)

type ConfigErr string

func (e ConfigErr) Error() string {
	return string(e)
}

type Workspaces map[string]*Workspace

type Global struct {
	ClockifyToken string `yaml:"clockify_token"`
	Period        int    `yaml:"period"`
}

type Config struct {
	Global           Global     `yaml:"global"`
	DefaultClient    Client     `yaml:"default_client"`
	DefaultWorkspace Workspace  `yaml:"default_workspace"`
	Workspaces       Workspaces `yaml:"workspaces"`
}

func (c *Config) GetWorkspace(workspaceId string) (*Workspace, error) {

	workspace, ok := c.Workspaces[workspaceId]

	if !ok {
		return &Workspace{}, ErrWorkspaceNotFound
	}

	return workspace, nil
}

func (c *Config) FindWorkspaces(workspaceSelector []string) (Workspaces, error) {

	workspaces := c.Workspaces

	if len(workspaces) == 0 {
		return Workspaces{}, ErrWorkspacesNotConfigured
	}

	if len(workspaceSelector) != 0 {
		err := workspaces.filterBasedOnSelector(workspaceSelector)

		if err != nil {
			return nil, err
		}
	}

	return workspaces, nil
}

func (c *Config) OverwritePeriodSetting(period int) {
	c.Global.Period = period
}

func (c *Config) OverwritePrecisionSetting(precision int) {
	c.DefaultClient.StachurskyMode = precision

	for key := range c.Workspaces {
		c.Workspaces[key].overwritePrecisionSetting(precision)
	}
}

func (c *Config) GetTimeInterval(now *time.Time) (time.Time, time.Time) {
	return now.AddDate(0, 0, -c.Global.Period), *now
}

func (c *Config) combineWithDefaultConfig() Workspaces {

	workspaceList := Workspaces{}

	for key, workspace := range c.Workspaces {
		workspaceList[key] = workspace.combineWithDefaultConfig(c.DefaultWorkspace, c.DefaultClient)
	}

	return workspaceList
}

func (w *Workspaces) filterBasedOnSelector(selector []string) error {
	for key := range *w {
		if slices.Contains(selector, key) {
			continue
		}

		delete(*w, key)
	}

	if len(*w) == 0 {
		return ErrWorkspacesNotMatchingSelector
	}

	return nil
}

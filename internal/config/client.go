package config

type Client struct {
	JiraClientUser string `yaml:"jira_client_user"`
	JiraHost       string `yaml:"jira_host"`
	JiraUsername   string `yaml:"jira_username"`
	JiraPassword   string `yaml:"jira_password"`
	StachurskyMode int    `yaml:"stachursky_mode"`
	Enabled        bool   `yaml:"enabled"`
}

func (c *Client) combineWithDefaultConfig(defaultClient Client) *Client {

	client := defaultClient

	if c.JiraClientUser != "" {
		client.JiraClientUser = c.JiraClientUser
	}

	if c.JiraPassword != "" {
		client.JiraPassword = c.JiraPassword
	}

	if c.JiraUsername != "" {
		client.JiraUsername = c.JiraUsername
	}

	if c.JiraHost != "" {
		client.JiraHost = c.JiraHost
	}

	if c.StachurskyMode != 0 {
		client.StachurskyMode = c.StachurskyMode
	}

	if c.Enabled {
		client.Enabled = c.Enabled
	}

	return &client
}

func (c *Client) overwritePrecisionSetting(precision int) {
	c.StachurskyMode = precision
}

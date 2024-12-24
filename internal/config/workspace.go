package config

const (
	ErrClientNotFound = ConfigErr("Cannot find client in given workspace")
)

type Clients map[string]*Client

type Workspace struct {
	WorkspaceId             string  `yaml:"workspace_id"`
	JiraMigrationFailedTag  string  `yaml:"jira_migration_failed_tag"`
	JiraMigrationSkipTag    string  `yaml:"jira_migration_skip_tag"`
	JiraMigrationSuccessTag string  `yaml:"jira_migration_success_tag"`
	Clients                 Clients `yaml:"clients"`
}

func (w *Workspace) GetClient(clientId string) (*Client, error) {

	client, ok := w.Clients[clientId]

	if !ok {
		return &Client{}, ErrClientNotFound
	}

	return client, nil
}

func (w *Workspace) combineWithDefaultConfig(defaultWorkspace Workspace, defaultClient Client) *Workspace {
	workspace := defaultWorkspace

	if w.WorkspaceId != "" {
		workspace.WorkspaceId = w.WorkspaceId
	}

	if w.JiraMigrationFailedTag != "" {
		workspace.JiraMigrationFailedTag = w.JiraMigrationFailedTag
	}

	if w.JiraMigrationSkipTag != "" {
		workspace.JiraMigrationSkipTag = w.JiraMigrationSkipTag
	}

	if w.JiraMigrationSuccessTag != "" {
		workspace.JiraMigrationSuccessTag = w.JiraMigrationSuccessTag
	}

	if workspace.Clients == nil {
		workspace.Clients = Clients{}
	}

	for id, client := range w.Clients {
		workspace.Clients[id] = client.combineWithDefaultConfig(defaultClient)
	}

	return &workspace
}

func (w *Workspace) overwritePrecisionSetting(precision int) {

	for id := range w.Clients {
		w.Clients[id].overwritePrecisionSetting(precision)
	}
}

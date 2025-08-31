package github

type InstallationResponse struct {
	ID             uint   `json:"id"`
	InstallationID string `json:"installation_id"`
	AccountLogin   string `json:"account_login"`
	AccountType    string `json:"account_type"`
}

type WebhookEvent struct {
	Action       string                 `json:"action"`
	Installation map[string]interface{} `json:"installation"`
	Sender       map[string]interface{} `json:"sender"`
}

type RepositoryDispatchPayload struct {
	EventType     string      `json:"event_type"`
	ClientPayload interface{} `json:"client_payload"`
}

type ConfigAPIPayload struct {
	Path   string      `json:"path"`
	Action string      `json:"action"`
	Spec   interface{} `json:"spec"`
}

type GitHubRepo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    Owner  `json:"owner"`
	Private  bool   `json:"private"`
}

type Owner struct {
	Login string `json:"login"`
}
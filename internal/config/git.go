package config

type GitConfig struct {
	Repos []Repository `json:"repositories"`
}

type MessageVariables struct {
	OrdinalNumber string
	Week          int
}

type Repository struct {
	Path            string `json:"path"`
	MessageTemplate string `json:"message_template"`
}

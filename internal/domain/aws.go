package domain

type Credential struct {
	Profile      string `json:"name"`
	Key          string `json:"key"`
	Secret       string `json:"secret"`
	SessionToken string `json:"session_token"`
}

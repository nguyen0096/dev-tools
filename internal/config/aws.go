package config

type AWSConfig struct {
	MFAs []MFA `yaml:"mfas"`
}

type Credential struct {
	Profile      string `json:"name"`
	Key          string `json:"key"`
	Secret       string `json:"secret"`
	SessionToken string `json:"session_token"`
}

type MFA struct {
	Profile         string `json:"profile"`
	Device          string `json:"device"`
	SessionDuration int    `json:"session_duration"`
	OutputProfile   string `json:"output_profile"`
}

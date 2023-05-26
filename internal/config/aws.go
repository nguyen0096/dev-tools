package config

type AWSConfig struct {
	MFAs []MFA `yaml:"mfas"`
}

type AWSProfile struct {
	Name   string `yaml:"name"`
	Params struct {
		Region        string `yaml:"region"`
		Output        string `yaml:"output"`
		RoleARN       string `yaml:"role_arn"`
		SourceProfile string `yaml:"source_profile"`
	} `yaml:"params"`
}

type Credential struct {
	Profile      string `yaml:"name"`
	Key          string `yaml:"key"`
	Secret       string `yaml:"secret"`
	SessionToken string `yaml:"session_token"`
}

type MFA struct {
	Profile         string `yaml:"profile"`
	Device          string `yaml:"device"`
	SessionDuration int64  `yaml:"session_duration"`
	OutputProfile   string `yaml:"output_profile"`
}

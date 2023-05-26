package aws

import (
	"dev-tools/pkg/shell"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dlclark/regexp2"
	"github.com/spf13/cobra"

	"dev-tools/internal/config"
)

const (
	// flags const
	flagNameProfile  = "profile"
	flagNameMFAToken = "mfa_token"
)

var (
	setupMfaFlags = struct {
		token string
	}{}
)

type awsGetSessionTokenResponse struct {
	Credential struct {
		AccessKeyID     string `json:"AccessKeyId"`
		SecretAccessKey string `json:"SecretAccessKey"`
		SessionToken    string `json:"SessionToken"`
	} `json:"Credentials"`
}

func (a *AWSCommander) cmdSetupMFA() *cobra.Command {
	setupMfaCmd := &cobra.Command{
		Use:   "mfa",
		Short: "setup mfa auth session",
		Run: func(cmd *cobra.Command, args []string) {
			// loop through MFA config and setup mfa session
			for _, mfaProfile := range config.Cfg.AWS.MFAs {
				a.setupMFA(mfaProfile)
			}
		},
	}

	// Setup flags
	setupMfaCmd.Flags().StringVarP(&setupMfaFlags.token, flagNameMFAToken, "t", "", "MFA code")
	setupMfaCmd.MarkFlagRequired(flagNameMFAToken)

	return setupMfaCmd
}

func (a *AWSCommander) setupMFA(mfa config.MFA) error {
	raw, err := shell.ExecOutput("aws", "sts", "get-session-token", "--serial-number", mfa.Device, "--profile", mfa.Profile, "--token-code", setupMfaFlags.token, "--duration-seconds", fmt.Sprint(mfa.SessionDuration))
	if err != nil {
		a.run.IO.Printf("failed to run aws sts get-session-token")
	}

	res := &awsGetSessionTokenResponse{}
	if err := json.Unmarshal(raw, &res); err != nil {
		return err
	}

	// reoutput config again along with mfa session credential
	if err := a.updateAWSCredentials([]config.Credential{
		{
			Profile:      mfa.OutputProfile,
			Key:          res.Credential.AccessKeyID,
			Secret:       res.Credential.SecretAccessKey,
			SessionToken: res.Credential.SessionToken,
		},
	}); err != nil {
		return err
	}

	if err := shell.Exec("aws", "sts", "get-caller-identity", "--profile", mfa.OutputProfile); err != nil {
		return err
	}

	raw, err = shell.ExecOutput("aws", "eks", "update-kubeconfig", "--name", "prophet-dev", "--profile", mfa.OutputProfile)
	if err != nil {
		return err
	}

	rg, err := regexp2.Compile("Updated context (arn:.*) in .*", regexp2.None)
	if err != nil {
		return err
	}

	matched, err := rg.FindStringMatch(string(raw))
	if err != nil {
		return err
	}

	if matched == nil || len(matched.Captures) == 0 {
		return errors.New("failed to get kubectl context from output")
	}

	err = shell.Exec("kubectl", "config", "use-context", matched.Captures[0].String())
	if err != nil {
		return err
	}

	return nil
}

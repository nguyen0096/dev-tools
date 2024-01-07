package aws

import (
	"dev-tools/config"
	"dev-tools/internal/domain"
	"dev-tools/pkg/shell"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/dlclark/regexp2"
	"github.com/spf13/cobra"
)

func GetAWSCommand() *cobra.Command {
	awsCmd := &cobra.Command{
		Use:   "aws",
		Short: "Setup AWS and Kubectl config",
	}

	awsCmd.AddCommand(cmdSetupMFA())

	return awsCmd
}

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

func cmdSetupMFA() *cobra.Command {
	setupMfaCmd := &cobra.Command{
		Use:   "mfa",
		Short: "setup mfa auth session",
		Run: func(cmd *cobra.Command, args []string) {
			// loop through MFA config and setup mfa session
			for _, mfaProfile := range config.Cfg.AWS.MFAs {
				setupMFA(mfaProfile)
			}
		},
	}

	// Setup flags
	setupMfaCmd.Flags().StringVarP(&setupMfaFlags.token, flagNameMFAToken, "t", "", "MFA code")
	setupMfaCmd.MarkFlagRequired(flagNameMFAToken)

	return setupMfaCmd
}

func setupMFA(mfa config.MFA) error {
	log.Printf("session duration: %v", mfa.SessionDuration)
	raw, err := shell.ExecOutput("aws", "sts", "get-session-token", "--serial-number", mfa.Device, "--profile", mfa.Profile, "--token-code", setupMfaFlags.token, "--duration-seconds", strconv.Itoa(mfa.SessionDuration))
	if err != nil {
		fmt.Printf("failed to run aws sts get-session-token. err: %v", err)
	}

	res := &awsGetSessionTokenResponse{}
	if err := json.Unmarshal(raw, &res); err != nil {
		return err
	}

	// reoutput config again along with mfa session credential
	if err := updateAWSCredentials([]domain.Credential{
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

	raw, err = shell.ExecOutput("aws", "eks", "update-kubeconfig", "--name", mfa.ClusterName, "--profile", mfa.OutputProfile)
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

func updateAWSCredentials(creds []domain.Credential) error {
	for _, cred := range creds {
		setAccessKeyCmd := []string{"aws", "configure", "set", "AWS_ACCESS_KEY_ID", cred.Key, "--profile", cred.Profile}
		if o, err := shell.ExecOutput(setAccessKeyCmd...); err != nil {
			err = fmt.Errorf("failed to update aws access key for profile [%s]\n%s", cred.Profile, o)
			return err
		}

		setAccessSecretCmd := []string{"aws", "configure", "set", "AWS_SECRET_ACCESS_KEY", cred.Secret, "--profile", cred.Profile}
		if o, err := shell.ExecOutput(setAccessSecretCmd...); err != nil {
			err = fmt.Errorf("failed to update aws access secret for profile [%s]\n%s", cred.Profile, o)
			return err
		}

		setSessionTokenCmd := []string{"aws", "configure", "set", "AWS_SESSION_TOKEN", cred.SessionToken, "--profile", cred.Profile}
		if o, err := shell.ExecOutput(setSessionTokenCmd...); err != nil {
			err = fmt.Errorf("failed to update aws session token for profile [%s]\n%s", cred.Profile, o)
			return err
		}
		fmt.Printf("updated aws credentials for profile [%s]", cred.Profile)
	}

	return nil
}

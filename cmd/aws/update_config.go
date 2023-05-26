package aws

import (
	"dev-tools/internal/config"
	"dev-tools/pkg/shell"
)

func (a *AWSCommander) updateAWSCredentials(creds []config.Credential) error {
	for _, cred := range creds {
		setAccessKeyCmd := []string{"aws", "configure", "set", "AWS_ACCESS_KEY_ID", cred.Key, "--profile", cred.Profile}
		if err := shell.Exec(setAccessKeyCmd...); err != nil {
			a.run.IO.Errorf("failed to update aws access key for profile [%s]", cred.Profile)
		}

		setAccessSecretCmd := []string{"aws", "configure", "set", "AWS_SECRET_ACCESS_KEY", cred.Secret, "--profile", cred.Profile}
		if err := shell.Exec(setAccessSecretCmd...); err != nil {
			a.run.IO.Errorf("failed to update aws access secret for profile [%s]", cred.Profile)
		}

		setSessionTokenCmd := []string{"aws", "configure", "set", "AWS_SESSION_TOKEN", cred.SessionToken, "--profile", cred.Profile}
		if err := shell.Exec(setSessionTokenCmd...); err != nil {
			a.run.IO.Errorf("failed to update aws session token for profile [%s]", cred.Profile)
		}

		a.run.IO.Printf("updated aws credentials for profile [%s]", cred.Profile)
	}

	return nil
}

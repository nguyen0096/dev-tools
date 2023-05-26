package aws

import (
	"github.com/spf13/cobra"

	"dev-tools/internal/cmd"
)

type AWSCommander struct {
	run *cmd.CmdRun
}

func NewAWSCommander(run *cmd.CmdRun) *AWSCommander {
	return &AWSCommander{
		run: run,
	}
}

func (a *AWSCommander) GetCommand() *cobra.Command {
	awsCmd := &cobra.Command{
		Use: "aws",
	}

	awsCmd.AddCommand(a.cmdSetupMFA())

	return awsCmd
}

package root

import (
	"github.com/spf13/cobra"

	"dev-tools/cmd/aws"
	"dev-tools/cmd/git"
	"dev-tools/internal/cmd"
)

func NewRootCmd(r *cmd.CmdRun) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "neyu <command> <subcommand> [flags]",
		Short: "Neyu dev-tools CLI",
	}

	cmd.AddCommand(aws.NewAWSCommander(r).GetCommand())
	cmd.AddCommand(git.NewGitCommander(r).GetCommand())

	return cmd
}

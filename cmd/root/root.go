package root

import (
	"github.com/spf13/cobra"

	"dev-tools/cmd/aws"
	"dev-tools/internal/cmd"
)

func NewRootCmd(r *cmd.CmdRun) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "neyu <command> <subcommand> [flags]",
		Short: "Neyu dev-tools CLI",
	}

	cmd.AddCommand(aws.NewAWSCommander(r).GetCommand())

	return cmd
}

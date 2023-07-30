package main

import (
	"dev-tools/config"
	"dev-tools/internal/aws"
	"dev-tools/internal/git"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "ndv",
		Short: "A binary tools for working process automation",
	}
	config.MustLoadConfig()
	cmd.AddCommand(aws.GetAWSCommand())
	cmd.AddCommand(git.GetGitCommand())
	cmd.Execute()
}

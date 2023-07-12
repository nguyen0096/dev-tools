package git

import (
	"dev-tools/internal/cmd"

	"github.com/spf13/cobra"
)

type GitCommander struct {
	run *cmd.CmdRun
}

func NewGitCommander(run *cmd.CmdRun) *GitCommander {
	return &GitCommander{
		run: run,
	}
}

func (g *GitCommander) GetCommand() *cobra.Command {
	gitCmd := &cobra.Command{
		Use: "aws",
	}

	return gitCmd
}

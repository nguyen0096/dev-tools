package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"dev-tools/internal/config"
	"dev-tools/pkg/iostreams"
)

func HasCommand(rootCmd *cobra.Command, args []string) bool {
	c, _, err := rootCmd.Traverse(args)
	return err == nil && c != rootCmd
}

// CmdRun is the context of root command run,
// contains all shared dependencies between commands
type CmdRun struct {
	IO iostreams.IOStreams
}

func New() *CmdRun {
	return &CmdRun{
		IO: iostreams.New(),
	}
}

func (c *CmdRun) Exec(root *cobra.Command) ExitCode {
	expandedArgs := []string{}
	if len(os.Args) > 0 {
		expandedArgs = os.Args[1:]
	}

	if err := config.LoadConfig(); err != nil {
		c.IO.Printf("failed to load config. err: %v", err)
		return ExitCodeError
	}

	if !HasCommand(root, expandedArgs) {
		c.IO.Errorf("Error: command not found")
		return ExitCodeError
	}

	root.SetArgs(expandedArgs)
	if _, err := root.ExecuteC(); err != nil {
		return ExitCodeError
	}

	return ExitCodeOK
}

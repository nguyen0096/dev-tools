package main

import (
	"os"

	"dev-tools/cmd/root"
	"dev-tools/internal/cmd"
)

func main() {
	run := cmd.New()
	rootCmd := root.NewRootCmd(run)
	exitCode := run.Exec(rootCmd)
	os.Exit(int(exitCode))
}

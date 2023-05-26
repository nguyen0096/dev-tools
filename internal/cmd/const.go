package cmd

type ExitCode int

const (
	ExitCodeOK     ExitCode = 0
	ExitCodeError  ExitCode = 1
	ExitCodeCancel ExitCode = 2
	ExitCodeAuth   ExitCode = 4
)

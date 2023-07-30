package shell

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/cli/safeexec"
)

func Exec(args ...string) error {
	exe, err := safeexec.LookPath(args[0])
	if err != nil {
		return err
	}
	announce(args...)
	cmd := exec.Command(exe, args[1:]...)
	return cmd.Run()
}

func ExecOutput(args ...string) ([]byte, error) {
	exe, err := safeexec.LookPath(args[0])
	if err != nil {
		return nil, err
	}
	announce(args...)
	cmd := exec.Command(exe, args[1:]...)
	return cmd.CombinedOutput()
}

func announce(args ...string) {
	fmt.Println(shellInspect(args))
}

func shellInspect(args []string) string {
	fmtArgs := make([]string, len(args))
	for i, arg := range args {
		if strings.ContainsAny(arg, " \t'\"") {
			fmtArgs[i] = fmt.Sprintf("%q", arg)
		} else {
			fmtArgs[i] = arg
		}
	}
	return strings.Join(fmtArgs, " ")
}

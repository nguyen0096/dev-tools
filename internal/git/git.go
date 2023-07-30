package git

import (
	"bytes"
	"dev-tools/internal/config"
	"dev-tools/pkg/shell"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

func GetGitCommand() *cobra.Command {
	gitCmd := &cobra.Command{
		Use:   "git",
		Short: "Automate repository processes",
	}

	gitCmd.AddCommand(commitUpdate())

	return gitCmd
}

func commitUpdate() *cobra.Command {
	commitUpdateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update repo with predefined commit message format",
		Run: func(cmd *cobra.Command, args []string) {
			for _, r := range config.Cfg.Git.Repos {
				if err := os.Chdir(r.Path); err != nil {
					log.Printf("failed to change working dir. err: %v", err)
				}
				if output, err := gitAddAll(); err != nil {
					log.Printf("failed to run git add. err: %v. output: %s", err, output)
				}
				if output, err := gitCommit(r.MessageTemplate); err != nil {
					log.Printf("failed to run git commit. err: %v. output: %s", err, output)
				}
				if output, err := gitPush(); err != nil {
					log.Printf("failed to run git commit. err: %v. output: %s", err, output)
				}
			}
		},
	}

	return commitUpdateCmd
}

func gitAddAll() (string, error) {
	o, err := shell.ExecOutput("git", "add", ".")
	if err != nil {
		return string(o), err
	}
	return string(o), nil
}

func gitCommit(msgTmpl string) (string, error) {
	now := time.Now()
	_, w := now.ISOWeek()
	tmplVars := config.MessageVariables{
		Week:          w,
		OrdinalNumber: "001",
		Date:          now.Format("060102_150405"),
	}

	t := template.Must(template.New("git").Parse(msgTmpl))
	buf := bytes.NewBuffer(nil)
	t.Execute(buf, tmplVars)

	o, err := shell.ExecOutput("git", "commit", "-m", fmt.Sprintf(`"%s"`, buf.String()))
	if err != nil {
		return string(o), err
	}
	return string(o), nil
}

func gitPush() (string, error) {
	o, err := shell.ExecOutput("git", "push")
	if err != nil {
		return string(o), err
	}
	return string(o), nil
}

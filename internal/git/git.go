package git

import (
	"bytes"
	"dev-tools/config"
	"dev-tools/internal/domain"
	"dev-tools/pkg/shell"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/go-git/go-git/v5"
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
			for _, rcfg := range config.Cfg.Git.Repos {
				if err := updateRepo(rcfg); err != nil {
					log.Printf("failed to update repo [%s]. err: %s", rcfg.Path, err)
				}
			}
		},
	}

	return commitUpdateCmd
}

func updateRepo(rcfg config.Repository) error {
	log.Printf("=== repo [%s]", rcfg.Path)
	r, err := git.PlainOpen(rcfg.Path)
	if err != nil {
		err = fmt.Errorf("failed to open git folder. %w", err)
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		err = fmt.Errorf("failed to get worktree. %w", err)
		return err
	}

	stt, err := w.Status()
	if err != nil {
		err = fmt.Errorf("failed to check worktree status. %w", err)
		return err
	}
	if stt.IsClean() {
		log.Printf("repo is clean. skipped")
		return nil
	}

	if err := gitAddDot(rcfg.Path); err != nil {
		err = fmt.Errorf("failed to exec git add. %w", err)
		return err
	}

	now := time.Now()
	_, week := now.ISOWeek()

	dateFormat := rcfg.DateFormat
	if dateFormat == "" {
		dateFormat = "060102_150405"
	}

	tmplVars := domain.MessageVariables{
		Week: strconv.Itoa(week),
		Date: now.Format(dateFormat),
	}

	t := template.Must(template.New("git").Parse(rcfg.MessageTemplate))
	buf := bytes.NewBuffer(nil)
	t.Execute(buf, tmplVars)

	if _, err := w.Commit(buf.String(), &git.CommitOptions{}); err != nil {
		err = fmt.Errorf("failed to commit. %w", err)
		return err
	}

	if err := r.Push(&git.PushOptions{}); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Printf("failed to push. already up-to-date")
			return nil
		}
		err = fmt.Errorf("failed to push. %w", err)
		return err
	}

	log.Printf("made new commit and pushed to remote")
	return nil
}

func gitAddDot(path string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(path); err != nil {
		return err
	}
	if o, err := shell.ExecOutput("git", "add", "."); err != nil {
		err = fmt.Errorf("%w. %s", err, o)
		return err
	}
	if err := os.Chdir(wd); err != nil {
		return err
	}
	return nil
}

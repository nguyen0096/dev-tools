package git

import (
	"bytes"
	"dev-tools/config"
	"dev-tools/internal/domain"
	"fmt"
	"log"
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
		log.Printf("repo [%s] is clean. skipped", rcfg.Path)
		return nil
	}

	if _, err := w.Add("."); err != nil {
		err = fmt.Errorf("failed to add files. %w", err)
		return err
	}

	now := time.Now()
	week, _ := now.ISOWeek()

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

	if _, err := w.Commit(buf.String(), nil); err != nil {
		err = fmt.Errorf("failed to commit. %w", err)
		return err
	}

	return nil
}

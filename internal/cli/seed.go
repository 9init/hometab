package cli

import (
	"github.com/spf13/cobra"

	"github.com/9init/hometab/internal/session"
)

func (a *app) copyCmd() *cobra.Command {
	return a.seedCmd("copy",
		"Copy files from your real home into a session (one-time snapshot)",
		"  hometab copy work .config/Claude .claude/projects",
		a.store.Copy)
}

func (a *app) linkCmd() *cobra.Command {
	return a.seedCmd("link",
		"Symlink files from your real home into a session (shared live)",
		"  hometab link personal .gitconfig .ssh",
		a.store.Link)
}

func (a *app) seedCmd(use, short, example string, seed func(session.Session, string) error) *cobra.Command {
	return &cobra.Command{
		Use:               use + " SESSION PATH...",
		Short:             short,
		Example:           example,
		Args:              cobra.MinimumNArgs(2),
		ValidArgsFunction: a.completeSession,
		RunE: func(_ *cobra.Command, args []string) error {
			s, err := a.store.Session(args[0])
			if err != nil {
				return err
			}
			for _, path := range args[1:] {
				if err := seed(s, path); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

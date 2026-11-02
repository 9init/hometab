package cli

import (
	"github.com/spf13/cobra"

	"github.com/9init/hometab/internal/launch"
	"github.com/9init/hometab/internal/session"
)

type app struct {
	store *session.Store
	exec  func(*launch.Cmd) error
}

func Execute(version string) error {
	store, err := session.Open()
	if err != nil {
		return err
	}
	a := &app{store: store, exec: (*launch.Cmd).Exec}
	return a.rootCmd(version).Execute()
}

func (a *app) rootCmd(version string) *cobra.Command {
	root := &cobra.Command{
		Use:   "hometab",
		Short: "Run any executable with its own home directory, kept per session",
		Long: `hometab runs a program with HOME and XDG dirs pointed at a per-session
directory, so the same app can run side by side with separate logins and
settings. Using the same session name again returns to the same state.

Sessions live in ~/.local/share/hometab (override with HOMETAB_DIR).`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(
		a.runCmd(), a.shellCmd(),
		a.lsCmd(), a.mvCmd(), a.rmCmd(),
		a.copyCmd(), a.linkCmd(),
		a.genManCmd(),
	)
	return root
}

func (a *app) launch(name string, argv []string, opts launch.Options) error {
	s, err := a.store.Session(name)
	if err != nil {
		return err
	}
	if err := s.Prepare(); err != nil {
		return err
	}
	cmd, err := launch.Command(s, a.store.RealHome, argv, opts)
	if err != nil {
		return err
	}
	return a.exec(cmd)
}

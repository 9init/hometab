package cli

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/9init/hometab/internal/launch"
)

func (a *app) runCmd() *cobra.Command {
	var opts launch.Options
	cmd := &cobra.Command{
		Use:   "run [flags] SESSION [--] COMMAND [ARGS...]",
		Short: "Run a command inside a session",
		Example: `  hometab run work claude-desktop
  hometab run test ./mybin --verbose
  hometab run --share ~/work sep claude-desktop`,
		Args:              cobra.MinimumNArgs(2),
		ValidArgsFunction: a.completeSessionThenCommand,
		RunE: func(_ *cobra.Command, args []string) error {
			argv := args[1:]
			// SetInterspersed(false) makes cobra keep the separator.
			if argv[0] == "--" {
				argv = argv[1:]
			}
			if len(argv) == 0 {
				return errors.New("missing command to run")
			}
			return a.launch(args[0], argv, opts)
		},
	}
	cmd.Flags().SetInterspersed(false)
	addSandboxFlags(cmd, &opts)
	return cmd
}

func (a *app) shellCmd() *cobra.Command {
	var opts launch.Options
	cmd := &cobra.Command{
		Use:               "shell [flags] SESSION",
		Short:             "Open $SHELL inside a session",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: a.completeSession,
		RunE: func(_ *cobra.Command, args []string) error {
			shell := os.Getenv("SHELL")
			if shell == "" {
				shell = "/bin/sh"
			}
			return a.launch(args[0], []string{shell}, opts)
		},
	}
	addSandboxFlags(cmd, &opts)
	return cmd
}

func addSandboxFlags(cmd *cobra.Command, opts *launch.Options) {
	cmd.Flags().BoolVar(&opts.Sandbox, "sandbox", false, "run under bwrap: the session is mounted over your real home, which the app cannot see")
	cmd.Flags().StringArrayVar(&opts.Shares, "share", nil, "path to keep visible inside the sandbox (repeatable, implies --sandbox)")
}

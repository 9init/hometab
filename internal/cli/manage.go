package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func (a *app) lsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List sessions",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sessions, err := a.store.List()
			if err != nil {
				return err
			}
			for _, s := range sessions {
				fmt.Fprintln(cmd.OutOrStdout(), s.Name)
			}
			return nil
		},
	}
}

func (a *app) mvCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "mv SESSION NEW_NAME",
		Aliases:           []string{"rename"},
		Short:             "Rename a session (close apps running in it first)",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: a.completeSession,
		RunE: func(_ *cobra.Command, args []string) error {
			return a.store.Rename(args[0], args[1])
		},
	}
}

func (a *app) rmCmd() *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:               "rm SESSION",
		Short:             "Delete a session and all its data (logins, settings)",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: a.completeSession,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes && !confirm(cmd, fmt.Sprintf("Delete session %q and all its data?", args[0])) {
				return nil
			}
			return a.store.Remove(args[0])
		},
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "do not ask for confirmation")
	return cmd
}

func confirm(cmd *cobra.Command, question string) bool {
	fmt.Fprintf(cmd.OutOrStdout(), "%s [y/N] ", question)
	answer, _ := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
	return strings.EqualFold(strings.TrimSpace(answer), "y")
}

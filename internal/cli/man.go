package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func (a *app) genManCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "gen-man DIR",
		Short:  "Write man pages to a directory",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := os.MkdirAll(args[0], 0o755); err != nil {
				return err
			}
			root := cmd.Root()
			root.DisableAutoGenTag = true
			return doc.GenManTree(root, &doc.GenManHeader{Title: "HOMETAB", Section: "1"}, args[0])
		},
	}
}

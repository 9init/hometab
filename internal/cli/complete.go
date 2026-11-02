package cli

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func (a *app) completeSession(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveDefault
	}
	sessions, err := a.store.List()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	names := make([]string, len(sessions))
	for i, s := range sessions {
		names[i] = s.Name
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func (a *app) completeSessionThenCommand(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch {
	case len(args) == 0:
		return a.completeSession(cmd, args, toComplete)
	case len(args) == 1 && !strings.Contains(toComplete, "/"):
		return commandsOnPath(toComplete), cobra.ShellCompDirectiveNoFileComp
	default:
		return nil, cobra.ShellCompDirectiveDefault
	}
}

func commandsOnPath(prefix string) []string {
	seen := map[string]bool{}
	var names []string
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			name := e.Name()
			if seen[name] || !strings.HasPrefix(name, prefix) {
				continue
			}
			if isExecutable(filepath.Join(dir, name)) {
				seen[name] = true
				names = append(names, name)
			}
		}
	}
	return names
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir() && info.Mode()&0o111 != 0
}

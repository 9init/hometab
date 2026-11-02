package launch

import (
	"fmt"
	"os"
	"path/filepath"
)

func sandboxMounts(sessionDir, realHome string, shares []string) ([]string, error) {
	args := []string{
		"--dev-bind", "/", "/",
		"--bind", sessionDir, realHome,
	}
	for _, p := range shares {
		abs, err := filepath.Abs(p)
		if err != nil {
			return nil, err
		}
		if abs == realHome {
			return nil, fmt.Errorf("cannot share the whole home %s", abs)
		}
		if _, err := os.Stat(abs); err != nil {
			return nil, fmt.Errorf("share: %w", err)
		}
		args = append(args, "--bind", abs, abs)
	}
	return append(args, "--unshare-pid", "--proc", "/proc"), nil
}

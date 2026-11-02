package session

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (st *Store) Copy(s Session, path string) error {
	src, dst, err := st.seedPaths(s, path)
	if err != nil {
		return err
	}
	// -T copies onto dst, not into it, so re-importing doesn't nest.
	if out, err := exec.Command("cp", "-a", "-T", src, dst).CombinedOutput(); err != nil {
		return fmt.Errorf("copy %s: %w: %s", path, err, strings.TrimSpace(string(out)))
	}
	return nil
}

func (st *Store) Link(s Session, path string) error {
	src, dst, err := st.seedPaths(s, path)
	if err != nil {
		return err
	}
	if err := os.Symlink(src, dst); err != nil {
		return fmt.Errorf("link %s: %w", path, err)
	}
	return nil
}

func (st *Store) seedPaths(s Session, path string) (src, dst string, err error) {
	rel, err := st.relToRealHome(path)
	if err != nil {
		return "", "", err
	}
	src, dst = filepath.Join(st.RealHome, rel), filepath.Join(s.Dir, rel)
	if _, err := os.Lstat(src); err != nil {
		return "", "", err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return "", "", err
	}
	return src, dst, nil
}

func (st *Store) relToRealHome(path string) (string, error) {
	rel := path
	if filepath.IsAbs(path) {
		var err error
		if rel, err = filepath.Rel(st.RealHome, path); err != nil {
			return "", err
		}
	}
	rel = filepath.Clean(rel)
	if rel == "." || rel == ".." || strings.HasPrefix(rel, "../") {
		return "", fmt.Errorf("%q must be inside %s", path, st.RealHome)
	}
	return rel, nil
}

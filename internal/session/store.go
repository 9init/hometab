package session

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
)

var (
	ErrNotFound = errors.New("session not found")
	ErrExists   = errors.New("session already exists")
)

var validName = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

type Store struct {
	Dir      string
	RealHome string
}

func Open() (*Store, error) {
	// Not $HOME: inside a session that points at the session itself.
	u, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("look up current user: %w", err)
	}
	dir := os.Getenv("HOMETAB_DIR")
	if dir == "" {
		dir = filepath.Join(u.HomeDir, ".local", "share", "hometab")
	}
	return &Store{Dir: dir, RealHome: u.HomeDir}, nil
}

func (st *Store) Session(name string) (Session, error) {
	if !validName.MatchString(name) {
		return Session{}, fmt.Errorf("invalid session name %q (allowed: letters, digits, _ -)", name)
	}
	return Session{Name: name, Dir: filepath.Join(st.Dir, name)}, nil
}

func (st *Store) List() ([]Session, error) {
	entries, err := os.ReadDir(st.Dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	var sessions []Session
	for _, e := range entries {
		if s, err := st.Session(e.Name()); err == nil && e.IsDir() {
			sessions = append(sessions, s)
		}
	}
	return sessions, nil
}

func (st *Store) Rename(from, to string) error {
	src, err := st.existing(from)
	if err != nil {
		return err
	}
	dst, err := st.Session(to)
	if err != nil {
		return err
	}
	// os.Rename silently replaces an empty directory.
	if dst.exists() {
		return fmt.Errorf("%s: %w", to, ErrExists)
	}
	return os.Rename(src.Dir, dst.Dir)
}

func (st *Store) Remove(name string) error {
	s, err := st.existing(name)
	if err != nil {
		return err
	}
	if err := os.RemoveAll(s.Dir); err == nil {
		return nil
	}
	// Some tools, like Go's module cache, leave directories read-only.
	if err := makeWritable(s.Dir); err != nil {
		return fmt.Errorf("remove %s: %w", name, err)
	}
	return os.RemoveAll(s.Dir)
}

func makeWritable(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		return os.Chmod(path, info.Mode().Perm()|0o700)
	})
}

func (st *Store) existing(name string) (Session, error) {
	s, err := st.Session(name)
	if err != nil {
		return Session{}, err
	}
	if !s.exists() {
		return Session{}, fmt.Errorf("%s: %w", name, ErrNotFound)
	}
	return s, nil
}

package session

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func newStore(t *testing.T) *Store {
	t.Helper()
	return &Store{Dir: t.TempDir(), RealHome: t.TempDir()}
}

func mustCreate(t *testing.T, st *Store, name string) Session {
	t.Helper()
	s, err := st.Session(name)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Prepare(); err != nil {
		t.Fatal(err)
	}
	return s
}

func names(t *testing.T, st *Store) []string {
	t.Helper()
	list, err := st.List()
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, s := range list {
		out = append(out, s.Name)
	}
	return out
}

func TestSessionRejectsInvalidNames(t *testing.T) {
	st := newStore(t)
	for _, bad := range []string{"", ".", "..", "../x", "a/b", "a b"} {
		if _, err := st.Session(bad); err == nil {
			t.Errorf("accepted %q", bad)
		}
	}
}

func TestPrepareCreatesXDGDirs(t *testing.T) {
	s := mustCreate(t, newStore(t), "w")
	for _, d := range XDGDirs {
		if _, err := os.Stat(filepath.Join(s.Dir, d.Path)); err != nil {
			t.Errorf("%s not created: %v", d.Path, err)
		}
	}
}

func TestList(t *testing.T) {
	st := newStore(t)
	if got := names(t, st); got != nil {
		t.Fatalf("empty store listed %v", got)
	}
	mustCreate(t, st, "b")
	mustCreate(t, st, "a")
	os.WriteFile(filepath.Join(st.Dir, "stray-file"), nil, 0o600)
	if got := names(t, st); !slices.Equal(got, []string{"a", "b"}) {
		t.Errorf("List = %v", got)
	}
}

func TestRename(t *testing.T) {
	st := newStore(t)
	a := mustCreate(t, st, "a")
	mustCreate(t, st, "taken")
	os.WriteFile(filepath.Join(a.Dir, "data"), []byte("1"), 0o600)

	if err := st.Rename("a", "taken"); !errors.Is(err, ErrExists) {
		t.Errorf("rename onto existing: %v", err)
	}
	if err := st.Rename("missing", "z"); !errors.Is(err, ErrNotFound) {
		t.Errorf("rename missing: %v", err)
	}
	if err := st.Rename("a", "../x"); err == nil {
		t.Error("accepted invalid new name")
	}
	if err := st.Rename("a", "b"); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(st.Dir, "b", "data")); string(b) != "1" {
		t.Error("data not moved")
	}
	if got := names(t, st); !slices.Equal(got, []string{"b", "taken"}) {
		t.Errorf("List = %v", got)
	}
}

func TestRemove(t *testing.T) {
	st := newStore(t)
	mustCreate(t, st, "a")
	if err := st.Remove("a"); err != nil {
		t.Fatal(err)
	}
	if err := st.Remove("a"); !errors.Is(err, ErrNotFound) {
		t.Errorf("second remove: %v", err)
	}
}

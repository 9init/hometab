package session

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	st := newStore(t)
	s := mustCreate(t, st, "s")
	os.MkdirAll(filepath.Join(st.RealHome, ".config/app"), 0o700)
	os.WriteFile(filepath.Join(st.RealHome, ".config/app/x"), []byte("1"), 0o600)

	for range 2 {
		if err := st.Copy(s, ".config/app"); err != nil {
			t.Fatal(err)
		}
	}
	if b, _ := os.ReadFile(filepath.Join(s.Dir, ".config/app/x")); string(b) != "1" {
		t.Error("file not copied")
	}
	if _, err := os.Stat(filepath.Join(s.Dir, ".config/app/app")); err == nil {
		t.Error("second copy nested the directory")
	}
}

func TestLink(t *testing.T) {
	st := newStore(t)
	s := mustCreate(t, st, "s")
	real := filepath.Join(st.RealHome, ".gitconfig")
	os.WriteFile(real, []byte("g"), 0o600)

	if err := st.Link(s, real); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.Readlink(filepath.Join(s.Dir, ".gitconfig")); got != real {
		t.Errorf("link points to %q", got)
	}
}

func TestSeedRejectsPathsOutsideHome(t *testing.T) {
	st := newStore(t)
	s := mustCreate(t, st, "s")
	for _, bad := range []string{"..", "../etc", "/etc/passwd", "missing"} {
		if err := st.Copy(s, bad); err == nil {
			t.Errorf("accepted %q", bad)
		}
	}
}

package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/9init/hometab/internal/launch"
	"github.com/9init/hometab/internal/session"
)

type harness struct {
	app      *app
	launched *launch.Cmd
	out      bytes.Buffer
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	h := &harness{}
	h.app = &app{
		store: &session.Store{Dir: t.TempDir(), RealHome: t.TempDir()},
		exec:  func(c *launch.Cmd) error { h.launched = c; return nil },
	}
	return h
}

func (h *harness) run(t *testing.T, stdin string, args ...string) error {
	t.Helper()
	root := h.app.rootCmd("test")
	root.SetArgs(args)
	root.SetOut(&h.out)
	root.SetIn(strings.NewReader(stdin))
	return root.Execute()
}

func TestRunPassesArgsToTheProgram(t *testing.T) {
	for _, args := range [][]string{
		{"run", "s", "sh", "-c", "true"},
		{"run", "s", "--", "sh", "-c", "true"},
	} {
		h := newHarness(t)
		if err := h.run(t, "", args...); err != nil {
			t.Fatalf("%v: %v", args, err)
		}
		if !slices.Equal(h.launched.Args, []string{"sh", "-c", "true"}) {
			t.Errorf("%v: launched %v", args, h.launched.Args)
		}
		if !slices.Contains(h.launched.Env, "HOME="+filepath.Join(h.app.store.Dir, "s")) {
			t.Errorf("%v: HOME not set to the session", args)
		}
	}
}

func TestRunRequiresACommand(t *testing.T) {
	h := newHarness(t)
	if err := h.run(t, "", "run", "s", "--"); err == nil {
		t.Error("ran with no command")
	}
}

func TestSessionLifecycle(t *testing.T) {
	h := newHarness(t)
	steps := [][]string{
		{"run", "a", "true"},
		{"mv", "a", "b"},
		{"ls"},
	}
	for _, args := range steps {
		if err := h.run(t, "", args...); err != nil {
			t.Fatalf("%v: %v", args, err)
		}
	}
	if got := strings.TrimSpace(h.out.String()); got != "b" {
		t.Errorf("ls after rename = %q", got)
	}

	if err := h.run(t, "n\n", "rm", "b"); err != nil {
		t.Fatal(err)
	}
	if sessions, _ := h.app.store.List(); len(sessions) != 1 {
		t.Error("rm deleted without confirmation")
	}
	if err := h.run(t, "", "rm", "-y", "b"); err != nil {
		t.Fatal(err)
	}
	if sessions, _ := h.app.store.List(); len(sessions) != 0 {
		t.Error("rm -y did not delete")
	}
}

func TestGenMan(t *testing.T) {
	h := newHarness(t)
	dir := t.TempDir()
	if err := h.run(t, "", "gen-man", dir); err != nil {
		t.Fatal(err)
	}
	for _, page := range []string{"hometab.1", "hometab-run.1"} {
		if _, err := os.Stat(filepath.Join(dir, page)); err != nil {
			t.Errorf("%s not generated", page)
		}
	}
}

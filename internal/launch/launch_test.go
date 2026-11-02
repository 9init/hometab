package launch

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"

	"github.com/9init/hometab/internal/session"
)

func testSession(t *testing.T) (session.Session, string) {
	t.Helper()
	return session.Session{Name: "w", Dir: t.TempDir()}, t.TempDir()
}

func TestEnvironmentOverridesHomeAndKeepsTheRest(t *testing.T) {
	s, _ := testSession(t)
	env := environment([]string{"HOME=/real", "XDG_CONFIG_HOME=/real/.config", "DISPLAY=:0"}, s, "/seen")

	for _, want := range []string{"HOME=/seen", "XDG_CONFIG_HOME=/seen/.config", "XDG_CACHE_HOME=/seen/.cache", "DISPLAY=:0", "HOMETAB_SESSION=w"} {
		if !slices.Contains(env, want) {
			t.Errorf("missing %s in %v", want, env)
		}
	}
	for _, leaked := range []string{"HOME=/real", "XDG_CONFIG_HOME=/real/.config"} {
		if slices.Contains(env, leaked) {
			t.Errorf("%s leaked", leaked)
		}
	}
}

func TestCommandDirect(t *testing.T) {
	s, realHome := testSession(t)
	cmd, err := Command(s, realHome, []string{"sh", "-c", "true"}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(cmd.Path) != "sh" || !slices.Equal(cmd.Args, []string{"sh", "-c", "true"}) {
		t.Errorf("got %s %v", cmd.Path, cmd.Args)
	}
	if !slices.Contains(cmd.Env, "HOME="+s.Dir) {
		t.Error("HOME is not the session dir")
	}
}

func TestCommandSandboxed(t *testing.T) {
	if _, err := exec.LookPath("bwrap"); err != nil {
		t.Skip("bwrap not installed")
	}
	s, realHome := testSession(t)
	work := filepath.Join(realHome, "work")
	os.Mkdir(work, 0o700)

	cmd, err := Command(s, realHome, []string{"sh", "-c", "true"}, Options{Shares: []string{work}})
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(cmd.Path) != "bwrap" {
		t.Errorf("not run through bwrap: %s", cmd.Path)
	}
	if !slices.Contains(cmd.Env, "HOME="+realHome) {
		t.Error("HOME should be the real home path, which bwrap backs with the session")
	}
	sep := slices.Index(cmd.Args, "--")
	if sep < 0 || filepath.Base(cmd.Args[sep+1]) != "sh" || !slices.Equal(cmd.Args[sep+2:], []string{"-c", "true"}) {
		t.Errorf("program not passed after --: %v", cmd.Args)
	}
}

func TestSandboxMountsSessionBeforeShares(t *testing.T) {
	s, realHome := testSession(t)
	work := filepath.Join(realHome, "work")
	os.Mkdir(work, 0o700)

	args, err := sandboxMounts(s.Dir, realHome, []string{work})
	if err != nil {
		t.Fatal(err)
	}
	home, share := slices.Index(args, s.Dir), slices.Index(args, work)
	if home < 0 || args[home+1] != realHome || share < home {
		t.Errorf("session must be mounted over home before shares: %v", args)
	}
}

func TestSandboxRejectsBadShares(t *testing.T) {
	s, realHome := testSession(t)
	for _, bad := range []string{realHome, filepath.Join(realHome, "missing")} {
		if _, err := sandboxMounts(s.Dir, realHome, []string{bad}); err == nil {
			t.Errorf("accepted share %q", bad)
		}
	}
}

package launch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/9init/hometab/internal/session"
)

type Options struct {
	Sandbox bool
	Shares  []string
}

func (o Options) sandboxed() bool {
	return o.Sandbox || len(o.Shares) > 0
}

type Cmd struct {
	Path string
	Args []string
	Env  []string
}

func (c *Cmd) Exec() error {
	return syscall.Exec(c.Path, c.Args, c.Env)
}

func Command(s session.Session, realHome string, argv []string, opts Options) (*Cmd, error) {
	program, err := resolve(argv[0])
	if err != nil {
		return nil, err
	}
	if !opts.sandboxed() {
		return &Cmd{
			Path: program,
			Args: argv,
			Env:  environment(os.Environ(), s, s.Dir),
		}, nil
	}

	bwrap, err := exec.LookPath("bwrap")
	if err != nil {
		return nil, fmt.Errorf("sandbox needs bubblewrap: %w", err)
	}
	mounts, err := sandboxMounts(s.Dir, realHome, opts.Shares)
	if err != nil {
		return nil, err
	}
	args := append([]string{"bwrap"}, mounts...)
	args = append(args, "--", program)
	return &Cmd{
		Path: bwrap,
		Args: append(args, argv[1:]...),
		Env:  environment(os.Environ(), s, realHome),
	}, nil
}

func resolve(name string) (string, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", err
	}
	return filepath.Abs(path)
}

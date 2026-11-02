package session

import (
	"fmt"
	"os"
	"path/filepath"
)

type Session struct {
	Name string
	Dir  string
}

type XDGDir struct {
	Env  string
	Path string
}

var XDGDirs = []XDGDir{
	{Env: "XDG_CONFIG_HOME", Path: ".config"},
	{Env: "XDG_DATA_HOME", Path: ".local/share"},
	{Env: "XDG_STATE_HOME", Path: ".local/state"},
	{Env: "XDG_CACHE_HOME", Path: ".cache"},
}

func (s Session) Prepare() error {
	for _, d := range XDGDirs {
		if err := os.MkdirAll(filepath.Join(s.Dir, d.Path), 0o700); err != nil {
			return fmt.Errorf("prepare session %s: %w", s.Name, err)
		}
	}
	return nil
}

func (s Session) exists() bool {
	_, err := os.Stat(s.Dir)
	return err == nil
}

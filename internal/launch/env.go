package launch

import (
	"path/filepath"
	"strings"

	"github.com/9init/hometab/internal/session"
)

func environment(base []string, s session.Session, home string) []string {
	overrides := map[string]string{
		"HOME":            home,
		"HOMETAB_SESSION": s.Name,
	}
	for _, d := range session.XDGDirs {
		overrides[d.Env] = filepath.Join(home, d.Path)
	}

	env := make([]string, 0, len(base)+len(overrides))
	for _, kv := range base {
		key, _, _ := strings.Cut(kv, "=")
		if _, replaced := overrides[key]; !replaced {
			env = append(env, kv)
		}
	}
	for key, value := range overrides {
		env = append(env, key+"="+value)
	}
	return env
}

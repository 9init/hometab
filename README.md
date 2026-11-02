# hometab

hometab runs any program with its own home directory. Each session has a
name, and the program keeps its logins, settings and data there. Use the same
name again and you're back where you left off.

This lets you run several copies of the same app side by side, each with its
own account, such as a work and a personal profile of a desktop app that
normally allows only one. It also keeps experiments and test runs from
touching your real config.

## Usage

```sh
hometab run work claude-desktop
hometab run personal claude-desktop
```

Two windows, two accounts. It works the same for any app or binary:

```sh
hometab run test ./mybin --verbose
hometab shell personal
```

Manage sessions:

```sh
hometab ls
hometab mv personal home
hometab rm home
```

Bring files over from your real home, either as a copy or as a live link:

```sh
hometab copy work .config/Claude
hometab link personal .gitconfig .ssh
```

To hide your real home from the app completely, use `--sandbox` (needs
`bwrap`), and `--share` for anything it should still see. Flags go before the
session name; everything after the program name is passed to the program:

```sh
hometab run --share ~/work sep claude-desktop
```

Sessions are stored in `~/.local/share/hometab`.

## Limitations

- Linux only.
- Without `--sandbox`, this separates data, not access. The app runs as you and
  can still open files in your real home by full path or through a file picker.
- Apps that find their home some other way than `$HOME` and the XDG variables
  will use your real home. `--sandbox` handles these.
- A session starts empty: no dotfiles, git config or SSH keys until you `copy`
  or `link` them.
- Some apps keep data in more than one place. Claude Desktop, for example,
  stores chat history in `.claude/projects` next to `.config/Claude`, so copy
  both.
- Inside `--sandbox`, links made with `link` point to your hidden real home and
  stop working. Use `--share` instead, and stick to one mode per session.
- Renaming or removing a session while an app is running in it will confuse the
  app. Quit it first.

## Install

```sh
make install
```

This installs the binary to `~/.local/bin` and the man pages, so `man hometab`
and `man hometab-run` work.

## License

Apache-2.0

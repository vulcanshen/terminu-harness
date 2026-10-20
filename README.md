# terminu-harness

[![License](https://img.shields.io/badge/license-GPL--3.0-blue)](LICENSE)

**Language**: English · [繁體中文](README-zh_TW.md)

**A Claude Code plugin that builds TUI apps with [terminu design](https://github.com/vulcanshen/terminu) (tdp) and [terminu-go](https://github.com/vulcanshen/terminu-go), from what you want to a finished app.**

## What it does

Type `/terminu` and say what you want to do:

- **Make a new app**: it asks what the app is for, one question at a time (only what you need; screens, keys and file formats it decides by tdp), writes a spec and gets going: it builds the app with terminu-go from the tdp project templates piece by piece, runs the tests and opens each piece in a terminal of its own, and shows you the app when it is done. Where tdp says nothing or terminu-go can't do it yet, it tells you and drafts an issue for tdp or terminu-go; once you agree it files it, works around it in the meantime and comes back when they are updated.
- **Carry on with an app**: type `/terminu` in the app's folder; it first checks the issues it is waiting for, new versions of tdp and terminu-go and the workarounds left in the code, tells you, then picks up where it left off.
- **Review changes**: a separate, read-only reviewer looks at a range of changes against tdp and the app's spec; every problem names where it comes from (a tdp item or a line of the spec), and it is fixed and looked at again.
- **Release**: runs the checks, reviews what changed since the last release, writes the CHANGELOG the tdp way, and tags and pushes only once you agree.
- **Report to tdp or terminu-go**: it drafts an issue for tdp or terminu-go and files it once you agree.

It needs no tmux or graphical desktop, and it asks before anything that leaves your machine (issues, GitHub repositories, releases).

## Requirements

- [Claude Code](https://claude.com/claude-code)
- macOS or Linux
- [Go](https://go.dev/dl/): the first time the plugin starts, it builds its terminal from source; new apps are written in Go
- [GitHub CLI](https://cli.github.com) (`gh`), signed in, to file issues; without it you get the issue's text to post yourself

## Install

In Claude Code:

```
/plugin marketplace add vulcanshen/terminu-harness
/plugin install terminu@terminu-harness
```

The first start builds the terminal into `~/.terminu-harness/bin`, which takes a little while. If Claude Code reports the plugin's MCP server failed to start, run `/mcp` to reconnect once the build is done.

To remove it, `/plugin uninstall terminu@terminu-harness`, then delete `~/.terminu-harness`.

## Usage

```
/terminu
```

It asks what you want to do; in the folder of an app it has been making, it carries on with it. You can say it right away: `/terminu make an app that …`. The steps are also commands of their own:

| Command | What it does |
|---|---|
| `/terminu:newapp [what the app is for]` | make a new app, or carry on with the one in this folder |
| `/terminu:review [range]` | review a range of changes; without one, the uncommitted changes and the commits not pushed yet |
| `/terminu:release [version]` | release: checks, a review of what changed since the last release, the CHANGELOG, and a tag and push once you agree |
| `/terminu:reportgap [what is missing]` | draft an issue for tdp or terminu-go and file it once you agree |

## Where things go

| What | Where |
|---|---|
| A new app's progress and issue drafts | `.local/terminu/` in the app's repo (git ignores it) |
| A new app's spec | `docs/spec.md` in the app's repo |
| tdp files fetched | `~/.terminu-harness/cache/tdp/` |
| The built terminal | `~/.terminu-harness/bin/` |

| Variable | What it does |
|---|---|
| `TERMINU_HARNESS__HOME` | Keep everything somewhere other than `~/.terminu-harness` |
| `TERMINU_HARNESS__TDP_REPO` | Read tdp from a local clone of terminu rather than GitHub |

## Limits

- No Windows for the plugin itself (the apps it makes can ship for Windows).
- Waiting on an issue is not a running session: once tdp or terminu-go answered, type `/terminu` in the app's folder again.

## Related

- [terminu](https://github.com/vulcanshen/terminu): terminu design (tdp) and its checklist
- [terminu-go](https://github.com/vulcanshen/terminu-go): the Go lib for terminu family apps
- [CHANGELOG](CHANGELOG.md)
- [docs/dev-remarks.md](docs/dev-remarks.md): how it works, for developers

## License

[GPL-3.0](LICENSE)

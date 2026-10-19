# terminu-harness

[![License](https://img.shields.io/badge/license-GPL--3.0-blue)](LICENSE)

**Language**: English · [繁體中文](README-zh_TW.md)

**A Claude Code plugin that builds TUI apps with [terminu design](https://github.com/vulcanshen/terminu) (tdp), from what you want to a finished app, and checks how many terminu genes an app has.** Every item on the tdp checklist an app meets is one gene.

## What it does

Type `/terminu` and say what you want to do:

- **Make a new app**: it asks what the app is for, one question at a time, writes a spec for you to confirm, then builds the app with [terminu-go](https://github.com/vulcanshen/terminu-go) from the tdp project templates, showing you each piece as it runs. Where tdp says nothing or terminu-go can't do it yet, it tells you and drafts an issue for tdp or terminu-go; once you agree it files it, works around it in the meantime and comes back when they are updated. When the app is done it judges it and gives you the report card.
- **Carry on with an app**: it picks up where it left off, checks the issues it is waiting for and upgrades when tdp or terminu-go has a new version.
- **Judge an app**: it opens the app in a terminal of its own, presses keys and looks at the screen the way each checklist item says, reads the repo when there is one, and writes a report card: whether the app **passed** (nothing violated), is **pending** (every violation is listed in its `docs/<app>-terminu-fix.md`, waiting on tdp or terminu-go) or **failed** (a violation is listed nowhere); every verdict (meets, departs, violates, not applicable, can't tell) with what it saw; the gene count as `tdp-0.7.4 860 gene`, to follow its progress; and what the app does that tdp does not cover yet. Any TUI app, not only terminu family members, with or without its source.

It needs no tmux or graphical desktop, it asks before anything that leaves your machine (issues, GitHub repositories, releases), and judging leaves the app's repo as it is.

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

It asks what you want to do. You can say it right away: `/terminu make an app that …`. The steps are also commands of their own:

| Command | What it does |
|---|---|
| `/terminu:new-app [what the app is for]` | make a new app, or carry on with the one in this folder |
| `/terminu:judge <how the app starts> [the app's repo] [tdp version]` | judge an app; without a repo, items that need the source are "can't tell"; without a version, the newest tdp |
| `/terminu:report-gap [what is missing]` | draft an issue for tdp or terminu-go and file it once you agree |

A whole judging run splits the checklist among ten judges working at the same time; it takes a while and uses many tokens.

## Where things go

| What | Where |
|---|---|
| Report cards | `~/.terminu-harness/judge/<app>/<time>-tdp-<version>/`: `report.md`, every verdict, the evidence |
| A new app's progress and issue drafts | `.local/terminu/` in the app's repo (git ignores it) |
| A new app's spec | `docs/spec.md` in the app's repo |
| tdp files fetched | `~/.terminu-harness/cache/tdp/` |
| The built terminal | `~/.terminu-harness/bin/` |

| Variable | What it does |
|---|---|
| `TERMINU_HARNESS__HOME` | Keep everything somewhere other than `~/.terminu-harness` |
| `TERMINU_HARNESS__OUT` | Put report cards somewhere else |
| `TERMINU_HARNESS__TDP_REPO` | Read tdp from a local clone of terminu rather than GitHub |

## Limits

- No Windows for the plugin itself (the apps it makes can ship for Windows).
- Waiting on an issue is not a running session: come back with `/terminu` and pick "carry on" once tdp or terminu-go answered.
- The terminal sends no mouse events: items that need a mouse are judged "can't tell".
- Times measured on the screen (animations, how long a toast stays) are judged within 33ms until tdp says how close is close enough.

## Related

- [terminu](https://github.com/vulcanshen/terminu): terminu design (tdp) and its checklist
- [terminu-go](https://github.com/vulcanshen/terminu-go): the Go lib for terminu family apps
- [CHANGELOG](CHANGELOG.md)
- [docs/dev-remarks.md](docs/dev-remarks.md): how it works, for developers

## License

[GPL-3.0](LICENSE)

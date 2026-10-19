# terminu-harness

[![License](https://img.shields.io/badge/license-GPL--3.0-blue)](LICENSE)

**Language**: English · [繁體中文](README-zh_TW.md)

**A Claude Code plugin that builds TUI apps with [terminu design](https://github.com/vulcanshen/terminu) (tdp), from what you want to a finished app, and checks how many terminu genes an app has.** Every item on the tdp checklist an app meets is one gene.

## What it does

Type `/terminu` and say what you want to do:

- **Make a new app**: it asks what the app is for, one question at a time, writes a spec for you to confirm, then builds the app with [terminu-go](https://github.com/vulcanshen/terminu-go) from the tdp project templates, showing you each piece as it runs. Where tdp says nothing or terminu-go can't do it yet, it tells you and drafts an issue for tdp or terminu-go; once you agree it files it, works around it in the meantime and comes back when they are updated. When the app is done it judges it and gives you the report card.
- **Carry on with an app**: it picks up where it left off, checks the issues it is waiting for and upgrades when tdp or terminu-go has a new version.
- **Judge an app**: it opens the app in a terminal of its own, presses keys and looks at the screen the way each checklist item says, reads the repo when there is one, and writes a report card: whether the app **passed** (nothing violated) or **failed** (anything violated, whether the app or terminu-go has to fix it; a departure the app wrote down with its reason does not fail it); every verdict (meets, departs, violates, not applicable, can't tell) with what it saw; the gene count as `tdp-0.8.2 860 gene`, to follow its progress; and what the app does that tdp does not cover yet. Any TUI app, not only terminu family members, with or without its source.

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

### Reading the report card

tdp says what to check; how to judge is the harness's:

- **Outcome**: any item violated and the app **failed**, whether the app has to fix it or it waits for tdp or terminu-go;
  nothing violated and it **passed**. Violations come in two lists: what the app has to fix, and what its
  `docs/<app>-terminu-fix.md` lists as waiting for tdp or terminu-go.
- **Verdicts**:

| Verdict | Meaning |
|---|---|
| meets | it does what the item says |
| departs | it does not, but the app's `docs/dev-remarks.md` names the item and the reason under departures from tdp; it does not fail the app and comes first in the report |
| violates | it does not, and it is not written down as a departure; when the app has none of something every app must have, each of its details is violated too |
| not applicable | the app has no feature the item is about, and tdp does not require one; the reason says how that is known. Items for every app are never not applicable |
| can't tell | it cannot be checked, and the reason says why: tdp is unclear or terminu-go stands in the way (reported to terminu), or something else (say, no repo; you are asked) |

- **Genes**: the items the app meets, written `tdp-<version> <n> gene` and tied to that tdp version. They show progress;
  apps with more features have more, so the numbers of two apps do not compare.
- A time measured on the screen meets the item within two frames (about 34ms) of what it says.
- Once an app has been judged with a tdp and terminu-go version, judging it again redoes only what changed since, and keeps
  the rest from the last report card.

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

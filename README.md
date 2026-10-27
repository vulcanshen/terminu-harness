# terminu-harness

[![License](https://img.shields.io/badge/license-GPL--3.0-blue)](LICENSE)

**Language**: English · [繁體中文](README-zh_TW.md)

**A Claude Code plugin that builds TUI apps with [terminu design](https://github.com/vulcanshen/terminu) (tdp) and [terminu-go](https://github.com/vulcanshen/terminu-go), from what you want to a finished app. You think about the features; the screens and keys follow tdp.**

## What it does

Type `/terminu` and say what you want to do:

- **Make a new app**: it asks what the app is for, one question at a time (only what you need; screens, keys and file formats it decides by tdp), writes your words down, matches each request to tdp and gets going: it builds the app with terminu-go from the tdp project templates one request at a time, runs the tests, opens each in a terminal of its own, writes down how it was built and which tests back it, and shows you the app when it is done.
- **When tdp has no answer, it does not make one up**: every decision the app makes has its ground in tdp, or lies where tdp says "the app may decide for itself". Where tdp says nothing, is unclear, or terminu-go can't do it yet, it tells you and drafts an issue for tdp or terminu-go; once you agree it files it, works around it in the meantime and comes back when they are updated.
- **Carry on, add requests, follow new versions**: type `/terminu` in the app's folder; it first checks for a newer terminu-go (every tdp version has a terminu-go version of its own), answers to the issues it is waiting for and the workarounds left in the code, and tells you; upgrading and new requests carry on from there.
- **Take over a terminu app made without it**: an app built with tdp and terminu-go before the plugin was added gets its records from its README (what it does) and git history (how each part was built and which tests back it), and goes the same way from then on.
- **Review**: a separate, read-only reviewer judges how each request was built, from the app's records, against tdp, without reading the code; the checks only you can do (say, ones that need a Windows machine) you confirm again each time. What does not follow tdp it fixes by itself, in the records or in the code; it asks you only whether to report something to tdp.
- **Release**: runs the checks, reviews, tells you what is still open and which requests are not built, writes the CHANGELOG the tdp way, and tags and pushes only once you agree.
- **Report to tdp or terminu-go**: it drafts an issue for tdp or terminu-go and files it once you agree.

It needs no tmux or graphical desktop, and it asks before anything that leaves your machine (issues, GitHub repositories, releases).

## Requirements

- [Claude Code](https://claude.com/claude-code)
- macOS or Linux
- [Go](https://go.dev/dl/): the first time the plugin starts, it builds its terminal and checks from source; apps are written in Go
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

It asks what you want to do; in the folder of an app it keeps, it carries on with it; in a terminu app it does not keep yet, it offers to take it over. You can say it right away: `/terminu make an app that …`. The steps are also commands of their own:

| Command | What it does |
|---|---|
| `/terminu:newapp [what the app is for]` | make a new app, or carry on, add requests, follow new versions |
| `/terminu:analysis` | take over a terminu app made without it: records from its README and git history |
| `/terminu:review [task IDs]` | judge the app's records against tdp; all of them without IDs |
| `/terminu:release [version]` | release: checks, a review, the CHANGELOG, and a tag and push once you agree |
| `/terminu:reportgap [what is missing]` | draft an issue for tdp or terminu-go and file it once you agree |

In an app it keeps, a commit that changes the code without updating the records (`harness-assets/work.yaml`) is stopped, and Claude is told to bring the records up to date; when the code changed but what the app does did not, the commit message ends with a line `Terminu-Work: unchanged`.

## Where things go

| What | Where |
|---|---|
| The app's records: your requests (`idea.md`), the tdp each follows (`spec.yaml`), how each was built and which tests back it (`work.yaml`), the last review (`review.yaml`) | `harness-assets/` in the app's repo |
| The rules for Claude writing code in the repo | `.claude/rules/terminu.md` in the app's repo |
| tdp files fetched | `~/.terminu-harness/cache/tdp/` |
| The built terminal and checks | `~/.terminu-harness/bin/` |

| Variable | What it does |
|---|---|
| `TERMINU_HARNESS__HOME` | Keep everything somewhere other than `~/.terminu-harness` |
| `TERMINU_HARNESS__TDP_REPO` | Read tdp from a local clone of terminu rather than GitHub |

## Limits

- **Review judges only what the records say**: what the code does but `harness-assets/work.yaml` does not say, review cannot see, and the harness does not vouch for.
- Commits you make yourself in a terminal are not stopped; the next review sees them.
- No Windows for the plugin itself (the apps it makes can ship for Windows).
- Waiting on an issue is not a running session: once tdp or terminu-go answered, type `/terminu` in the app's folder again.

## Related

- [terminu](https://github.com/vulcanshen/terminu): terminu design (tdp) and its checklist
- [terminu-go](https://github.com/vulcanshen/terminu-go): the Go lib for terminu family apps
- [CHANGELOG](CHANGELOG.md)
- [docs/dev-remarks.md](docs/dev-remarks.md): how it works, for developers

## License

[GPL-3.0](LICENSE)

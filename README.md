# terminu-harness

[![License](https://img.shields.io/badge/license-GPL--3.0-blue)](LICENSE)

**Language**: English · [繁體中文](README-zh_TW.md)

**A Claude Code plugin for TUI apps made with [terminu design](https://github.com/vulcanshen/terminu) (tdp).** It judges how many terminu genes an app has: every item on the tdp checklist the app meets is one gene.

## What it does

`/terminu:judge` opens the app in a terminal of its own, reads its help, presses keys and looks at the screen, the way each checklist item says to check it. Items that need the source or other files are judged when the app's repo is at hand. It writes a report card:

- the verdict on every item (meets, departs, violates, not applicable, or can't tell) with what it saw: the screen as text, recordings of animations, colours, pictures
- the gene count for that version of tdp, as `tdp-0.7.4 860 gene`
- what the app does that tdp does not cover yet, to report back to tdp

It judges any TUI app, not only terminu family members, whatever it is written in, with or without its source. It needs no tmux or graphical desktop, and it leaves the app's repo as it is.

## Requirements

- [Claude Code](https://claude.com/claude-code)
- macOS or Linux
- [Go](https://go.dev/dl/): the first time the plugin starts, it builds its terminal from source

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
/terminu:judge <how the app starts> [the app's repo] [tdp version]
```

For example `/terminu:judge demou ~/src/demou v0.7.4`. Without a repo, items that need the source are judged "can't tell"; without a version, it judges by the newest tdp.

It first learns the app (its help, `Space`, `?`, every screen and popup), then splits the checklist among ten judges working at the same time, and tells you the gene count, what the app violates, what it could not tell and where the report card is. A whole run takes a while and uses many tokens.

## Where things go

| What | Where |
|---|---|
| Report cards | `~/.terminu-harness/judge/<app>/<time>-tdp-<version>/`: `report.md`, every verdict, the evidence |
| tdp files fetched | `~/.terminu-harness/cache/tdp/` |
| The built terminal | `~/.terminu-harness/bin/` |

| Variable | What it does |
|---|---|
| `TERMINU_HARNESS__HOME` | Keep everything somewhere other than `~/.terminu-harness` |
| `TERMINU_HARNESS__OUT` | Put report cards somewhere else |
| `TERMINU_HARNESS__TDP_REPO` | Read tdp from a local clone of terminu rather than GitHub |

## Limits

- No Windows.
- The terminal sends no mouse events: items that need a mouse are judged "can't tell".
- Times measured on the screen (animations, how long a toast stays) are judged within 33ms until tdp says how close is close enough.

## Related

- [terminu](https://github.com/vulcanshen/terminu): terminu design (tdp) and its checklist
- [terminu-go](https://github.com/vulcanshen/terminu-go): the Go lib for terminu family apps
- [CHANGELOG](CHANGELOG.md)
- [docs/dev-remarks.md](docs/dev-remarks.md): how it works, for developers

## License

[GPL-3.0](LICENSE)

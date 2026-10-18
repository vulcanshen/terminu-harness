# terminu-harness

[![License](https://img.shields.io/badge/license-GPL--3.0-blue)](LICENSE)

**Language**: English · [繁體中文](README-zh_TW.md)

**A Claude Code plugin for TUI apps made with [terminu design](https://github.com/vulcanshen/terminu) (tdp).** It judges how many terminu genes an app has: every item on the tdp checklist the app meets is one gene.

## What it does

`/terminu:judge` opens the app in a terminal of its own, reads its help, presses keys and looks at the screen, the way each checklist item says to check it. Items that need the source or other files are judged when the app's repo is at hand. It writes a report card:

- the verdict on every item (meets, departs, violates, not applicable, or can't tell) with what it saw
- the gene count for that version of tdp, as `tdp-0.7.4 860 gene`
- what the app does that tdp does not cover yet, to report back to tdp

It judges any TUI app, not only terminu family members, whatever it is written in, with or without its source. It needs no tmux or graphical desktop, and it leaves the app's repo as it is.

## Requirements

- [Claude Code](https://claude.com/claude-code)

## Install

In Claude Code:

```
/plugin marketplace add vulcanshen/terminu-harness
/plugin install terminu@terminu-harness
```

## Usage

```
/terminu:judge
```

## Related

- [terminu](https://github.com/vulcanshen/terminu): terminu design (tdp) and its checklist
- [terminu-go](https://github.com/vulcanshen/terminu-go): the Go lib for terminu family apps
- [CHANGELOG](CHANGELOG.md)

## License

[GPL-3.0](LICENSE)

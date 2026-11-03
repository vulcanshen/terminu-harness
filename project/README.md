# terminu project

**Language**: English · [繁體中文](README-zh_TW.md)

terminu project is what a family member's repo looks like — the directories, the Makefile, the version, the CHANGELOG,
the release — and its template files. [tdp](https://github.com/vulcanshen/terminu) says what an app should be, whatever
the language or the build tool; this page is the family's upkeep, done with Go and goreleaser, so it lives in
terminu-harness, the tool that makes and releases the family's apps. tdp's start of chapter E and
[E6](https://github.com/vulcanshen/terminu/blob/main/principle/rules.md) hand a family member's repo to this page.

- **Moved from tdp**: from tdp v1.2.0-beta.1 on (the user, 2026-11-03). It used to be tdp's `project/`, last in tdp
  v1.1.10; the section IDs project.1–8 and the checklist's `project-` IDs stay as they were.
- **Which one applies**: while the tdp an app follows (the one its go.mod's terminu-go follows) still has `project/`
  (v1.1.10 and before), that tdp's, checked by that terminu-go's tdptest; once it has none, this page, checked by
  harness's `check` against the templates of the harness installed.
- **Citing**: an app's files cite tdp as `tdp` plus the ID, and this page as `harness` plus the ID: a section
  `harness project.3`, the whole page `harness project`, a template `harness project/install.sh`, a checklist item
  `harness project-13` ([`items.md`](items.md)). The parts of this page that get copied into an app are written that
  way too.

## project.1 What is a requirement

This page binds family members only (tdp's start of chapter E). For a family member, everything written here
is a requirement. Only three things are not: examples (however they are introduced: "e.g.",
"for example", "Example:", "such as", an "Examples" column), the "Why", and paragraphs marked "Implementation
reference (not a requirement)". Where it says "the app may decide for itself", that part is left to the app
explicitly, with what still holds inside that scope written beside it
([tdp P5](https://github.com/vulcanshen/terminu/blob/main/principle/README.md#p5-leaving-it-to-the-app-explicitly)).

The template files are requirements too: they are copied as they are, changed only where the last column of the
"Template files" table says.

**An app made with tdp** (not a family member): its repo's directories, Makefile, where the version lives, CHANGELOG
and release steps, the app may decide for itself. This page and its templates are a starting point it may copy and
change as it likes. What still holds is tdp's chapter E (e.g. E1 `version`, E6 `install.sh` and the tap, E6's
screen tests).

**Why**: what this page is worth lands on whoever keeps the whole family: a release problem one app runs into is fixed
once for all of them, the family shares one Homebrew tap, and a diff across the repos shows which one drifted. An app
made on its own gets none of that, so binding it gains nothing (the user, 2026-10-27).

## project.2 Template files

| File | Copied to | What the app changes |
|---|---|---|
| [`family.mk`](family.mk) | `family.mk` | Nothing |
| [`release.yml`](release.yml) | `.github/workflows/release.yml` | With a Windows build, `windows-latest` added to the test matrix |
| [`goreleaser.yaml`](goreleaser.yaml) | `.goreleaser.yaml` | `<app>`, the lines marked `app:`; with a Windows build, the comments marked `Windows` opened |
| [`install.sh`](install.sh) | `install.sh` | `<app>`, `WINDOWS`, the blocks marked `app's own` |
| [`uninstall.sh`](uninstall.sh) | `uninstall.sh` | `<app>`, the blocks marked `app's own` |
| [`version.go`](version.go) | `internal/version/version.go` | `<app>` |
| [`gitignore`](gitignore) | `.gitignore` | `<app>`; the app's own lines go after the last one |

- **An app made with tdp** (not a family member; project.1) starting from the templates replaces the `vulcanshen` in
  them (the repo path, the tap's owner, the homepage) with its own account.
- **How far a mark reaches**: a `# app:` comment in `goreleaser.yaml` stays as it is and covers the value of the key
  right below it, a list or a whole block included (`ignore`, `description`, `caveats`); when the key it names is not
  in the template yet (`dependencies`), the app adds that key where the comment is; other keys (`test:`, say) stay as
  they are. Opening `Windows` takes the line marked `# Windows` with the commented lines of its block below it,
  removing only the leading `# `; the trailing `# Windows` stays. The `<APP>` in a comment of `uninstall.sh` is not
  `<app>` and stays as it is.
- **Checking for drift**: put `<app>` back in place of the app's name and diff against the templates here; only what
  the last column of the table names should differ. Harness's `check` compares them this way.
- **A template changes here first**: a harness release (its CHANGELOG names the template changed), then each member
  copies it the next time it runs `check` (upgrading, releasing), never an edit in an app. A problem an app runs into
  with a template is reported to terminu-harness.

**Why**: the family used to share one Makefile template, and copying it about let it drift — the directories `fmt`
covers, whether `check` runs `-race`, the column width of `make help`, the brackets in the section names were each a
little different from app to app. With one file, a diff shows what differs, and a change has one place to go.

## project.3 Directories

```
<app>/
├── cmd/<app>/main.go               the main package
├── internal/                       the app's own packages
│   └── version/version.go          the version (see "Version")
├── docs/                           icon.svg, the README's demo gif (tdp E6, E7)
├── .local/                         not in git; demos/ holds the VHS tapes
├── .github/workflows/release.yml
├── .goreleaser.yaml
├── .gitignore
├── go.mod, go.sum
├── family.mk
├── Makefile
├── install.sh, uninstall.sh
├── CHANGELOG.md
├── README.md, README-zh_TW.md
└── LICENSE
```

- The app's own Go code lives only in `cmd/<app>/` and `internal/`. The `go install` path is
  `github.com/vulcanshen/<app>/cmd/<app>@latest` (an app made with tdp uses its own account:
  `github.com/<account>/<app>`; tdp's start of chapter E).
- **From v2 on, the module path ends in the major version**, as Go requires: kbu v4 is `github.com/vulcanshen/kbu/v4`,
  its imports and its `go install` path (`github.com/vulcanshen/kbu/v4/cmd/kbu@latest`) with it. The release that
  moves to a new major (v2 and up) changes them in the same commit. **Why**: without it Go does not take the tags
  from v2 on as the module's versions: `go install …@latest` installed kbu v1.7.12 while kbu was v4.0.0, and the
  binaries' build information read a pseudo-version (the user, 2026-10-28).
- **Small development tools** (Go ones too) may go in a top-level `tools/`; they are not the app itself.
- Other top-level directories and files not in the tree (e.g. end-to-end tests), the app may decide for itself.
- **Claude Code's files stay out of git**: `CLAUDE.md` and what is under `.claude/` are not tracked (the template
  [`gitignore`](gitignore) leaves them out), except harness's `.claude/rules/terminu.md`. They may sit in the
  working tree. **Why**: they are the maintainer's own instructions and settings, full of the paths of one machine;
  a public repo has no use for them, and they gave the maintainer's paths away (the user, 2026-10-28).
- **`docs/icon.svg` is the family mark's original.** Go's `go:embed` cannot reach `docs/`, so the copy the splash uses
  lives in the code (a copy of the file under `internal/`, or the cells written out), and the tdp E6 test keeps it
  identical to `docs/icon.svg`.
- Whether to build for Windows, the app may decide for itself; an app with a Windows build also has `install.ps1` and
  `uninstall.ps1`, which are not templates, and how to write them the app may decide for itself.

**Why**: most family apps were already `cmd/<app>` plus `internal/`, so following them moves the fewest apps. With code
in two places and only repo-level files at the root, anyone opening a family repo knows where things are.

## project.4 Version

- The version lives in `internal/version` (template [`version.go`](version.go)): `Version` is the version without its v,
  and `Display()` adds the v when it is printed.
- goreleaser injects the tag without its v when it releases (e.g. tag `v0.2.0` → `0.2.0`). `<app> version` and the
  splash print `Display()` (e.g. `v0.2.0`; tdp E1, S5).
- A local `make build` or `go build` injects nothing; the version is `dev`.

**Why**: `dev` shows at a glance that a build is not a release. Where versions are compared (e.g. ccdrawu's update
check), `Version` without the v is used, and `dev` tells it the build is local, so it does not compare.

## project.5 Makefile

An app's `Makefile` only sets variables and includes [`family.mk`](family.mk); its own targets go in a last section:

```make
# <app> — build / run / test / package（家族共用的 target 在 family.mk）
#
# 跑 `make`（或 `make help`）列出所有指令。

BINARY := <app>
include family.mk

##@ <app> 自己的
```

**Variables** (set before the `include`):

| Variable | Meaning | Default |
|---|---|---|
| `BINARY` | The executable's name, which is also the name of the directory under `cmd/` | Required |
| `CHECK` | The app's own targets that `check` runs last (e.g. ccdrawu's `mod-check`) | Empty |
| `GOFMT_DIRS` | The directories gofmt covers | `cmd internal` |
| `TAPES` | The VHS tapes `make gif` runs | `.local/demos/demo.tape` |
| `VHS` | The path to vhs | `vhs` |

**family.mk's targets**, which `make help` lists in this order, the app's own section last:

| Section | Targets |
|---|---|
| 編譯（build） | `build`, `install`, `uninstall`, `tidy` |
| 發布（release） | `release-check`, `snapshot` |
| 執行（run） | `run`, `dev` |
| 測試 / 檢查（test） | `test`, `test-race`, `check-lib`, `vet`, `fmt`, `fmt-check`, `check` |
| 打包（package） | `package` |
| 展示（demo） | `gif` |
| 其他 | `clean`, `help` |

- **`check` is run before a commit**: `fmt-check`, `vet`, `test-race`; with a `go.work`, `check-lib` as well (`go.work`
  off, built and tested with the versions `go.mod` names); `CHECK` last.
- **`tidy` and `snapshot` turn `go.work` off**: CI has no `go.work` (git ignores it), and the local run matches it.
- **The app's own targets all go in the last section**, `##@ <app> 自己的`, with names that are not family.mk's.
- family.mk's descriptions are in Traditional Chinese, as the family's Makefiles were.

**Why**: CI runs `go test -race`, and `check` runs the same, so a local pass means a CI pass. `check-lib` catches "it
passes locally with `go.work` pointing at a lib that is not released yet, and fails with `go.mod`" — it must pass before a
release. With the app's own targets together at the end, `make help` shows which are the family's.

## project.6 CHANGELOG

```markdown
# Changelog

Each release names the version of terminu design (tdp) it follows.

## Unreleased

- …

## v0.1.2 — tdp v0.6.11

- …
```

- **Each release's heading** is the tag pushed (with its v), then ` — tdp <version>`: the version of tdp this release
  follows.
- **A beta** has a heading of its own, its tag with the pre-release part (`## v0.2.0-beta.1 — tdp v1.1.0-beta.2`); the
  stable release after it covers everything since the previous stable release, the betas' changes included, so that
  someone who skipped the betas can read it.
- **`## Unreleased`**: a change goes into this section when it lands (opened at the top when there is none); a release
  turns its heading into the release's heading.
- **Releases made before tdp v0.7.0 keep the headings they were written with.**
- Written in English. A release's notes are that release's section of the CHANGELOG (see "Release").

**Why**: with the tag at the start of the heading, `release.yml` finds the section with `$2 == tag`, and one awk serves
the whole family; with the tdp version written, it is plain which version of the spec each app has reached.
`Unreleased` records a change as it lands, so nothing is left out when a session is cut off or a release carries several
changes. Old headings stay: `release.yml` only looks for the new tag's section, so they do no harm, and which version of
tdp an old release followed often cannot be found; filling it in would get it wrong.

## project.7 Release

1. `make check` passes.
2. The CHANGELOG's `## Unreleased` becomes `## vX.Y.Z — tdp vA.B.C`.
3. A commit, and the tag `vX.Y.Z`.
4. The commit and the tag are pushed. [`release.yml`](release.yml) then:
   - tests: `go test -race` on ubuntu and macOS (an app with a Windows build adds Windows, without `-race`);
   - runs goreleaser: builds, makes the GitHub release, pushes to the Homebrew tap `vulcanshen/homebrew-tap` (an app
     made with tdp pushes to its own tap; and the scoop bucket, with a Windows build);
   - sets the release's body to that section of the CHANGELOG. This step runs even when goreleaser's tap push failed,
     and skips when the release has not been made.

- **A beta**: before a stable release, the app may release a beta: tag `vX.Y.Z-beta.N` (N from 1), the same steps.
  goreleaser makes it a GitHub pre-release and leaves the Homebrew tap and the scoop bucket alone (`skip_upload:
  auto`), so brew, scoop and `install.sh` (the latest release) keep the stable one; `go install …@vX.Y.Z-beta.N`
  installs the beta. **Why**: v1.0.0 needed v1.0.1 the same day (the user, 2026-10-27: release a beta before a
  stable one); a beta goes through the whole release without touching what brew and scoop users get.
- **The release notes are written outside the checkout** (`$RUNNER_TEMP`): a file in the checkout leaves the tree
  dirty, and Go writes the binaries' build information as `vX.Y.Z+dirty` (`go version -m`), as the family's v1.0.0
  releases were.
- **The token**: goreleaser reads `GITHUB_TOKEN` for the release and for the tap. The workflow hands it the repo secret
  `ACTION_TOKEN` (a PAT) when one is set, and the built-in token otherwise — the built-in one cannot write to the tap. A
  release run by hand takes the shell's `GITHUB_TOKEN`.
- **`make snapshot`** packages locally, without releasing or pushing to the tap.
- **`make release-check`** validates `.goreleaser.yaml`. goreleaser marks `brews` deprecated, so `goreleaser check`
  fails; `release-check` still prints the warning, but that failure alone does not count.

**Why**: Homebrew gets a formula, not a cask: files a cask installs are stamped `com.apple.quarantine`, and macOS
Gatekeeper stops an unsigned executable; a formula's are not. One token throughout, so a release by hand and the
workflow use the same name. The release body is set a second time because goreleaser's `--release-notes` is
unreliable; it is set even when the tap push failed because kbu had a tap push fail with 401 and the release's body was
left empty.

## project.8 Reference implementation

demou: the family's app for learning the operations, built from scratch to this page. It meets problems with the
templates first, and reports them to terminu-harness.

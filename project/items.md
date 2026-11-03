# terminu project: items

**Language**: English · [繁體中文](items-zh_TW.md)

The items of [terminu project](README.md), moved from the checklist of tdp v1.1.10 with their IDs (cited as `harness project-13`). What each column means and
what an item applies to: tdp's [checklist](https://github.com/vulcanshen/terminu/blob/main/checklist/README.md); they check family members only (project.1).

| ID | Source | What to check | Applies to | How to check |
|---|---|---|---|---|
| project-1 | project.2 | `family.mk` is copied from the template as it is to `family.mk`, unchanged | everyone | Read the file: a diff against harness's `project/family.mk` shows no difference |
| project-2 | project.2 | `release.yml` is copied from the template to `.github/workflows/release.yml`; the only change is `windows-latest` added to the test matrix, with a Windows build | everyone | Read the file: a diff against harness's template differs only by `windows-latest` in the matrix (with a Windows build) |
| project-3 | project.2 | `goreleaser.yaml` is copied from the template to `.goreleaser.yaml`, changing only `<app>` and the lines marked `app:`, and, with a Windows build, opening the comments marked `Windows` | everyone | Read the file: put `<app>` back in place of the app's name and diff against harness's template; only these places differ |
| project-4 | project.2 | `install.sh` is copied from the template, changing only `<app>`, `WINDOWS` and the blocks marked `app's own` | everyone | Read the file: put `<app>` back in place of the app's name and diff against harness's template; only these places differ |
| project-5 | project.2 | `uninstall.sh` is copied from the template, changing only `<app>` and the blocks marked `app's own` | everyone | Read the file: put `<app>` back in place of the app's name and diff against harness's template; only these places differ |
| project-6 | project.2 | `version.go` is copied from the template to `internal/version/version.go`, changing only `<app>` | everyone | Read the file: put `<app>` back in place of the app's name and diff against harness's template; nothing else differs |
| project-7 | project.2 | `gitignore` is copied from the template to `.gitignore`, changing only `<app>`, with the app's own lines after the last one | everyone | Read the file: put `<app>` back in place of the app's name and diff against harness's template; only the lines added at the end differ |
| project-8 | project.2 | With a Windows build, `windows-latest` is added to the test matrix of `release.yml` | Windows build | Read the file: the matrix in `.github/workflows/release.yml` |
| project-9 | project.2 | With a Windows build, the comments marked `Windows` in `.goreleaser.yaml` are all opened | Windows build | Read the file: `.goreleaser.yaml` |
| project-10 | project.2 | With a Windows build, `WINDOWS` in `install.sh` is set to `yes` (as the template's comment says) | Windows build | Read the file: `WINDOWS` in `install.sh` |
| project-11 | project.3 | The app's own Go code lives only in `cmd/<app>/` and `internal/`; development tools may go in a top-level `tools/` | everyone | Read the files: every `.go` file outside `tools/` is under these two directories |
| project-12 | project.3 (diagram) | The main package is in `cmd/<app>/main.go` | everyone | Read the files |
| project-13 | project.3 | `go install github.com/<account>/<app>/cmd/<app>@latest` installs its latest release, with `/vN` after `<app>` from v2 on (a family member's account is `vulcanshen`, an app made by tdp uses its own) | everyone | Read the files: the module in `go.mod` is `github.com/<account>/<app>`, ending in `/vN` when the latest release is vN.x.y with N ≥ 2, and `cmd/<app>` is the main package; `go list -m <module>@latest` gives the latest release |
| project-14 | project.3 (diagram) | The README's demo gif is in `docs/` | everyone | Read the files: the gif the README refers to is under `docs/` |
| project-15 | project.3 (diagram) | `go.mod` and `go.sum` are at the root | everyone | Read the files |
| project-16 | project.3 (diagram) | `LICENSE` is at the root | everyone | Read the files |
| project-18 | project.3 | The family mark's original is `docs/icon.svg`; the copy the splash uses lives in the code (a copy of the file under `internal/`, or the cells written out) | everyone | Read the files: `docs/icon.svg` exists; read the code: the splash draws the copy under `internal/` or the cells written out |
| project-19 | project.3 | An app with a Windows build also has `install.ps1` and `uninstall.ps1` | Windows build | Read the files |
| project-37 | project.3 | No `CLAUDE.md` and no file under `.claude/` is tracked, except `.claude/rules/terminu.md` | everyone | Run `git ls-files \| grep -E '(^\|/)CLAUDE\.md$\|^\.claude/'`: nothing but `.claude/rules/terminu.md` |
| project-20 | project.4 | `<app> version` and the splash's version print `Display()` of `internal/version` | everyone | Read the code: both places call `Display()` |
| project-21 | project.4 | A release's version is the tag injected without its v, and printed with the v (e.g. tag `v0.2.0` prints `v0.2.0`) | everyone | Read the files: `.goreleaser.yaml` injects the tag without its v and `internal/version` prints it with the v; for an app that has released, install a release and run `<app> version`: it prints the same as the tag |
| project-22 | project.4 | A local `make build` or `go build` injects no version; the version is `dev` | everyone | After a local `make build`, run `./<app> version`; it prints `dev` |
| project-23 | project.5 | `Makefile` only sets variables and does `include family.mk`; the app's own targets all go in the last section, `##@ <app> 自己的` | everyone | Read the file; run `make help`: the family's sections in family.mk's order, the app's own section last |
| project-24 | project.5 | The variables are set before `include family.mk` | everyone | Read the file |
| project-25 | project.5 | `BINARY` is set to the app's name, which is also the name of the directory under `cmd/` | everyone | Read the file |
| project-26 | project.5 | The two comment lines at the start of the Makefile are the template's (`<app> — build / run / test / package（家族共用的 target 在 family.mk）` and 「跑 make（或 make help）列出所有指令。」) | everyone | Read the start of the file |
| project-27 | project.5 | The app's own targets have names that are not family.mk's | everyone | Read the files: compare the target names on both sides |
| project-28 | project.6 | `CHANGELOG.md` starts as in the template: `# Changelog`, then `Each release names the version of terminu design (tdp) it follows.` | everyone | Read the start of the file |
| project-29 | project.6 | Each release's heading is the tag pushed (with its v), then ` — tdp <version>` | has released | Read the file: the heading of every release made after tdp v0.7.0 |
| project-30 | project.6 | Changes not yet released are written in `## Unreleased` at the top | everyone | Read the file and the commits after the latest tag: every change is recorded in `## Unreleased` at the top |
| project-31 | project.6 | Releases made before tdp v0.7.0 keep the headings they were written with | has released | Read the file: in `git log -p CHANGELOG.md`, the old releases' headings have not been rewritten |
| project-32 | project.6 | The CHANGELOG is written in English | everyone | Read the file |
| project-33 | project.6 | A GitHub release's notes are that release's section of the CHANGELOG | has released | Compare `gh release view <tag>` with that section of the CHANGELOG |
| project-34 | project.7 | A release's tag is `vX.Y.Z`, a beta's `vX.Y.Z-beta.N` | has released | Run `git tag` |
| project-35 | project.7 | `make check` passes on the repo as it is now | everyone | Run `make check` on the repo as it is now |
| project-36 | project.7 | A release with a Windows build also pushes to the scoop bucket | Windows build | Check that the scoop bucket has this release |
| project-38 | project.7 | A beta, once released, is a GitHub pre-release, and the Homebrew tap (and the scoop bucket) still has the latest stable release | has released | `gh release view <beta tag> --json isPrerelease`; the version in the tap's formula |
| project-39 | project.7 | The latest release's binaries have no `+dirty` in their build information | has released | Download a binary of the latest release and run `go version -m` on it: the main module's version is the tag, `vcs.modified=false` |

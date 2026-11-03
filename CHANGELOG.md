# Changelog

Each release names the version of terminu design (tdp) it follows.

## Unreleased

- The terminal's `keys` types `text` `gap_ms` after the last key, as it spaces the keys, where it used to follow at once: an app could read a key and the text as one input (`/` and `line` as `/line`, one key of several characters to bubbletea), and the key never acted.

## v1.2.0-beta.3 — tdp v1.2.0-beta.1

- A fourth backing in `work.yaml`, `check: repo`: harness's `check` guards the task, reading a family member's repo against the family's templates. It counts only where harness reads the repo: a family member's, its tdp without `project/`.
- A subtest named as a backing (`TestX/sub`) has to run: `check` runs it with `go test -run` and looks for it in the output, where it used to look up `TestX` alone, so a subtest that is gone (terminu-go's `TestTDP/repo`, say) no longer passes. A gone `…/repo` says to write `check: repo`; `/terminu:newapp`'s first upgrade to a tdp without `project/` swaps it, and `/terminu:analysis` writes `check: repo` for the repo's task.
- `check` no longer passes what it could not check: tdp it could not read, tests that would not run head the report as not checked, and the command exits 1 (the MCP tool reports an error). Reading tdp from GitHub tries again after a network error, a 5xx or a 429.

## v1.2.0-beta.2 — tdp v1.2.0-beta.1

- The newest terminu-go no longer counts a beta, as Go's `@latest` does not: `/terminu`, `/terminu:newapp`, `/terminu:release` and `/terminu:analysis` take the newest release (`go list -m …@latest`), a beta only when the app is on a beta already (the last of `go list -m -versions`, so a release replaces its betas) or when you name one. A family's beta round is for the members it names; the rest are no longer upgraded to it when they release.
- `tdp_file` and `tdp_section`, given no version, read the newest tdp release rather than a beta.
- `/terminu:release` reads the release rules from tdp's `project/` while the app's tdp still has it (v1.1.10 and before), and from harness's after that.

## v1.2.0-beta.1 — tdp v1.2.0-beta.1

- The family's repo templates move here from tdp: tdp says what an app should be, whatever the language, and templates written for Go and goreleaser are the family's upkeep, kept with the tool that makes and releases the family's apps. `project/` holds terminu project as tdp v1.1.10 had it: its page (sections project.1–8), its checklist items (the 38 `project-` items, their IDs kept) and the seven templates, whose first comment now cites `harness project/…`.
- Citing them: `harness project.3`, `harness project` (the whole page), `harness project/install.sh`, `harness project-13`, in the same `tdp` list of `spec.yaml`. `tdp_section` reads them; `check` checks them against `project/`, and a `tdp project…` citation the app's tdp no longer has says to write it with `harness`. An app whose tdp still has `project/` (v1.1.10 and before) may cite it as before.
- `check` compares a family member's repo with the templates, the files, the Makefile and the CHANGELOG, once the tdp it follows has no `project/`: what terminu-go's tdptest checked until now, moved here with the templates, with messages in Traditional Chinese.
- `/terminu:newapp`, `/terminu:analysis` and `/terminu:release` read the templates from harness; the first upgrade to a tdp without `project/` rewrites the citations and the templates' first comments. `/terminu:release` runs `check` after `make check`.
- `/terminu:reportgap` reports a problem with the family's templates to terminu-harness.

## v1.1.0 — tdp v1.1.0

Everything since v1.0.1, the betas v1.1.0-beta.1 to beta.3 included:

- tdp's pre-releases: `tdp_file`, `tdp_section` and `check` take a version such as `v1.1.0-beta.1`, so an app on a terminu-go beta reads its tdp and its `spec.yaml` may follow it; the newest tag counts pre-releases, each before its release.
- `/terminu:release` releases a beta: the version may be `vX.Y.Z-beta.N` and the CHANGELOG heading is its tag; it goes out as a GitHub pre-release that the Homebrew tap and the scoop bucket do not get, and the stable release after it covers everything since the previous stable one, the betas' changes too.
- `/terminu:release`, on a new major version from v2 on, ends the module path in `/vN` and changes every import and the `go install` paths in the README and the docs, in the same commit as the CHANGELOG heading, so `go install …@latest` gets the release.
- An app made with tdp (not one of vulcanshen's): `/terminu:newapp` starts it from tdp's project templates with your account in them and in `go mod init`, and you may change them afterwards; an upgrade does not carry tdp's project changes into it, and `/terminu:release` follows the app's own CHANGELOG and release steps where it changed them.

## v1.1.0-beta.3 — tdp v1.1.0-beta.3

- `/terminu:release`, on a new major version from v2 on, ends the module path in `/vN` and changes every import and the `go install` paths in the README and the docs, in the same commit as the CHANGELOG heading, so `go install …@latest` gets the release.

## v1.1.0-beta.2 — tdp v1.1.0-beta.2

- `/terminu:release` releases a beta: the version may be `vX.Y.Z-beta.N` and the CHANGELOG heading is its tag; it goes out as a GitHub pre-release that the Homebrew tap and the scoop bucket do not get, and the stable release after it covers everything since the previous stable one, the betas' changes too.

## v1.1.0-beta.1 — tdp v1.1.0-beta.1

- tdp's pre-releases: `tdp_file`, `tdp_section` and `check` take a version such as `v1.1.0-beta.1`, so an app on a terminu-go beta reads its tdp and its `spec.yaml` may follow it; the newest tag counts pre-releases, each before its release.
- An app made with tdp (not one of vulcanshen's): `/terminu:newapp` starts it from tdp's project templates with your account in them and in `go mod init`, and you may change them afterwards (tdp v1.0.2); an upgrade does not carry tdp's project changes into it, and `/terminu:release` follows the app's own CHANGELOG and release steps where it changed them.

## v1.0.1 — tdp v1.0.1

- `/terminu`: a newer terminu-go is offered as an upgrade only when it changes the lib or what the app has to do; one that changes only documents (such as tdp v1.0.1, a family member's README) is told of and waits for the release.
- `/terminu:release`: brings the app up to the newest terminu-go before the checks, so the upgrades left waiting go in with the release.

## v1.0.0 — tdp v1.0.0

- `/terminu`: in the folder of an app the harness keeps, checks what the app is waiting for (a newer terminu-go and the tdp it follows, answers to reports, TODOs) and carries on; in a terminu-go app it has not kept yet, offers to take it over; otherwise asks what you want to do (make a new app, review, release, or report to tdp or terminu-go) and takes you there.
- `/terminu:newapp`: asks what the app is for, one question at a time, and writes it down in your words (`idea.md`); matches each request to tdp (`spec.yaml`) and builds the app with terminu-go from the tdp project templates task by task, running the tests, opening each piece in a terminal of its own and writing down how each task was built and what tests back it (`work.yaml`). Where tdp has no answer it drafts an issue for tdp or terminu-go and works around it meanwhile. It picks up where it left off when you come back, takes new requests, and follows a newer tdp when you agree.
- `/terminu:review`: a separate, read-only reviewer judges every task in the app's records against tdp, without reading the code; you reconfirm the checks only you can do; what does not follow tdp is listed in `review.yaml` and fixed without asking, in the records or in the code; you are asked only whether to report something to tdp.
- `/terminu:release`: checks, reviews, tells you what is still open or not built, turns `## Unreleased` into the version's heading the tdp way, and tags and pushes once you agree.
- `/terminu:reportgap`: drafts an issue for tdp or terminu-go, files it once you agree, and notes it on the task it holds up.
- `/terminu:analysis`: takes over a terminu-go app made before the harness: writes what it does from the README (`idea.md`, `spec.yaml`) and how each part was built and which tests back it from git history (`work.yaml`), without reading the code; what neither says is marked as not known, for review to list. An app on an older terminu-go is first brought up to the newest, with what tdp changed in between that its tests do not check (the app's own words, README, project files). It adds the app's rules (`.claude/rules/terminu.md`), deletes the old developer docs, and tells you what was made.
- A terminal of its own (a PTY and an emulator, no tmux or desktop needed) to open the app being made, press keys and read the screen, its colours or a picture of it; and tdp's files at any version.
- A commit hook: in an app the harness keeps, a commit that changes the app's code without `harness-assets/work.yaml` is stopped, unless its message says `Terminu-Work: unchanged` (the code changed, what the app does did not).
- `tdp_section` (an MCP tool): parts of tdp by the IDs an app cites them with (`K4`, `P5`, `terms.3`, `panel.5`, a whole file such as `confirm`, a template such as `project/install.sh`, or a checklist item such as `K9-6`, which comes with the sections it points at), each with the checklist items that come from it.
- `check` (an MCP tool, and `terminu-harness check` on the command line): checks the app's `harness-assets/` — `idea.md`, `spec.yaml`, `work.yaml`, `review.yaml` — against each other; every tdp citation in `spec.yaml` against the tdp the app's terminu-go follows (read from `go.mod` and terminu-go's CHANGELOG); every test in `work.yaml` against `go test -list` and every commit against git. It lists what is left to do: tasks not built yet, tasks built to an older spec, a review behind the spec.

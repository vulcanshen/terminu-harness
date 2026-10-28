# Changelog

Each release names the version of terminu design (tdp) it follows.

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

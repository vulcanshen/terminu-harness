# Changelog

## [Unreleased]
### Added
- `/terminu`: in the folder of an app it has been making, checks what the app is waiting for and carries on; otherwise asks what you want to do (make a new app, review changes, release, or report to tdp or terminu-go) and takes you there.
- `/terminu:newapp`: asks what the app is for, one question at a time, writes a spec and builds the app with terminu-go from the tdp project templates piece by piece, running the tests and opening each piece in a terminal of its own; where it gets stuck it drafts an issue for tdp or terminu-go and works around it meanwhile, and it shows you the app when it is done. It picks up where it left off when you come back.
- `/terminu:review`: a separate, read-only reviewer looks at a range of changes against tdp and the app's spec, every problem with where it comes from; problems are fixed and looked at again, at most twice, then you decide.
- `/terminu:release`: checks, reviews what changed since the last tag, turns `## Unreleased` into the version's heading the tdp way, and tags and pushes once you agree.
- `/terminu:reportgap`: drafts an issue for tdp or terminu-go and files it once you agree.
- A terminal of its own (a PTY and an emulator, no tmux or desktop needed) to open the app being made, press keys and read the screen, its colours or a picture of it; and tdp's files at any version.

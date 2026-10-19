# Changelog

## [Unreleased]
### Added
- `/terminu`: asks what you want to do (make a new app, carry on with one, or judge one) and takes you there.
- `/terminu:new-app`: asks what the app is for, one question at a time, writes a spec for you to confirm, builds the app with terminu-go from the tdp project templates piece by piece, drafts an issue for tdp or terminu-go where it gets stuck and works around it meanwhile, and judges the app once it is done. It picks up where it left off when you come back.
- `/terminu:report-gap`: drafts an issue for tdp or terminu-go and files it once you agree.
- `/terminu:judge`: judges a TUI app against the checklist of a version of terminu design (tdp), every item, and writes a report card: passed or failed (any item violated fails it), each verdict with its evidence, the gene count (`tdp-0.8.0 860 gene`) and what tdp does not cover yet. How to judge is the harness's own: the verdicts, the pass line, the tolerance for times measured on the screen. It runs the app in a terminal of its own, no tmux or desktop needed, and reads the app's repo when you give one.

# Changelog

Each release names the version of terminu design (tdp) it follows.

## v1.3.9 — tdp v1.3.9

Follows tdp v1.3.9 and terminu-go v1.3.11, where a README's first section shows the app's icon.

- `harness project.3`: the tree's `docs/` says `icon.svg` is also the icon atop the README.
- Checked against tdp v1.3.9, which needs no other change: `/terminu:newapp` writes both READMEs by tdp E7, so a new app's README shows the icon, and a member adds it at its next release as `/terminu:newapp`'s upgrade already says for what no task covers. `check` does not read a README, and still does not; the reviewer and the templates stay as they are.

## v1.3.8 — tdp v1.3.8

Follows tdp v1.3.8 and terminu-go v1.3.10, where data naming something outside the app is kept when its thing is not there, and the one panel with a number on its screen may take the digits.

- `/terminu:newapp`: an entry of the user's data (a list the user writes for the app: pinned containers or kinds, favourites, names to watch) that names something outside the app is not left out when that thing is not there (tdp E3-43). It stays as written: no `Notice`, no toast, no WARN, and it is still there when the app writes the file again, since what goes back is the list the user wrote, not the entries that work now; it works again once its thing is back, and only the user removes it. Warning once and leaving the entry out (tdp E3-41) stays for a value in the state.
- Checked against tdp v1.3.8, which needs no other change: the digits of the one panel with a number are `Panel.Digits`, and `Run` already refuses a digit as a panel's hotkey, so there is no way to take them by hand; its hint, its row in `?`, a digit past N and a second panel taking them back are the lib's. The reviewer reads K1 from tdp, so such a panel is no longer judged against K1-7. `check`, the reviewer and the templates stay as they are.

## v1.3.7 — tdp v1.3.7

Follows tdp v1.3.7 and terminu-go v1.3.9, where a mode is named by the -ing word of what is done in it, an app's own mode too, and its layers are joined with `:`.

- `/terminu:newapp`: a mode of the app's own takes such a name in `Mode.Name`, one word (a game's `Playing`), not `Play`, `Drag` or `Visual`, and so does a layer it goes into with `Panel.EnterLayer` (`selecting`); the lib joins them with `:` (tdp K11-29). The lib does not check the name, so a wrong one still builds and passes its tests. A Drag mode is now a Dragging mode, named `Dragging`, and the list's Diff mode is the Diffing mode; the Diff popup, `[D]iff` and the `D` key keep their name, and so do the APIs (`DiffMode`, `App.DiffModeAction`).
- `harness-assets.md`: the implementation of a mode of the app's own says what it is called, so a review can judge its name.
- Checked against tdp v1.3.7, which needs no other change: a terminal's scrolling back with no History mode is the lib's, and nothing in harness names the History mode or joins layers with `>`. `check`, the reviewer and the templates stay as they are.

## v1.3.6 — tdp v1.3.6

Follows tdp v1.3.6 and terminu-go v1.3.8, where the user moves the lines between panels from the focused one (tdp screen.9) and the lib keeps the new places for the app to write to its state.

- `/terminu:newapp`: on a screen with more than one panel, the app keeps where the user left the lines in its own state (tdp E3): it reads them into `App.Lines` before `Run` and writes what `App.OnLines` gets back to its state, with an Error toast when that fails. Without it the lines last only until the app quits. When to write is up to the lib: it calls `OnLines` once the user moves a line or `=` changes something, and also after it puts back kept lines that no longer fit. Those are not a value that does not fit: no `KeepAside`, no `Notice` (tdp E3's one exception). The five keys, the menu rows, the toast when nothing moves and finding the lines in `Layout` are the lib's, and no panel, screen or global hotkey takes `>` `<` `+` `-` `=` (`Run` refuses them). `Panel.MinWidth` and `MinHeight` set a panel's least size. The terminu-go README's section on moving the lines joins what is read first.
- Checked against tdp v1.3.6, which needs no other change: the toast's place (toast.2), the disabled rows (M6), leaving resizing out of the App Log (note.7) and the pointer in L1 are the lib's or need nothing from the app. `check`, the reviewer and the templates name none of the five keys and stay as they are.

## v1.3.5 — tdp v1.3.5

Follows terminu-go v1.3.7, whose `SelfUpdate` no longer asks GitHub's API for the latest release: to a network that has not signed in, the API answers 60 times an hour, every program on it counted together, and then says 403.

- Templates: `install.sh` finds the latest release the same way, from where `github.com/<owner>/<app>/releases/latest` leads (`.../releases/tag/v<version>`), never a pre-release, not from GitHub's API. With no release found it says so and exits 1; before, it went on with an empty version and failed on a download that was not there. Every app copies the new `install.sh` again at its next release (`check` says so); `<app>`, `WINDOWS` and the app's own blocks stay the app's.
- `check`, `tdp_section` and `tdp_file` no longer ask GitHub's API either: the latest tdp comes from `git ls-remote`, and the files of a tdp version from that tag's source archive, still read once and kept.
- Checked against terminu-go v1.3.6 and v1.3.7, which need no other change: a PTY on Windows needs no API, and nothing in harness says a PTY is unix-only; `/terminu:newapp`'s `update` step names no API.

## v1.3.4 — tdp v1.3.5

Follows terminu-go v1.3.5, whose `List.Hold` sets a list's filter aside while an item is held and whose `List.Drop` brings it back.

- `/terminu:newapp`: a list with a Drag mode builds it with `Panel.EnterMode`, picking the item up with `List.Hold`, moving it with `List.MoveHeld` and dropping it with `List.Drop`; putting it back with `Esc` drops it in `OnExit`, then sets the rows as they were with `SetRows`. The app does not move rows itself with `SetRows`, and does not call `ClearFilter` before `Hold`: the step that clears the filter before the app's own actions does not cover a Drag. On a filtered list, setting the filter aside while the item is held and bringing it back as it was, with the cursor on that item, are the lib's (tdp panel-199).
- Checked against tdp v1.3.4 and v1.3.5, which need no other change: how a fixed-width slot is padded is the app's own (L2-2), and the flat start colour with no OSC 11 answer and the order of the OSC 11 query are the lib's (panel-200, color.4). `check`, the reviewer and the templates stay as they are.

## v1.3.3 — tdp v1.3.3

Follows terminu-go v1.3.3, which adds `Mode.PTY` for a mode of the app's own that holds its panel's PTY, `Terminal.Send` and `Panel.EscDesc`.

- `/terminu:newapp`: a panel's PTY gone into and out of as a mode of the app's own, such as an agent the user goes into to type to, is entered with `Panel.EnterMode` and a `Mode` with `PTY` set, not made by the app forwarding keys in a mode of its own. Every key but `Esc` going to the program, `Esc` leaving with no confirm and two Escs read as `Alt-Esc` leaving and sending one Esc are the lib's (tdp K11-28). Whether to have it is the app's pick; going in with `Terminal.Enter` stays as it was. Out of the mode, an `Esc` that cancels what the program is doing is sent with `Terminal.Send`, and any panel whose `OnEsc` does something other than going up a level sets `Panel.EscDesc`, so the key reference does not say "one level up" (tdp K4, M5).
- Checked against tdp v1.3.3, which needs no other change: F6's third case, asking only when a run would carry something plainly wrong into what it makes, leaves the delete-key step as it is, and `check`, the reviewer and the templates stay as they are.

## v1.3.2 — tdp v1.3.2

Follows terminu-go v1.3.2, which adds `List.LocationWork` and `Panel.LocationWork` for work running on the place a location row names, and keeps a panel's key the panel's while its action is hidden.

- `/terminu:newapp`, the delete key: only an action whose point is removing for good goes on `x` or `X`, and it may have no hotkey at all, reached from its menu (tdp M3). One that does something else and loses something on the way, such as abandoning a merge, takes the letter of what it does (`[A]bort`) and asks with `Action.AskFor` only when this run would lose something, the same state always asking the same way (tdp K13-6, F6-6). The step used to put every action that removes something on `x`.
- `/terminu:newapp`, the panel filter: an action the user runs on an item of the filtered list that moves it or takes it away, such as staging or deleting, keeps the filter. The app sets the new rows and puts the cursor on an item the filter shows with `List.Select` (tdp panel-198). The step used to clear the filter before any action that changes what the panel lists.
- `/terminu:newapp`: work running on the place a location row names, such as a push, is written in `LocationWork`, not in `Location`. The Green, the turning loading icon and the row's own glyph coming back are the lib's (tdp panel.10, panel-197).
- `/terminu:newapp`: an action with no target for now, such as Mark all on an empty list, or one that only some tabs have, stays in `Panel.Actions` with `Hidden`. Taken out, its key would run the global action on the same key there (tdp M4-5, M4-6).
- Checked against tdp v1.3.2, which needs no other change: `Enter` by each item as it stands (K3-2 removed) and the cursor going to a new item only for an action whose point is adding it (K3-15) need nothing from the harness, and `check`, the reviewer and the templates stay as they are.

## v1.3.1 — tdp v1.3.1

Follows terminu-go v1.3.1, which adds `Screen.Following`, `Screen.Leads` and `Screen.Hidden` for panels that follow a leading panel.

- `/terminu:newapp`: panels that follow whichever leading panel last had the focus, such as a git tool's right side showing commits for a branch and a diff and history for a changed file, are written in `Screen.Following` and `Screen.Leads`, not swapped by the app as the focus moves (tdp L5, L5-4). The lib keeps them the same each time, changes nothing as the focus moves into them, numbers them by place and leaves a hidden one's digit doing nothing. Whether to have them is the app's pick.
- Checked against tdp v1.3.1, which needs no other change: an item's other rows fading out, the terminal background asked with OSC 11 and a note's opening in the App Log are the lib's, and `check`, the reviewer and the templates stay as they are.

## v1.3.0 — tdp v1.3.0

Follows terminu-go v1.3.0, which adds `SelfUpdate` for `<app> update`, `Panel.Note` for an item's note and `Panel.Sign` for a sign, and stacks toasts, so `App.Notify` no longer merges them.

- `/terminu:newapp`: every app has `<app> update`, run by terminu-go's `SelfUpdate`, never written by the app (tdp E1-9, E1-10). `Owner` and `Tap` follow tdp E6, and an app that changed its release steps keeps the release file names of the family's goreleaser template, which `update` downloads.
- `/terminu:newapp`, upgrading: something the new tdp wants in every app, such as `<app> update`, goes into the task that covers that part, or into a new task whose request is what tdp asks for, so the review sees it.
- `/terminu:newapp`: each E3 notice is a toast of its own, stacked like every toast. The start's notices no longer merge into one, and one found later is a toast of its own too (tdp E3-37; E3-43 is gone).
- `/terminu:newapp`: a value naming something outside the app is used at once by its name. An entry the answer shows gone is warned about and dropped at the next write, so the next start says nothing (tdp E3-41). What can't be used until found holds the default until the finding is done (tdp E3-42). A value naming the user's own data is kept as before.
- `/terminu:newapp`: what an item is, when its columns can't hold it, goes in the item's note (`Panel.Note`, tdp panel.18), not in a toast or the deprecated `Panel.Inspect`. `i` opens notes and nothing else, and `x` and `X` only remove for good: closing what can be opened again takes another letter (tdp K13-4, K13-5). An app with no Diff mode may give `D` to something that isn't a delete.
- `/terminu:newapp`: a panel that only tells is a sign (`Panel.Sign`, tdp panel.17), never a panel that takes the focus. The README sections 「項目的附註」 and 「說明牌」 are now required reading.
- `/terminu:review`: when there is no `terminu:reviewer` agent, as after installing the plugin into a session already open without running `/reload-plugins` afterwards, the review goes to a general-purpose agent instead of being judged in the session itself. The agent is given `agents/reviewer.md` and the same limits: it only reads, with Read, Grep, Glob, `tdp_section` and `tdp_file`, runs no commands and changes no files. The user is told that `/reload-plugins` or a new session brings the reviewer.
- The install notes say to run `/reload-plugins` after installing into a session already open: the skills can show up without it, but the reviewer agent does not.
- Checked against tdp v1.3.0, which needs no other change: stacked toasts, the App Log's short text (tdp note-108) and the bracketed letter (tdp M5-1) are drawn by the lib, and `check`, the reviewer and the templates stay as they are.

## v1.2.13 — tdp v1.2.20

Follows terminu-go v1.2.21, which adds `DiffMode.InDetails` and `DiffMode.Noun` and has `App.Notify` join a late toast to the start's.

- `/terminu:newapp`: anything that removes for good (delete, discard, drop) is on `x` for the item under the cursor and on `X` for many at once, never on `d` or `D`, and the Diff mode's key is `D` (tdp K13). The example in `harness-assets.md`, which deleted on `[D]`, now deletes on `x`.
- `/terminu:newapp`: a value naming something outside the app that is not there, such as a namespace gone from a cluster, is kept like one naming the user's own data, also when the app writes the file again, and a value of several entries uses those still there (tdp E3-41). Something the app finds only after start is not missing until the finding is done (tdp E3-42). Then it goes to `App.Notify`, which joins it to the start's toast while that toast still shows (tdp E3-43).
- `/terminu:newapp`: a Diff mode shown in a list's details panel uses `DiffMode.InDetails`; the app does not draw it there itself (tdp panel-185 to panel-187). Whether to show it there is the app's pick.
- Checked against tdp v1.2.20, which needs no other change: `check`, the reviewer and the templates stay as they are.

## v1.2.12 — tdp v1.2.19

Follows terminu-go v1.2.20, which renames the Compare mode the Diff mode (`DiffMode`, `App.DiffModeAction`, with `Compare` and `CompareAction` kept as deprecated aliases) and adds `List.Spread`.

- `/terminu:newapp` uses the new names: rows that hold many in the Diff mode (`App.DiffModeAction`) use `DiffMode.Many` and `DiffMode.Items`, and the deprecated names are not used, as a later minor version removes them. The required README section 「比較」 is now 「Diff」.
- Checked against tdp v1.2.19, which needs no other change. Where a list's room left over goes (tdp panel-184) is the app's pick, and the default, `SpreadEnd`, is one of the four.

## v1.2.11 — tdp v1.2.18

Follows terminu-go v1.2.19, which adds `Notice.Lines`, `NoLongerUsedLines`, `CannotReadText`, `ProbeResult.Lines`, `Pager.SetDiff`, `Pager.SetDiffItems` and `Row.Below`.

- `/terminu:newapp`: with no screen, each `Notice`'s `Lines()` goes to stderr, each line starting `<app>: ` (tdp E3-40). So do the Info entries for keys not read and for files no longer used (`NoLongerUsedLines`). A file that can't be read is treated as `MoveAside`'s `Esc`: no file is touched, `CannotReadText` goes to stderr and the command exits 1 (tdp E3-27).
- `/terminu:newapp`: a command prints as tdp E1 says. Results go to stdout. Errors and warnings go to stderr, starting `<app>: `, and a wrong use prints the usage and exits 2. The app writes no `log.Fatal` and doesn't leave messages to the `flag` package. `iconwidth` prints `ProbeResult.Lines()`, and the command line is written as `command(args, stdout, stderr) int`, as terminu-go's showcase does, so tests can read what it prints and its exit status. In `harness-assets.md`, a command's implementation records what it prints, where, and its exit status, so the review can judge E1.
- `/terminu:newapp`: a panel showing a difference uses `Pager.SetDiff`, or `Pager.SetDiffItems` for something that holds many, never a tool's own output, and its `Enter` goes through `App.OpenDiff` (tdp note-103, note-106).
- `/terminu:newapp`: an action of the app's own that changes what a panel lists or moves its cursor calls `List.ClearFilter()` first (tdp panel-180). An item taking several rows uses `Row.Below` (tdp panel.5).
- `/terminu:newapp`: content the app reads again on its own is logged once (tdp note-104, note-105). A change found or a recovery gets an Info entry. A failure is said only the first time and when it changes, and a repeated one calls `CannotReadAgain(v)` with no words.
- The README sections "命令列" and "占好幾列的項目" are now required reading. Checked against tdp v1.2.18, which needs no other change: `check`, the reviewer and the templates stay as they are.

## v1.2.10 — tdp v1.2.17

- A family release is tagged with an annotated tag whose message is the release's CHANGELOG heading without `## `, such as `v1.3.0 — tdp v1.2.17`, and a beta's the same way (harness project.7, project-34). Lightweight tags made before this release and already pushed stay as they are.
- `/terminu:release`: step 8 tags with `git tag -a <version> -m "<heading>"`. Step 9 then checks the remote, where `git ls-remote origin "refs/tags/<version>^{}"` prints a line only for an annotated tag.

## v1.2.9 — tdp v1.2.17

Follows terminu-go v1.2.18, which adds `App.OpenDiff`, `Compare.Many` and `Compare.Items`.

- `/terminu:newapp`: outside the TUI (a screen that is not part of it, or output printed for another program to show, such as a status line), every E3 case leaves the file alone (tdp E3). There is no `KeepAside`, `MoveAside` or rewrite, and the defaults hold for this run. The problem is shown there, for example as a mark in the status line, not only on stderr. The next TUI start does the rest. "With no screen" now means a command run in a terminal or a pipe.
- `/terminu:newapp`: a key in a settings file that the app does not read, whether a typo or a key it no longer reads, is not a value it can't use (tdp E3-38). It is not copied aside, collected as a `Notice` or given a toast, and it stays in the file when the app writes the file again. At every start the app logs an Info entry with `App.Log`, naming the key and the file: `No longer used` if the app knows it used to read the key, `Not read` otherwise.
- `/terminu:newapp`: showing two versions of one item, such as a file in a commit against the commit before it, goes through `App.OpenDiff` instead of a viewer showing a patch (tdp note.8). Rows that hold many in a compare mode use `Compare.Many` and `Compare.Items` (tdp panel.9). The list of items that differ is the lib's, and `Panel.LeaveDiff` comes before the app shows something else in that panel. The README section "比較" is now required reading.
- Checked against tdp v1.2.17, which needs no other change.

## v1.2.8 — tdp v1.2.16

Follows terminu-go v1.2.17, which adds `Notice`, `MisfitNotice` and `App.Notify`.

- `/terminu:newapp`: a value naming an entry of the user's own data that is no longer there (tdp E3) is not a value the app can't use, so it is neither copied aside nor written again. The file and the value stay as they are and the built-in default holds for this run. The value is shown as written with its label Red (`Field.Err` in a values panel), and its select opens with `Selected: -1`, so picking any entry is a change.
- `/terminu:newapp`: what E3 tells the user at start comes as one toast (tdp E3-37). The app collects a `Notice` for each case, in the order they happened. These are the location move, `MisfitNotice(path, copied, misfits…)` for values that don't fit, an Error `Notice` with `KeepAside`'s error when the copy can't be made (followed by `MisfitNotice(path, "", …)` for the WARNs alone), and a Warning `Notice` for each value naming an entry that is gone. It then calls `App.Notify` once, with no `App.Toast`, `App.Log` or `DidNotFit` of its own. With no screen, each `Notice` with text is a line on stderr.
- Checked against tdp v1.2.16, which needs no other change.

## v1.2.7 — tdp v1.2.15

Follows terminu-go v1.2.16, which adds `KeepAside`, `App.DidNotFit` and `App.Report`.

- `/terminu:newapp`: a file of the app's own that reads but has values it can't use or entries that don't fit (tdp E3) is first copied aside with `terminu.KeepAside`, written again with the defaults only when that worked, and then reported with `App.DidNotFit` (`DidNotFitText` on stderr with no screen). When the copy can't be made, the file is not written again: the defaults hold for this run, `KeepAside`'s error is the Error toast, and `App.DidNotFit(path, "", …)` logs the WARNs. The Warning toast and the WARNs are `DidNotFit`'s; the app adds none of its own.
- `/terminu:newapp` reads terminu-go's README section on success reports, and a success report the user must read whole that a toast can't hold (tdp toast.2) goes through `App.Report`, not a viewer the app opens or a Success it logs itself.
- Checked against tdp v1.2.15, which needs no other change.

## v1.2.6 — tdp v1.2.14

Follows terminu-go v1.2.15, whose `App.CannotReadAgain` now takes the app's words and shows the Error toast or popup itself.

- `/terminu:newapp`: when a viewer's source cannot be read again, the app calls `App.CannotReadAgain(v, <its words>)`, which shows the Error toast, or the error popup over the selection mode's Space menu or `?` when one is open (tdp note-100); the app shows no toast of its own.
- Checked against tdp v1.2.14, which needs no other change.

## v1.2.5 — tdp v1.2.13

Follows terminu-go v1.2.14, which adds `App.ReadAgain` and `Pager.ReadAgain` for content read again while it is shown (tdp note.9).

- `/terminu:newapp` reads terminu-go's README section on content that changes while it is open, and puts a source read again (a watched file, say) into an open viewer or panel with `App.ReadAgain` or `Pager.ReadAgain`, not by replacing it with `SetLines` or opening it again, so search, the selection and following the end behave as tdp note.9 says. Such content is not streaming, so `Panel.Stream` stays off (tdp T2). A `Pager` that takes new lines while it is open (`Append` or `ReadAgain`) sets `Follow`, or it does not follow its end (tdp note-90).
- Checked against tdp v1.2.7 to v1.2.13, which need no other change.

## v1.2.4 — tdp v1.2.6

Follows terminu-go v1.2.7, the first terminu-go release whose number differs from its tdp's: it still follows tdp v1.2.6 and only adds `App.RunAction`.

- A terminu-go release can follow the same tdp as the one before. `check` already reads the tdp from terminu-go's CHANGELOG heading (`## v1.2.7 — tdp v1.2.6`), so it needs no change. `/terminu:newapp`'s upgrade, when the newer terminu-go follows the tdp the spec is on, opens no new spec: it reads only terminu-go's CHANGELOG, changes what that says an app must change, and commits once `make check` passes. "Carry on" tells you what the newer terminu-go adds and whether the app must change, not what tdp changed. `/terminu:analysis` skips summing up tdp's CHANGELOG when the tdp is the same. The README no longer says each tdp version has exactly one terminu-go version.
- `/terminu:newapp` reads terminu-go's README section on running actions from a popup of the app's own, and runs an app action from such a popup with `App.RunAction` (tdp popup.8) rather than redoing it by hand. From a popup opened in the app's own mode the action leaves the mode without running `OnExit`, so what `OnExit` would do (pausing, say) and asking before losing what the user made (tdp K11) come before the call.
- harness-assets (work.yaml): an app's own mode lists one more way it is left, an app action run from a popup opened in the mode (tdp popup-74).

## v1.2.3 — tdp v1.2.6

- harness-assets (spec.yaml, work.yaml): tdp v1.2.6 makes an app's own mode one of P5's openings. A task that uses one cites `tdp P5` and `tdp K11` for the mode's keys and content, and its implementation says what content the mode is on, how it is entered, and each way it is left (`Esc`, a key of the mode, the app ending it), so the review can judge the door and the frame.

## v1.2.2 — tdp v1.2.5

- harness project.3: besides the files the tree names, what else goes in `docs/` (e.g. the social preview image GitHub shows when the repo's link is shared) the app may decide for itself, as with the top-level files the tree does not name (issue #1, from tetru). Checked against tdp v1.2.1 to v1.2.5, which need no other change.

## v1.2.1 — tdp v1.2.0

- The family no longer offers `go install` as a way to install an app: a binary it builds gets no version from goreleaser and prints `dev`. harness project.3 drops the `go install` path and keeps the `/vN` rule, for the module's versions and the binaries' build information, however the app is installed; project.7 names the ways in (brew, `install.sh`, scoop and `install.ps1` with a Windows build, a download from the GitHub release) and has a beta downloaded from its GitHub pre-release. `project-13`, its ID kept, checks the `/vN` without `go install`. `/terminu:release`, moving to a new major, changes the module and the imports; the README has no `go install` path to change. `check`'s message on a module without its `/vN` says why without `go install`.

## v1.2.0 — tdp v1.2.0

Everything since v1.1.0, the betas v1.2.0-beta.1 to beta.3 included:

- The family's repo templates move here from tdp: tdp says what an app should be, whatever the language, and templates written for Go and goreleaser are the family's upkeep, kept with the tool that makes and releases the family's apps. `project/` holds terminu project as tdp v1.1.10 had it: its page (sections project.1–8), its checklist items (the 38 `project-` items, their IDs kept) and the seven templates, whose first comment now cites `harness project/…`.
- Citing them: `harness project.3`, `harness project` (the whole page), `harness project/install.sh`, `harness project-13`, in the same `tdp` list of `spec.yaml`. `tdp_section` reads them; `check` checks them against `project/`, and a `tdp project…` citation the app's tdp no longer has says to write it with `harness`. An app whose tdp still has `project/` (v1.1.10 and before) may cite it as before.
- `check` compares a family member's repo with the templates, the files, the Makefile and the CHANGELOG, once the tdp it follows has no `project/`: what terminu-go's tdptest checked until then, moved here with the templates, with messages in Traditional Chinese.
- Four backings in `work.yaml`: a test, `check: repo`, `user: {date, how}`, `commit:`. `check: repo` means harness's `check` guards the task, reading a family member's repo against the templates; it counts only where harness reads the repo (a family member's, its tdp without `project/`). A subtest named as a backing (`TestX/sub`) has to run: `check` runs it with `go test -run` and looks for it in the output, so a subtest that is gone (terminu-go's `TestTDP/repo`, say) no longer passes, and a gone `…/repo` says to write `check: repo`.
- `check` no longer passes what it could not check: tdp it could not read, tests that would not run head the report as not checked, and the command exits 1 (the MCP tool reports an error). Reading tdp from GitHub tries again after a network error, a 5xx or a 429. The tool's description says all of this.
- The newest terminu-go no longer counts a beta, as Go's `@latest` does not: `/terminu`, `/terminu:newapp`, `/terminu:release` and `/terminu:analysis` take the newest release (`go list -m …@latest`), a beta only when the app is on a beta already (the last of `go list -m -versions`, so a release replaces its betas) or when you name one. `tdp_file` and `tdp_section`, given no version, read the newest tdp release rather than a beta.
- `/terminu:newapp`, `/terminu:analysis` and `/terminu:release` read the templates from harness. The first upgrade to a tdp without `project/` rewrites the citations and the templates' first comments and swaps a `TestTDP/repo` backing for `check: repo`; `/terminu:analysis` writes `check: repo` for the repo's task. `/terminu:release` reads the release rules from tdp's `project/` while the app's tdp still has it, from harness's after that, and runs `check` after `make check`.
- `/terminu:newapp`'s upgrade gives a task the new spec id only when the new tdp changes it, its record or what was built from it; a task that cites a changed section for something else keeps its id, as each review judges every task against the new tdp anyway.
- `/terminu:review` handles a `check` that is not complete by checking again first.
- `/terminu:reportgap` reports a problem with the family's templates to terminu-harness.
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

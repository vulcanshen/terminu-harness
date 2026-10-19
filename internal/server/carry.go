package server

import (
	"fmt"

	"github.com/vulcanshen/terminu-harness/internal/appgit"
	"github.com/vulcanshen/terminu-harness/internal/run"
	"github.com/vulcanshen/terminu-harness/internal/tdp"
)

// carryOver keeps what the last run of the app found where nothing it looked
// at changed (user, 2026-10-19: judge again only what changed, roughly). With
// a repo, the run records its commit and terminu-go version; when an earlier
// finished run for the same tdp and terminu-go has a commit, the files
// changed since decide: changed code judges every item but those read from
// files again, changed other files those read from files, and the rest keep
// the earlier verdicts. It says what it did, "" when it did nothing.
func carryOver(r *run.Run, base, repo, tag string) (string, error) {
	if repo == "" {
		return "", nil
	}
	commit := appgit.Commit(repo)
	if commit == "" {
		return "", nil
	}
	tgo := appgit.Tgo(repo)
	if err := r.SetRepoState(commit, tgo); err != nil {
		return "", err
	}
	prev := run.Latest(base, r.Meta.App, tag, r.Dir)
	if prev == nil {
		return "no earlier run to carry from: every item is judged\n", nil
	}
	if prev.Meta.Tgo != tgo {
		return fmt.Sprintf("terminu-go changed since %s (%s → %s): every item is judged\n", prev.Dir, prev.Meta.Tgo, tgo), nil
	}
	changed, err := appgit.Changed(repo, prev.Meta.Commit)
	if err != nil {
		return fmt.Sprintf("can't tell what changed since %s: every item is judged\n", prev.Meta.Commit), nil
	}
	code, other := false, false
	for _, f := range changed {
		if appgit.IsCode(f) {
			code = true
		} else {
			other = true
		}
	}
	n, err := r.Carry(prev, func(it tdp.Item) bool {
		if run.ReadsFiles(it) {
			return !other
		}
		return !code
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("carried %d verdicts from %s (changed since: code %v, other files %v); judge only what judge_status lists as not judged\n",
		n, prev.Dir, code, other), nil
}

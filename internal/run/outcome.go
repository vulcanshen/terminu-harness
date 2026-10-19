package run

import (
	"os"
	"path/filepath"
	"regexp"
)

// The outcomes of a judging: harness's line for whether an app is done
// (user, 2026-10-19: judging is the harness's, tdp only says what an app is).
const (
	Passed  = "通過"  // no item violated
	Pending = "待修"  // every violation is listed in the app's terminu-fix file, waiting on tdp or terminu-go
	Failed  = "不通過" // a violation is listed nowhere: the app has work of its own left
)

// FixFile is where an app lists what it does not meet yet (tdp E7).
func (r *Run) FixFile() string {
	if r.Meta.Repo == "" {
		return ""
	}
	return filepath.Join(r.Meta.Repo, "docs", r.Meta.App+"-terminu-fix.md")
}

// Outcome is the run's outcome over the violated items, and the violated
// items the fix file does not name. Without a repo there is no fix file to
// read, so any violation fails.
func (r *Run) Outcome(violated []string) (string, []string) {
	if len(violated) == 0 {
		return Passed, nil
	}
	fix := ""
	if f := r.FixFile(); f != "" {
		if b, err := os.ReadFile(f); err == nil {
			fix = string(b)
		}
	}
	var unlisted []string
	for _, id := range violated {
		if !regexp.MustCompile(`(^|[^A-Za-z0-9-])` + regexp.QuoteMeta(id) + `($|[^A-Za-z0-9-])`).MatchString(fix) {
			unlisted = append(unlisted, id)
		}
	}
	if len(unlisted) > 0 {
		return Failed, unlisted
	}
	return Pending, nil
}

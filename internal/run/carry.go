package run

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/vulcanshen/terminu-harness/internal/tdp"
)

// ReadsFiles reports an item judged by reading files, not the screen or the
// code: its way to check starts with 讀檔案.
func ReadsFiles(it tdp.Item) bool { return strings.HasPrefix(it.HowTo, "讀檔案") }

// SetRepoState records the app's commit and terminu-go version the run judges.
func (r *Run) SetRepoState(commit, tgo string) error {
	r.Meta.Commit, r.Meta.Tgo = commit, tgo
	return r.saveMeta()
}

func (r *Run) saveMeta() error {
	b, _ := json.MarshalIndent(r.Meta, "", "  ")
	return os.WriteFile(filepath.Join(r.Dir, metaFile), b, 0o644)
}

// Latest is the newest finished run of the app under base for the same tdp,
// with the app's commit recorded, other than skip; nil for none.
func Latest(base, app, tag, skip string) *Run {
	dirs, _ := filepath.Glob(filepath.Join(base, app, "*-tdp-"+strings.TrimPrefix(tag, "v")))
	var runs []*Run
	for _, d := range dirs {
		if filepath.Clean(d) == filepath.Clean(skip) {
			continue
		}
		if _, err := os.Stat(filepath.Join(d, ReportFile)); err != nil {
			continue
		}
		if r, err := Open(d); err == nil && r.Meta.Commit != "" {
			runs = append(runs, r)
		}
	}
	if len(runs) == 0 {
		return nil
	}
	sort.Slice(runs, func(i, j int) bool { return runs[i].Meta.Started.After(runs[j].Meta.Started) })
	return runs[0]
}

// Carry copies prev's verdicts on the items keep says nothing changed for,
// the reason saying where they come from and the evidence pointing into
// prev's directory. It reports how many it copied.
func (r *Run) Carry(prev *Run, keep func(tdp.Item) bool) (int, error) {
	judged, err := prev.Judged()
	if err != nil {
		return 0, err
	}
	rel, err := filepath.Rel(r.Dir, prev.Dir)
	if err != nil {
		return 0, err
	}
	short := prev.Meta.Commit
	if len(short) > 7 {
		short = short[:7]
	}
	n := 0
	for _, it := range r.items {
		v, ok := judged[it.ID]
		if !ok || !keep(it) {
			continue
		}
		var ev []string
		for _, e := range v.Evidence {
			if m, _ := filepath.Glob(filepath.Join(prev.Dir, "evidence", e+".*")); len(m) > 0 {
				ev = append(ev, filepath.ToSlash(filepath.Join(rel, "evidence", filepath.Base(m[0]))))
			}
		}
		v.Evidence = ev
		v.Reason = fmt.Sprintf("沿用 %s 的判定：%s", short, v.Reason)
		if err := appendJSON(filepath.Join(r.Dir, verdictFile), v); err != nil {
			return n, err
		}
		n++
	}
	r.Meta.Previous, r.Meta.Carried = prev.Dir, n
	return n, r.saveMeta()
}

// Items is the checklist the run judges by.
func (r *Run) Items() []tdp.Item { return r.items }

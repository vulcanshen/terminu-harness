// Package run keeps one judging of one app against one version of tdp: what
// the judges found, item by item, and the report card made from it.
package run

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/vulcanshen/terminu-harness/internal/tdp"
)

// The five verdicts (checklist README, 判法). Only 符合 is a gene.
const (
	Meets       = "符合"
	Departs     = "偏離"
	Violates    = "違反"
	NotApply    = "不適用"
	CannotTell  = "判不了"
	notJudged   = "沒判到"
	itemsFile   = "items-zh_TW.md"
	metaFile    = "run.json"
	verdictFile = "verdicts.jsonl"
	gapFile     = "gaps.jsonl"
)

// Verdicts in the order a report lists them.
var Verdicts = []string{Meets, Departs, Violates, NotApply, CannotTell}

// Meta is what the run is about.
type Meta struct {
	App     string       `json:"app"`
	Command string       `json:"command"`
	Repo    string       `json:"repo,omitempty"`
	TDP     string       `json:"tdp"`
	Started time.Time    `json:"started"`
	Batches []BatchFiles `json:"batches"`
	// Commit and Tgo are the app's commit and terminu-go version, when there
	// is a repo; Previous is the run whose verdicts the Carried items reuse.
	Commit   string `json:"commit,omitempty"`
	Tgo      string `json:"tgo,omitempty"`
	Previous string `json:"previous,omitempty"`
	Carried  int    `json:"carried,omitempty"`
}

// BatchFiles is one batch and the file listing its items for its judge.
type BatchFiles struct {
	ID       string   `json:"id"`
	Sections []string `json:"sections"`
	Count    int      `json:"count"`
	File     string   `json:"file"`
}

// Verdict is a judge's verdict on one item.
type Verdict struct {
	ID      string `json:"id"`
	Verdict string `json:"verdict"`
	Reason  string `json:"reason"`
	// Cause is why an item can't be told (user, 2026-10-19), where it goes:
	// tdp or terminu-go are reported back to terminu; other is put to the
	// user in the app's session.
	Cause    string    `json:"cause,omitempty"`
	Evidence []string  `json:"evidence,omitempty"`
	At       time.Time `json:"at"`
}

// The causes of 判不了.
const (
	CauseTDP   = "tdp"   // tdp says nothing, or nothing clear
	CauseTgo   = "tgo"   // terminu-go can't do it, or gets in the way
	CauseOther = "other" // anything else: no repo, a special environment
)

// Causes in the order a report lists them.
var Causes = []string{CauseTDP, CauseTgo, CauseOther}

// Gap is something the app does that tdp does not cover yet.
type Gap struct {
	What     string    `json:"what"`
	Where    string    `json:"where,omitempty"`
	Evidence []string  `json:"evidence,omitempty"`
	At       time.Time `json:"at"`
}

// Run is one judging, kept in its own directory.
type Run struct {
	Dir   string
	Meta  Meta
	items []tdp.Item
	index map[string]int

	mu       sync.Mutex
	evidence int
}

var unsafeName = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

// Start makes the run's directory under base: <app>/<time>-tdp-<version>, with
// the checklist it judges by, one file of items per batch and the meta.
func Start(base, command, repo, tag string, itemsMD []byte, batches int) (*Run, error) {
	secs, err := tdp.Parse(string(itemsMD))
	if err != nil {
		return nil, fmt.Errorf("the checklist of %s: %w", tag, err)
	}
	app := unsafeName.ReplaceAllString(filepath.Base(strings.Fields(command)[0]), "_")
	now := time.Now()
	stamp := now.Format("20060102-150405")
	dir := filepath.Join(base, app, stamp+"-tdp-"+strings.TrimPrefix(tag, "v"))
	for n := 2; ; n++ { // two runs started in the same second
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			break
		}
		dir = filepath.Join(base, app, fmt.Sprintf("%s-%d-tdp-%s", stamp, n, strings.TrimPrefix(tag, "v")))
	}
	if err := os.MkdirAll(filepath.Join(dir, "evidence"), 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(dir, "batches"), 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(dir, itemsFile), itemsMD, 0o644); err != nil {
		return nil, err
	}
	meta := Meta{App: app, Command: command, Repo: repo, TDP: tag, Started: now}
	for _, b := range tdp.Batches(secs, batches) {
		file := filepath.Join("batches", b.ID+".md")
		if err := os.WriteFile(filepath.Join(dir, file), []byte(batchMarkdown(tag, b)), 0o644); err != nil {
			return nil, err
		}
		meta.Batches = append(meta.Batches, BatchFiles{ID: b.ID, Sections: b.Sections, Count: len(b.Items), File: file})
	}
	b, _ := json.MarshalIndent(meta, "", "  ")
	if err := os.WriteFile(filepath.Join(dir, metaFile), b, 0o644); err != nil {
		return nil, err
	}
	return Open(dir)
}

func batchMarkdown(tag string, b tdp.Batch) string {
	var s strings.Builder
	fmt.Fprintf(&s, "# %s：tdp %s 的 %d 條\n\n", b.ID, tag, len(b.Items))
	section := ""
	for _, it := range b.Items {
		if it.Section != section {
			section = it.Section
			fmt.Fprintf(&s, "\n## %s\n\n| 編號 | 出處 | 檢查什麼 | 適用 | 怎麼驗 |\n|---|---|---|---|---|\n", section)
		}
		esc := func(v string) string { return strings.ReplaceAll(v, "|", `\|`) }
		fmt.Fprintf(&s, "| %s | %s | %s | %s | %s |\n", it.ID, esc(it.Source), esc(it.Check), esc(it.Applies), esc(it.HowTo))
	}
	return s.String()
}

// Open reads a run's directory.
func Open(dir string) (*Run, error) {
	b, err := os.ReadFile(filepath.Join(dir, metaFile))
	if err != nil {
		return nil, fmt.Errorf("%s is not a run: %w", dir, err)
	}
	r := &Run{Dir: dir, index: map[string]int{}}
	if err := json.Unmarshal(b, &r.Meta); err != nil {
		return nil, err
	}
	md, err := os.ReadFile(filepath.Join(dir, itemsFile))
	if err != nil {
		return nil, err
	}
	secs, err := tdp.Parse(string(md))
	if err != nil {
		return nil, err
	}
	for _, s := range secs {
		for _, it := range s.Items {
			r.index[it.ID] = len(r.items)
			r.items = append(r.items, it)
		}
	}
	entries, _ := os.ReadDir(filepath.Join(dir, "evidence"))
	for _, e := range entries {
		var n int
		if _, err := fmt.Sscanf(e.Name(), "e%d", &n); err == nil && n > r.evidence {
			r.evidence = n
		}
	}
	return r, nil
}

// Total is how many items the checklist has: the most genes there are.
func (r *Run) Total() int { return len(r.items) }

// NewEvidence names the next piece of evidence and the file it goes in.
func (r *Run) NewEvidence(ext string) (id, path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.evidence++
	id = fmt.Sprintf("e%04d", r.evidence)
	return id, filepath.Join(r.Dir, "evidence", id+ext)
}

func (r *Run) hasEvidence(id string) bool {
	m, _ := filepath.Glob(filepath.Join(r.Dir, "evidence", id+".*"))
	return len(m) > 0
}

// SetProfile keeps what the judge learned about the app.
func (r *Run) SetProfile(text string) error {
	return os.WriteFile(filepath.Join(r.Dir, "profile.md"), []byte(text), 0o644)
}

// AddVerdicts keeps the verdicts that are well formed and says what is wrong
// with the rest. A later verdict on an item replaces an earlier one.
func (r *Run) AddVerdicts(vs []Verdict) (kept int, problems []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	f, err := os.OpenFile(filepath.Join(r.Dir, verdictFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return 0, []string{err.Error()}
	}
	defer f.Close()
	for _, v := range vs {
		v.ID = strings.TrimSpace(v.ID)
		switch {
		case !r.known(v.ID):
			problems = append(problems, fmt.Sprintf("%s: no such item in tdp %s", v.ID, r.Meta.TDP))
			continue
		case !valid(v.Verdict):
			problems = append(problems, fmt.Sprintf("%s: %q is not one of %s", v.ID, v.Verdict, strings.Join(Verdicts, "、")))
			continue
		case strings.TrimSpace(v.Reason) == "":
			problems = append(problems, fmt.Sprintf("%s: no reason", v.ID))
			continue
		case v.Verdict == CannotTell && v.Cause != CauseTDP && v.Cause != CauseTgo && v.Cause != CauseOther:
			problems = append(problems, fmt.Sprintf("%s: why it can't be told: cause is tdp, tgo or other", v.ID))
			continue
		}
		for _, e := range v.Evidence {
			if !r.hasEvidence(e) {
				problems = append(problems, fmt.Sprintf("%s: evidence %s not found (kept the verdict)", v.ID, e))
			}
		}
		v.At = time.Now()
		b, _ := json.Marshal(v)
		if _, err := f.Write(append(b, '\n')); err != nil {
			return kept, append(problems, err.Error())
		}
		kept++
	}
	return kept, problems
}

func (r *Run) known(id string) bool { _, ok := r.index[id]; return ok }

func valid(v string) bool {
	for _, w := range Verdicts {
		if v == w {
			return true
		}
	}
	return false
}

// AddGap keeps one thing tdp does not cover.
func (r *Run) AddGap(g Gap) error {
	if strings.TrimSpace(g.What) == "" {
		return errors.New("what is empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	g.At = time.Now()
	return appendJSON(filepath.Join(r.Dir, gapFile), g)
}

func appendJSON(path string, v any) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	b, _ := json.Marshal(v)
	_, err = f.Write(append(b, '\n'))
	return err
}

// Judged is the latest verdict on each item.
func (r *Run) Judged() (map[string]Verdict, error) {
	out := map[string]Verdict{}
	err := readLines(filepath.Join(r.Dir, verdictFile), func(b []byte) error {
		var v Verdict
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		out[v.ID] = v
		return nil
	})
	return out, err
}

// Gaps is what the judges found tdp does not cover.
func (r *Run) Gaps() ([]Gap, error) {
	var out []Gap
	err := readLines(filepath.Join(r.Dir, gapFile), func(b []byte) error {
		var g Gap
		if err := json.Unmarshal(b, &g); err != nil {
			return err
		}
		out = append(out, g)
		return nil
	})
	return out, err
}

func readLines(path string, f func([]byte) error) error {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	for sc.Scan() {
		if len(strings.TrimSpace(sc.Text())) == 0 {
			continue
		}
		if err := f(sc.Bytes()); err != nil {
			return err
		}
	}
	return sc.Err()
}

// Missing is, batch by batch, the items no judge has given a verdict yet.
func (r *Run) Missing() (map[string][]string, error) {
	judged, err := r.Judged()
	if err != nil {
		return nil, err
	}
	out := map[string][]string{}
	secBatch := map[string]string{}
	for _, b := range r.Meta.Batches {
		for _, s := range b.Sections {
			secBatch[s] = b.ID
		}
	}
	for _, it := range r.items {
		if _, ok := judged[it.ID]; !ok {
			b := secBatch[it.Section]
			out[b] = append(out[b], it.ID)
		}
	}
	return out, nil
}

// Counts is how many items have each verdict, 沒判到 for none.
func (r *Run) Counts() (map[string]int, error) {
	judged, err := r.Judged()
	if err != nil {
		return nil, err
	}
	c := map[string]int{}
	for _, it := range r.items {
		if v, ok := judged[it.ID]; ok {
			c[v.Verdict]++
		} else {
			c[notJudged]++
		}
	}
	return c, nil
}

// Status says how far the judging has come.
func (r *Run) Status() (string, error) {
	c, err := r.Counts()
	if err != nil {
		return "", err
	}
	missing, err := r.Missing()
	if err != nil {
		return "", err
	}
	var s strings.Builder
	fmt.Fprintf(&s, "%s, tdp %s: %d items\n", r.Meta.App, r.Meta.TDP, r.Total())
	for _, v := range append(append([]string{}, Verdicts...), notJudged) {
		fmt.Fprintf(&s, "  %s %d\n", v, c[v])
	}
	ids := make([]string, 0, len(missing))
	for b := range missing {
		ids = append(ids, b)
	}
	sort.Slice(ids, func(i, j int) bool { return batchNum(ids[i]) < batchNum(ids[j]) })
	for _, b := range ids {
		fmt.Fprintf(&s, "not judged in %s (%d): %s\n", b, len(missing[b]), strings.Join(missing[b], " "))
	}
	return s.String(), nil
}

func batchNum(id string) int {
	var n int
	fmt.Sscanf(id, "b%d", &n)
	return n
}

// Genes is the score: tdp-<version> <n> gene.
func Genes(tag string, n int) string {
	return fmt.Sprintf("tdp-%s %d gene", strings.TrimPrefix(tag, "v"), n)
}

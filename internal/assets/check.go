package assets

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/vulcanshen/terminu-harness/internal/tdp"
)

// TGO is terminu-go's module path.
const TGO = "github.com/vulcanshen/terminu-go"

// Env is what Check needs from outside the files in Dir.
type Env struct {
	// Versions is the terminu-go version the app requires and the tdp
	// version that terminu-go version follows.
	Versions func(repo string) (tgo, tdp string, err error)
	// IDs is what can be cited in tdp at a version.
	IDs func(tdpVersion string) (tdp.IDs, error)
	// Tests is the names of the app's tests; nil skips checking them.
	Tests func(repo string) (map[string]bool, error)
	// Commit reports whether the repo has the commit.
	Commit func(repo, hash string) bool
}

// Default reads go.mod and terminu-go's CHANGELOG for the versions, tdp from
// src, and runs go test -list and git in the repo.
func Default(src tdp.Source) Env {
	return Env{Versions: Versions, IDs: src.IDs, Tests: Tests, Commit: Commit}
}

// Problem is something in the files to fix.
type Problem struct {
	File string // a file in Dir, or go.mod
	Task int    // 0 when it is not about one task
	Msg  string
}

func (p Problem) String() string {
	if p.Task > 0 {
		return fmt.Sprintf("%s task %d：%s", p.File, p.Task, p.Msg)
	}
	return fmt.Sprintf("%s：%s", p.File, p.Msg)
}

// Report is what Check found: problems to fix, and what is left to do.
type Report struct {
	TGO, TDP string
	Problems []Problem
	Todo     []string
}

// ReasonKinds are how a reason in review.yaml starts.
var ReasonKinds = []string{"違反", "沒有測試", "tdp 不清楚"}

var dateRE = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// Check checks the files in Dir in the repo against each other, against the
// tdp the app follows, and against the app's tests and history.
func Check(repo string, env Env) Report {
	var r Report
	add := func(file string, task int, format string, a ...any) {
		r.Problems = append(r.Problems, Problem{file, task, fmt.Sprintf(format, a...)})
	}
	f, err := Load(repo)
	if err != nil {
		add(Dir, 0, "%v", err)
		return r
	}
	if strings.TrimSpace(f.Idea) == "" {
		add(IdeaFile, 0, "沒有這個檔，或是空的：使用者描述的需求寫在這裡")
	}
	if f.Spec == nil {
		add(SpecFile, 0, "沒有這個檔")
		return r
	}
	if f.Work == nil {
		add(WorkFile, 0, "沒有這個檔：還沒做任何 task 時寫 tasks: []")
	}

	var ids *tdp.IDs
	tgo, tdpVersion, err := env.Versions(repo)
	if err != nil {
		add("go.mod", 0, "%v", err)
	} else {
		r.TGO, r.TDP = tgo, tdpVersion
		if f.Spec.To != tdpVersion {
			add(SpecFile, 0, "to 是 %q，但 go.mod 的 terminu-go %s 對到 tdp %s：升版時 spec 的 id 加一，from 改成原本的 to，to 改成 %s",
				f.Spec.To, tgo, tdpVersion, tdpVersion)
		}
		got, err := env.IDs(tdpVersion)
		if err != nil {
			add(SpecFile, 0, "讀不到 tdp %s 的編號，tdp 欄位沒有檢查：%v", tdpVersion, err)
		} else {
			ids = &got
		}
	}

	checkSpec(f.Spec, ids, add)
	spec := map[int]SpecTask{}
	for _, t := range f.Spec.Tasks {
		spec[t.ID] = t
	}

	work := map[int]WorkTask{}
	if f.Work != nil {
		var tests map[string]bool
		if env.Tests != nil {
			if tests, err = env.Tests(repo); err != nil {
				add(WorkFile, 0, "go test -list 跑不起來，測試的名字沒有檢查：%v", err)
				tests = nil
			}
		}
		for _, w := range f.Work.Tasks {
			if _, dup := work[w.ID]; dup {
				add(WorkFile, w.ID, "id 重複")
				continue
			}
			work[w.ID] = w
			checkWork(repo, w, spec, tests, env.Commit, add)
		}
	}

	if f.Review != nil {
		checkReview(f.Review, f.Spec.ID, spec, add)
	}

	for _, t := range f.Spec.Tasks {
		if t.Canceled != "" {
			continue
		}
		w, ok := work[t.ID]
		switch {
		case !ok:
			r.Todo = append(r.Todo, fmt.Sprintf("task %d：還沒做（spec %d）", t.ID, t.Spec))
		case w.Spec < t.Spec:
			r.Todo = append(r.Todo, fmt.Sprintf("task %d：需求在 spec %d 改過，work 還是照 spec %d 做的", t.ID, t.Spec, w.Spec))
		}
	}
	switch {
	case f.Review == nil:
		r.Todo = append(r.Todo, "還沒 review 過")
	case f.Review.Spec < f.Spec.ID:
		r.Todo = append(r.Todo, fmt.Sprintf("上次 review 對的是 spec %d，現在是 spec %d", f.Review.Spec, f.Spec.ID))
	}
	return r
}

type addFunc func(file string, task int, format string, a ...any)

func checkSpec(s *Spec, ids *tdp.IDs, add addFunc) {
	if s.ID < 1 {
		add(SpecFile, 0, "id 從 1 開始")
	}
	if s.From != "" && s.From != "0" && !isVersion(s.From) {
		add(SpecFile, 0, "from 是 %q：寫 tdp 的版本（vX.Y.Z），第一份 spec 寫 0", s.From)
	}
	if !isVersion(s.To) {
		add(SpecFile, 0, "to 是 %q：寫 tdp 的版本（vX.Y.Z）", s.To)
	}
	seen := map[int]bool{}
	for _, t := range s.Tasks {
		if t.ID < 1 {
			add(SpecFile, 0, "task 的 id 從 1 開始（%q）", t.Request)
			continue
		}
		if seen[t.ID] {
			add(SpecFile, t.ID, "id 重複：task 的 id 整個 app 不重複使用")
		}
		seen[t.ID] = true
		if t.Spec < 1 || t.Spec > s.ID {
			add(SpecFile, t.ID, "spec 是 %d：寫最後新增或修改它的 spec（1～%d）", t.Spec, s.ID)
		}
		if strings.TrimSpace(t.Request) == "" {
			add(SpecFile, t.ID, "沒有 request")
		}
		reported := false
		for _, rm := range t.Remark {
			if strings.TrimSpace(rm.Note) == "" && strings.TrimSpace(rm.Issue) == "" {
				add(SpecFile, t.ID, "remark 有一筆是空的")
			}
			if strings.TrimSpace(rm.Issue) != "" {
				reported = true
			}
		}
		if len(t.TDP) == 0 && !reported && t.Canceled == "" {
			add(SpecFile, t.ID, "沒有 tdp 根據：在 tdp 找不到根據時回報 tdp，把 issue 記在 remark")
		}
		for _, c := range t.TDP {
			id, ok := strings.CutPrefix(c, "tdp ")
			id = strings.TrimSpace(id)
			switch {
			case !ok || id == "":
				add(SpecFile, t.ID, "tdp 寫成 %q：一律用 tdp 加編號，例 \"tdp K4\"、\"tdp panel.5\"、\"tdp confirm\"、\"tdp project/install.sh\"、\"tdp K9-6\"", c)
			case ids != nil && !ids.Has(id):
				add(SpecFile, t.ID, "%q 在這一版 tdp 裡找不到", c)
			}
		}
	}
}

func checkWork(repo string, w WorkTask, spec map[int]SpecTask, tests map[string]bool, commit func(repo, hash string) bool, add addFunc) {
	if w.ID < 1 {
		add(WorkFile, 0, "task 的 id 從 1 開始")
		return
	}
	t, ok := spec[w.ID]
	if !ok {
		add(WorkFile, w.ID, "spec.yaml 沒有這個 task")
		return
	}
	if w.Spec < 1 || w.Spec > t.Spec {
		add(WorkFile, w.ID, "spec 是 %d：寫這筆實作照的是哪一份 spec（1～%d）", w.Spec, t.Spec)
	}
	if (t.Canceled != "") != w.Canceled {
		if w.Canceled {
			add(WorkFile, w.ID, "標了 canceled，但 spec.yaml 的這個 task 沒有取消")
		} else {
			add(WorkFile, w.ID, "spec.yaml 取消了這個 task，這裡也要標 canceled: true")
		}
	}
	if w.Canceled {
		return
	}
	if strings.TrimSpace(w.Implementation) == "" {
		add(WorkFile, w.ID, "沒有 implementation")
	}
	if len(w.Test) == 0 {
		add(WorkFile, w.ID, "沒有測試：test 寫測試的名字；寫不成測試的寫使用者確認（user: {date, how}）")
	}
	for _, b := range w.Test {
		n := 0
		for _, set := range []bool{b.Test != "", b.Commit != "", b.User != nil} {
			if set {
				n++
			}
		}
		if n != 1 {
			add(WorkFile, w.ID, "test 的每一筆只寫一種：測試的名字、commit、或 user")
			continue
		}
		switch {
		case b.Test != "":
			name, _, _ := strings.Cut(b.Test, "/")
			if tests != nil && !tests[name] {
				add(WorkFile, w.ID, "測試 %s 不在 go test -list 裡", b.Test)
			}
		case b.Commit != "":
			if commit != nil && !commit(repo, b.Commit) {
				add(WorkFile, w.ID, "commit %s 不在這個 repo 裡", b.Commit)
			}
		case b.User != nil:
			if _, err := time.Parse("2006-01-02", b.User.Date); err != nil || !dateRE.MatchString(b.User.Date) {
				add(WorkFile, w.ID, "使用者確認的 date 是 %q：寫 YYYY-MM-DD", b.User.Date)
			}
			if strings.TrimSpace(b.User.How) == "" {
				add(WorkFile, w.ID, "使用者確認沒有寫 how（怎麼確認）")
			}
		}
	}
}

func checkReview(rv *Review, specID int, spec map[int]SpecTask, add addFunc) {
	if rv.Spec < 1 || rv.Spec > specID {
		add(ReviewFile, 0, "spec 是 %d：寫這次 review 對的 spec（1～%d）", rv.Spec, specID)
	}
	seen := map[int]bool{}
	for _, t := range rv.Tasks {
		if seen[t.ID] {
			add(ReviewFile, t.ID, "id 重複")
		}
		seen[t.ID] = true
		if _, ok := spec[t.ID]; !ok {
			add(ReviewFile, t.ID, "spec.yaml 沒有這個 task")
		}
		if t.Accepted {
			continue
		}
		kind := false
		for _, k := range ReasonKinds {
			if strings.HasPrefix(t.Reason, k+"：") || strings.HasPrefix(t.Reason, k+":") {
				kind = true
			}
		}
		if !kind {
			add(ReviewFile, t.ID, "reason 開頭寫是哪一種：%s，接「：」", strings.Join(ReasonKinds, "／"))
		}
	}
}

func isVersion(v string) bool {
	return versionRE.MatchString(v)
}

var (
	versionRE   = regexp.MustCompile(`^v\d+\.\d+\.\d+(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`)
	requireRE   = regexp.MustCompile(`(?m)^\s*(?:require\s+)?` + regexp.QuoteMeta(TGO) + `\s+(v\S+)`)
	changelogRE = regexp.MustCompile(`(?m)^## (v\S+) — tdp (v\S+)`)
)

// Versions reads the terminu-go version in the repo's go.mod, and the tdp
// version that version follows from terminu-go's CHANGELOG.
func Versions(repo string) (string, string, error) {
	b, err := os.ReadFile(filepath.Join(repo, "go.mod"))
	if err != nil {
		return "", "", fmt.Errorf("讀不到 go.mod：%w", err)
	}
	m := requireRE.FindSubmatch(b)
	if m == nil {
		return "", "", fmt.Errorf("go.mod 沒有 require %s：harness 只管用 terminu-go 做的 app", TGO)
	}
	tgo := string(m[1])
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", TGO)
	cmd.Dir = repo
	out, err := cmd.Output()
	if err != nil {
		return tgo, "", fmt.Errorf("找不到 terminu-go %s 的原始碼（go list -m）：%w", tgo, err)
	}
	cl, err := os.ReadFile(filepath.Join(strings.TrimSpace(string(out)), "CHANGELOG.md"))
	if err != nil {
		return tgo, "", fmt.Errorf("讀不到 terminu-go 的 CHANGELOG：%w", err)
	}
	for _, m := range changelogRE.FindAllSubmatch(cl, -1) {
		if string(m[1]) == tgo {
			return tgo, string(m[2]), nil
		}
	}
	return tgo, "", fmt.Errorf("terminu-go 的 CHANGELOG 沒有寫 %s 對到哪一版 tdp", tgo)
}

var testRE = regexp.MustCompile(`^(?:Test|Example|Fuzz|Benchmark)\w*$`)

// Tests runs go test -list in the repo and returns the names it lists.
func Tests(repo string) (map[string]bool, error) {
	cmd := exec.Command("go", "test", "-list", ".", "./...")
	cmd.Dir = repo
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%w\n%s", err, out)
	}
	names := map[string]bool{}
	for _, line := range strings.Split(string(out), "\n") {
		if line = strings.TrimSpace(line); testRE.MatchString(line) {
			names[line] = true
		}
	}
	return names, nil
}

// Commit reports whether the repo has the commit.
func Commit(repo, hash string) bool {
	return exec.Command("git", "-C", repo, "cat-file", "-e", hash+"^{commit}").Run() == nil
}

// Text is the report as the harness prints it.
func (r Report) Text() string {
	var b strings.Builder
	if r.TGO != "" {
		fmt.Fprintf(&b, "terminu-go %s，tdp %s\n", r.TGO, r.TDP)
	}
	if len(r.Problems) == 0 {
		b.WriteString("問題：沒有\n")
	} else {
		fmt.Fprintf(&b, "問題（%d）：\n", len(r.Problems))
		for _, p := range r.Problems {
			fmt.Fprintf(&b, "- %s\n", p)
		}
	}
	if len(r.Todo) > 0 {
		fmt.Fprintf(&b, "待做（%d）：\n", len(r.Todo))
		for _, t := range r.Todo {
			fmt.Fprintf(&b, "- %s\n", t)
		}
	}
	return b.String()
}

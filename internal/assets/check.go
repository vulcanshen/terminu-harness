package assets

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	repocheck "github.com/vulcanshen/terminu-harness/internal/repo"
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
	// Subtests runs the subtests named (TestX/sub) and returns the ones that
	// ran; nil skips checking them.
	Subtests func(repo string, names []string) (map[string]bool, error)
	// Commit reports whether the repo has the commit.
	Commit func(repo, hash string) bool
}

// Default reads go.mod and terminu-go's CHANGELOG for the versions, tdp from
// src, and runs go test -list, go test -run and git in the repo.
func Default(src tdp.Source) Env {
	return Env{Versions: Versions, IDs: src.IDs, Tests: Tests, Subtests: Subtests, Commit: Commit}
}

// Problem is something in the files to fix.
type Problem struct {
	File string // a file in Dir, go.mod, or a file of the repo
	Task int    // 0 when it is not about one task
	Msg  string
}

func (p Problem) String() string {
	if p.Task > 0 {
		return fmt.Sprintf("%s task %d：%s", p.File, p.Task, p.Msg)
	}
	return fmt.Sprintf("%s：%s", p.File, p.Msg)
}

// Report is what Check found: problems to fix, what it could not check this
// time (tdp not read, the tests not run), and what is left to do.
type Report struct {
	TGO, TDP   string
	Incomplete []string
	Problems   []Problem
	Todo       []string
}

// Failed reports whether the report is not a pass: a problem, or something
// left unchecked.
func (r Report) Failed() bool {
	return len(r.Problems) > 0 || len(r.Incomplete) > 0
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
			r.Incomplete = append(r.Incomplete, fmt.Sprintf("讀不到 tdp %s（%v）：spec.yaml 的 tdp 欄位、家族成員的 repo 都沒有檢查，這次不算通過；網路好了再跑一次 check", tdpVersion, err))
		} else {
			ids = &got
		}
	}
	// The repo against harness's project/, once the tdp the app follows has
	// none (v1.2.0-beta.1 on); before, that terminu-go's tdptest checks it.
	repoChecked := false
	if ids != nil && !ids.Project() {
		repoChecked = repocheck.IsFamily(repo)
		for _, p := range repocheck.Check(repo) {
			add(p.File, 0, "%s", p.Msg)
		}
	}

	checkSpec(f.Spec, ids, add)
	spec := map[int]SpecTask{}
	for _, t := range f.Spec.Tasks {
		spec[t.ID] = t
	}

	work := map[int]WorkTask{}
	if f.Work != nil {
		wc := workCheck{repo: repo, spec: spec, commit: env.Commit, ids: ids, repoChecked: repoChecked}
		if env.Tests != nil {
			if wc.tests, err = env.Tests(repo); err != nil {
				r.Incomplete = append(r.Incomplete, fmt.Sprintf("go test -list 跑不起來（%v）：work.yaml 的測試名字沒有檢查", err))
				wc.tests = nil
			}
		}
		if subs := subtests(f.Work); len(subs) > 0 && env.Subtests != nil {
			if wc.ran, err = env.Subtests(repo, subs); err != nil {
				r.Incomplete = append(r.Incomplete, fmt.Sprintf("子測試跑不起來（%v）：work.yaml 的子測試沒有檢查", err))
				wc.ran = nil
			}
		}
		for _, w := range f.Work.Tasks {
			if _, dup := work[w.ID]; dup {
				add(WorkFile, w.ID, "id 重複")
				continue
			}
			work[w.ID] = w
			wc.task(w, add)
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
			inHarness, id := tdp.Cited(c)
			switch {
			case !inHarness && !strings.HasPrefix(c, "tdp ") || id == "":
				add(SpecFile, t.ID, "tdp 寫成 %q：一律用 tdp 加編號，例 \"tdp K4\"、\"tdp panel.5\"、\"tdp confirm\"、\"tdp K9-6\"；家族的 repo 範本用 harness 加編號，例 \"harness project.3\"、\"harness project/install.sh\"、\"harness project-13\"", c)
			case inHarness && !tdp.HarnessIDs().Has(id):
				add(SpecFile, t.ID, "%q 在 harness 的 project 裡找不到", c)
			case !inHarness && ids != nil && !ids.Has(id) && tdp.HarnessIDs().Has(id):
				add(SpecFile, t.ID, "%q 在這一版 tdp 裡找不到：家族的 repo 範本從 tdp v1.2.0-beta.1 起在 harness，寫成 \"harness %s\"", c, id)
			case !inHarness && ids != nil && !ids.Has(id):
				add(SpecFile, t.ID, "%q 在這一版 tdp 裡找不到", c)
			}
		}
	}
}

// workCheck is what checking work.yaml's tasks needs: the spec, the app's
// tests and the subtests that ran, how to look up a commit, the tdp the app
// follows, and whether check read the repo against harness's project/.
type workCheck struct {
	repo        string
	spec        map[int]SpecTask
	tests       map[string]bool // nil: not checked
	ran         map[string]bool // subtests that ran; nil: not checked
	commit      func(repo, hash string) bool
	ids         *tdp.IDs // nil: tdp not read
	repoChecked bool
}

// subtests are the subtests work.yaml's live tasks name (TestX/sub), each once.
func subtests(work *Work) []string {
	var out []string
	seen := map[string]bool{}
	for _, w := range work.Tasks {
		if w.Canceled {
			continue
		}
		for _, b := range w.Test {
			if strings.Contains(b.Test, "/") && !seen[b.Test] {
				seen[b.Test] = true
				out = append(out, b.Test)
			}
		}
	}
	return out
}

func (wc workCheck) task(w WorkTask, add addFunc) {
	repo, spec, tests, commit := wc.repo, wc.spec, wc.tests, wc.commit
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
		for _, set := range []bool{b.Test != "", b.Check != "", b.Commit != "", b.User != nil} {
			if set {
				n++
			}
		}
		if n != 1 {
			add(WorkFile, w.ID, "test 的每一筆只寫一種：測試的名字、check、commit、或 user")
			continue
		}
		switch {
		case b.Test != "":
			name, _, isSub := strings.Cut(b.Test, "/")
			switch {
			case tests != nil && !tests[name]:
				add(WorkFile, w.ID, "測試 %s 不在 go test -list 裡", b.Test)
			case isSub && wc.ran != nil && !wc.ran[b.Test] && strings.HasSuffix(b.Test, "/repo") && wc.ids != nil && !wc.ids.Project():
				add(WorkFile, w.ID, "子測試 %s 跑不出來：tdptest 的 repo 子測試從 terminu-go v1.2.0-beta.1 起拿掉了，repo 改由 harness 的 check 查，作保寫成 check: repo", b.Test)
			case isSub && wc.ran != nil && !wc.ran[b.Test]:
				add(WorkFile, w.ID, "子測試 %s 跑不出來：go test -run 沒有跑到這一個", b.Test)
			}
		case b.Check != "":
			switch {
			case b.Check != "repo":
				add(WorkFile, w.ID, "check 寫成 %q：只有 check: repo（harness 的 check 拿 harness 的範本比對家族成員的 repo）", b.Check)
			case wc.ids == nil:
				// tdp not read: whether harness checks the repo is not known.
			case wc.ids.Project():
				add(WorkFile, w.ID, "check: repo 守不住：這一版 tdp 還有 project/，repo 由 terminu-go 的 tdptest 查，harness 不查；作保寫那個測試")
			case !wc.repoChecked:
				add(WorkFile, w.ID, "check: repo 守不住：go.mod 的 module 不在 github.com/vulcanshen 底下，harness 不查這個 repo")
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

var runRE = regexp.MustCompile(`(?m)^\s*=== RUN\s+(\S+)\s*$`)

// Subtests runs the subtests named (TestX/sub, TestX/sub/deeper) in the repo
// and returns the names go test -v reported running: one go test -run per
// top-level test, each level anchored, the names at a level ORed. Whether
// they pass is make check's business; a build that fails is an error.
func Subtests(repo string, names []string) (map[string]bool, error) {
	var tops []string
	paths := map[string][][]string{}
	for _, n := range names {
		parts := strings.Split(n, "/")
		if _, ok := paths[parts[0]]; !ok {
			tops = append(tops, parts[0])
		}
		paths[parts[0]] = append(paths[parts[0]], parts[1:])
	}
	ran := map[string]bool{}
	for _, top := range tops {
		levels := []string{"^" + regexp.QuoteMeta(top) + "$"}
		for depth := 0; ; depth++ {
			var alts []string
			seen := map[string]bool{}
			for _, p := range paths[top] {
				if depth < len(p) && !seen[p[depth]] {
					seen[p[depth]] = true
					alts = append(alts, regexp.QuoteMeta(p[depth]))
				}
			}
			if len(alts) == 0 {
				break
			}
			levels = append(levels, "^(?:"+strings.Join(alts, "|")+")$")
		}
		run := strings.Join(levels, "/")
		cmd := exec.Command("go", "test", "-run", run, "-v", "-count=1", "./...")
		cmd.Dir = repo
		out, err := cmd.CombinedOutput()
		var exit *exec.ExitError
		if err != nil && !errors.As(err, &exit) {
			return nil, fmt.Errorf("go test -run %s: %w", run, err)
		}
		if strings.Contains(string(out), "[build failed]") || strings.Contains(string(out), "[setup failed]") {
			return nil, fmt.Errorf("go test -run %s 編譯不過\n%s", run, out)
		}
		for _, m := range runRE.FindAllStringSubmatch(string(out), -1) {
			ran[m[1]] = true
		}
	}
	return ran, nil
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
	if len(r.Incomplete) > 0 {
		fmt.Fprintf(&b, "沒有檢查完（%d），這次不算通過：\n", len(r.Incomplete))
		for _, s := range r.Incomplete {
			fmt.Fprintf(&b, "- %s\n", s)
		}
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

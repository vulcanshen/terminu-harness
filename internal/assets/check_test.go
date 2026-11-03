package assets

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vulcanshen/terminu-harness/internal/tdp"
)

func fakeEnv() Env {
	return Env{
		Versions: func(string) (string, string, error) { return "v0.5.1", "v0.10.1", nil },
		IDs: func(string) (tdp.IDs, error) {
			return tdp.IDs{
				Sections:  map[string]string{"K9": "", "P5": "", "panel.1": "", "panel.5": "", "form.4": ""},
				Files:     map[string]string{"confirm": "", "form": ""},
				Templates: map[string]bool{"project/install.sh": true},
				Items:     map[string][]string{"K9-6": {"K9"}},
			}, nil
		},
		Tests:  func(string) (map[string]bool, error) { return map[string]bool{"TestList": true, "TestQuit": true}, nil },
		Commit: func(_, hash string) bool { return hash == "abc1234" },
	}
}

func write(t *testing.T, files map[string]string) string {
	t.Helper()
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, Dir), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(repo, Dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return repo
}

const goodSpec = `id: 2
from: 0
to: v0.10.1
tasks:
  - id: 1
    spec: 1
    request: 書是一個清單
    tdp: [tdp panel.1, tdp panel.5]
  - id: 2
    spec: 2
    request: 按 q 離開，先問
    tdp: [tdp K9, tdp confirm, tdp K9-6, tdp project/install.sh]
    remark:
      - 離開前一律問，因為書只在記憶體裡
  - id: 3
    spec: 1
    request: 匯出全部的書
    tdp: []
    remark:
      - note: 進度要用哪個元件顯示，tdp 沒寫
        issue: tdp #12
  - id: 4
    spec: 2
    request: 書可以加星號
    tdp: [tdp P5]
    canceled: 使用者說不要了
`

const goodWork = `tasks:
  - id: 1
    spec: 1
    implementation: 清單 panel，每本書一列
    test: [TestList]
  - id: 2
    spec: 1
    implementation: q 打開 confirm「Quit demou?」
    test:
      - TestQuit/asks
      - commit: abc1234
      - user: {date: 2026-10-26, how: 按 q 看有沒有問}
  - id: 4
    spec: 2
    implementation: 星號
    test: [TestList]
    canceled: true
`

const goodReview = `spec: 1
tasks:
  - id: 1
    accepted: true
  - id: 2
    accepted: false
    reason: 違反：confirm 的問句不在最後
`

func TestCheckGood(t *testing.T) {
	repo := write(t, map[string]string{IdeaFile: "# 書架\n", SpecFile: goodSpec, WorkFile: goodWork, ReviewFile: goodReview})
	r := Check(repo, fakeEnv())
	if len(r.Problems) > 0 {
		t.Fatalf("problems in a good repo:\n%s", r.Text())
	}
	want := []string{
		"task 2：需求在 spec 2 改過，work 還是照 spec 1 做的",
		"task 3：還沒做（spec 1）",
		"上次 review 對的是 spec 1，現在是 spec 2",
	}
	if strings.Join(r.Todo, "\n") != strings.Join(want, "\n") {
		t.Errorf("todo:\n%s\nwant:\n%s", strings.Join(r.Todo, "\n"), strings.Join(want, "\n"))
	}
	if r.TGO != "v0.5.1" || r.TDP != "v0.10.1" {
		t.Errorf("versions %s %s", r.TGO, r.TDP)
	}
}

func TestCheckProblems(t *testing.T) {
	spec := `id: 1
from: 0
to: v0.10.0
tasks:
  - id: 1
    spec: 1
    request: 清單
    tdp: [panel.5, tdp panel.99, tdp K9-99, tdp project/README-zh_TW.md]
  - id: 2
    spec: 3
    request: 沒有根據
  - id: 2
    spec: 1
    request: 重複的 id
    tdp: [tdp K9]
  - id: 3
    spec: 1
    request: 取消了
    tdp: [tdp K9]
    canceled: 不要了
`
	work := `tasks:
  - id: 1
    spec: 2
    implementation: 清單
    test:
      - TestGone
      - commit: deadbeef
      - user: {date: 10/26, how: ""}
      - {test: TestList, commit: abc1234}
  - id: 3
    spec: 1
    implementation: 還在
    test: [TestList]
  - id: 9
    spec: 1
    implementation: spec 沒有
    test: [TestList]
`
	review := `spec: 1
tasks:
  - id: 1
    accepted: false
    reason: confirm 的問句不在最後
`
	repo := write(t, map[string]string{SpecFile: spec, WorkFile: work, ReviewFile: review})
	r := Check(repo, fakeEnv())
	got := r.Text()
	for _, want := range []string{
		"idea.md：沒有這個檔",
		"spec.yaml：to 是 \"v0.10.0\"，但 go.mod 的 terminu-go v0.5.1 對到 tdp v0.10.1",
		"spec.yaml task 1：tdp 寫成 \"panel.5\"",
		"spec.yaml task 1：\"tdp panel.99\" 在這一版 tdp 裡找不到",
		"spec.yaml task 1：\"tdp K9-99\" 在這一版 tdp 裡找不到",
		"spec.yaml task 1：\"tdp project/README-zh_TW.md\" 在這一版 tdp 裡找不到",
		"spec.yaml task 2：spec 是 3",
		"spec.yaml task 2：沒有 tdp 根據",
		"spec.yaml task 2：id 重複",
		"work.yaml task 1：spec 是 2",
		"work.yaml task 1：測試 TestGone 不在 go test -list 裡",
		"work.yaml task 1：commit deadbeef 不在這個 repo 裡",
		"work.yaml task 1：使用者確認的 date 是 \"10/26\"",
		"work.yaml task 1：使用者確認沒有寫 how",
		"work.yaml task 1：test 的每一筆只寫一種",
		"work.yaml task 3：spec.yaml 取消了這個 task，這裡也要標 canceled: true",
		"work.yaml task 9：spec.yaml 沒有這個 task",
		"review.yaml task 1：reason 開頭寫是哪一種",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
}

func TestCheckMissing(t *testing.T) {
	r := Check(write(t, nil), fakeEnv())
	got := r.Text()
	for _, want := range []string{"idea.md：沒有這個檔", "spec.yaml：沒有這個檔"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	r = Check(write(t, map[string]string{SpecFile: "id: [\n"}), fakeEnv())
	if len(r.Problems) != 1 || !strings.Contains(r.Problems[0].Msg, "spec.yaml") {
		t.Errorf("a broken spec.yaml: %s", r.Text())
	}
}

func TestVersionPatterns(t *testing.T) {
	for _, gomod := range []string{
		"module x\n\nrequire github.com/vulcanshen/terminu-go v0.5.1\n",
		"module x\n\nrequire (\n\tgithub.com/creack/pty v1.1.24\n\tgithub.com/vulcanshen/terminu-go v0.5.1 // indirect\n)\n",
	} {
		m := requireRE.FindStringSubmatch(gomod)
		if m == nil || m[1] != "v0.5.1" {
			t.Errorf("require in %q: %v", gomod, m)
		}
	}
	if requireRE.MatchString("require github.com/vulcanshen/terminu-gox v1.0.0\n") {
		t.Error("another module matched")
	}
	cl := "# Changelog\n\n## v0.5.1 — tdp v0.10.1\n\n- …\n\n## v0.5.0 — tdp v0.10.0\n"
	ms := changelogRE.FindAllStringSubmatch(cl, -1)
	if len(ms) != 2 || ms[0][1] != "v0.5.1" || ms[0][2] != "v0.10.1" {
		t.Errorf("changelog: %v", ms)
	}
	m := changelogRE.FindStringSubmatch("## v1.1.0-beta.1 — tdp v1.1.0-beta.1\n")
	if m == nil || m[1] != "v1.1.0-beta.1" || m[2] != "v1.1.0-beta.1" {
		t.Errorf("a pre-release in the changelog: %v", m)
	}
	for v, want := range map[string]bool{
		"v1.0.2": true, "v1.1.0-beta.1": true, "v1.1.0-rc.2": true,
		"v1.1": false, "1.0.2": false, "v1.1.0-": false, "v1.1.0-beta..1": false,
	} {
		if isVersion(v) != want {
			t.Errorf("isVersion(%q) = %v", v, !want)
		}
	}
}

// TestCheckHarnessCitations: the family's repo templates are cited with the
// harness prefix (harness project.3, harness project/install.sh, harness
// project-13), checked against harness's project/; a tdp project ID the tdp
// no longer has says to write it with harness.
func TestCheckHarnessCitations(t *testing.T) {
	spec := `id: 1
from: 0
to: v0.10.1
tasks:
  - id: 1
    spec: 1
    request: repo 照家族的範本
    tdp: [harness project.3, harness project/install.sh, harness project-13, harness project]
  - id: 2
    spec: 1
    request: 範本
    tdp: [harness project.99, tdp project.3, harness K9]
`
	r := Check(write(t, map[string]string{IdeaFile: "# 書架\n", SpecFile: spec, WorkFile: "tasks: []\n"}), fakeEnv())
	got := r.Text()
	if strings.Contains(got, "spec.yaml task 1：") {
		t.Errorf("harness citations taken as problems:\n%s", got)
	}
	for _, want := range []string{
		"spec.yaml task 2：\"harness project.99\" 在 harness 的 project 裡找不到",
		"spec.yaml task 2：\"tdp project.3\" 在這一版 tdp 裡找不到：家族的 repo 範本從 tdp v1.2.0-beta.1 起在 harness，寫成 \"harness project.3\"",
		"spec.yaml task 2：\"harness K9\" 在 harness 的 project 裡找不到",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
}

// TestCheckRepo: once the tdp the app follows has no project/, check reads
// the family member's repo against harness's project/; while it has, it
// leaves the repo to that terminu-go's tdptest.
func TestCheckRepo(t *testing.T) {
	dir := write(t, map[string]string{IdeaFile: "# 書架\n", SpecFile: goodSpec, WorkFile: goodWork, ReviewFile: goodReview})
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module github.com/vulcanshen/shelfu\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	r := Check(dir, fakeEnv())
	for _, want := range []string{"Makefile：沒有這個檔（harness project.5", "family.mk：沒有這個檔（harness project.2、project-1）", "LICENSE：沒有這個檔"} {
		if !strings.Contains(r.Text(), want) {
			t.Errorf("missing %q in:\n%s", want, r.Text())
		}
	}
	env := fakeEnv()
	env.IDs = func(v string) (tdp.IDs, error) {
		ids, err := fakeEnv().IDs(v)
		ids.Files["project"] = "project/README-zh_TW.md"
		return ids, err
	}
	if r := Check(dir, env); len(r.Problems) > 0 {
		t.Errorf("tdp with project/, the repo checked:\n%s", r.Text())
	}
}

const backingSpec = `id: 1
from: 0
to: v0.10.1
tasks:
  - {id: 1, spec: 1, request: 子測試, tdp: [tdp K9]}
  - {id: 2, spec: 1, request: 不見的子測試, tdp: [tdp K9]}
  - {id: 3, spec: 1, request: tdptest 的 repo, tdp: [tdp K9]}
  - {id: 4, spec: 1, request: repo, tdp: [harness project.3]}
  - {id: 5, spec: 1, request: 別的 check, tdp: [tdp K9]}
`

const backingWork = `tasks:
  - {id: 1, spec: 1, implementation: a, test: [TestList/sub]}
  - {id: 2, spec: 1, implementation: b, test: [TestList/gone]}
  - {id: 3, spec: 1, implementation: c, test: [TestList/repo]}
  - {id: 4, spec: 1, implementation: d, test: [{check: repo}]}
  - {id: 5, spec: 1, implementation: e, test: [{check: lint}, {check: repo, commit: abc1234}]}
`

// TestCheckBackings: a subtest backs a task only when go test runs it, and a
// tdptest repo subtest gone with tdp's project/ says to write check: repo;
// check: repo backs a task only where harness's check reads the repo: a
// family member's, its tdp without project/.
func TestCheckBackings(t *testing.T) {
	dir := write(t, map[string]string{IdeaFile: "# 書架\n", SpecFile: backingSpec, WorkFile: backingWork})
	gomod := func(module string) {
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module "+module+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gomod("github.com/vulcanshen/shelfu")
	env := fakeEnv()
	var asked []string
	env.Subtests = func(_ string, names []string) (map[string]bool, error) {
		asked = names
		return map[string]bool{"TestList/sub": true}, nil
	}
	got := Check(dir, env).Text()
	if strings.Join(asked, " ") != "TestList/sub TestList/gone TestList/repo" {
		t.Errorf("subtests asked: %q", asked)
	}
	for _, want := range []string{
		"work.yaml task 2：子測試 TestList/gone 跑不出來：go test -run 沒有跑到這一個",
		"work.yaml task 3：子測試 TestList/repo 跑不出來：tdptest 的 repo 子測試從 terminu-go v1.2.0-beta.1 起拿掉了，repo 改由 harness 的 check 查，作保寫成 check: repo",
		"work.yaml task 5：check 寫成 \"lint\"：只有 check: repo",
		"work.yaml task 5：test 的每一筆只寫一種：測試的名字、check、commit、或 user",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	for _, not := range []string{"work.yaml task 1：", "work.yaml task 4："} {
		if strings.Contains(got, not) {
			t.Errorf("%q in:\n%s", not, got)
		}
	}

	withProject := env
	withProject.IDs = func(v string) (tdp.IDs, error) {
		ids, err := fakeEnv().IDs(v)
		ids.Files["project"] = "project/README-zh_TW.md"
		return ids, err
	}
	got = Check(dir, withProject).Text()
	for _, want := range []string{
		"work.yaml task 3：子測試 TestList/repo 跑不出來：go test -run 沒有跑到這一個",
		"work.yaml task 4：check: repo 守不住：這一版 tdp 還有 project/",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("tdp with project/: missing %q in:\n%s", want, got)
		}
	}

	gomod("github.com/someone/shelfu")
	if got := Check(dir, env).Text(); !strings.Contains(got, "work.yaml task 4：check: repo 守不住：go.mod 的 module 不在 github.com/vulcanshen 底下") {
		t.Errorf("not a family member's: %s", got)
	}
}

// TestCheckIncomplete: tdp that cannot be read, tests that do not run, are
// not a pass: they head the report as what was not checked, and the repo is
// not taken as checked.
func TestCheckIncomplete(t *testing.T) {
	dir := write(t, map[string]string{IdeaFile: "# 書架\n", SpecFile: backingSpec, WorkFile: backingWork})
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module github.com/vulcanshen/shelfu\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	env := fakeEnv()
	env.IDs = func(string) (tdp.IDs, error) { return tdp.IDs{}, errors.New("api.github.com: timeout") }
	env.Tests = func(string) (map[string]bool, error) { return nil, errors.New("exit 1") }
	env.Subtests = func(string, []string) (map[string]bool, error) { return nil, errors.New("build failed") }
	r := Check(dir, env)
	if !r.Failed() || len(r.Incomplete) != 3 {
		t.Fatalf("incomplete %q", r.Incomplete)
	}
	got := r.Text()
	for _, want := range []string{
		"沒有檢查完（3），這次不算通過：",
		"- 讀不到 tdp v0.10.1（api.github.com: timeout）：spec.yaml 的 tdp 欄位、家族成員的 repo 都沒有檢查，這次不算通過；網路好了再跑一次 check",
		"- go test -list 跑不起來（exit 1）：work.yaml 的測試名字沒有檢查",
		"- 子測試跑不起來（build failed）：work.yaml 的子測試沒有檢查",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	if strings.Index(got, "沒有檢查完") > strings.Index(got, "問題") {
		t.Errorf("not checked comes after the problems:\n%s", got)
	}
	for _, not := range []string{"Makefile：", "check: repo 守不住", "跑不出來"} {
		if strings.Contains(got, not) {
			t.Errorf("%q with nothing read:\n%s", not, got)
		}
	}
}

// TestSubtests runs go test in a module of its own: the subtests named come
// back when go test runs them, passing or failing, deeper ones too; a name
// nothing runs does not; a build that fails is an error.
func TestSubtests(t *testing.T) {
	if testing.Short() {
		t.Skip("runs go test")
	}
	dir := t.TempDir()
	put := func(name, s string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	put("go.mod", "module example.com/shelf\n\ngo 1.22\n")
	put("shelf_test.go", `package shelf

import "testing"

func TestA(t *testing.T) {
	t.Run("sub", func(t *testing.T) {})
	t.Run("other", func(t *testing.T) { t.Run("deep", func(t *testing.T) {}) })
}

func TestB(t *testing.T) { t.Run("fails", func(t *testing.T) { t.Fatal("no") }) }
`)
	ran, err := Subtests(dir, []string{"TestA/sub", "TestA/other/deep", "TestA/gone", "TestB/fails", "TestC/x"})
	if err != nil {
		t.Fatal(err)
	}
	for name, want := range map[string]bool{"TestA/sub": true, "TestA/other/deep": true, "TestB/fails": true, "TestA/gone": false, "TestC/x": false} {
		if ran[name] != want {
			t.Errorf("%s ran: %v", name, ran[name])
		}
	}
	put("broken.go", "package shelf\n\nfunc {\n")
	if _, err := Subtests(dir, []string{"TestA/sub"}); err == nil || !strings.Contains(err.Error(), "編譯不過") {
		t.Errorf("a build that fails: %v", err)
	}
}

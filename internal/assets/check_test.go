package assets

import (
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
}

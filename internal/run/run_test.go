package run

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const items = "# items\n\n## rules — K Core key\n\n| 編號 | 出處 | 檢查什麼 | 適用 | 怎麼驗 |\n|---|---|---|---|---|\n" +
	"| K1-1 | K1 | a | 都要 | b |\n| K1-2 | K1 | c | 表單 | d |\n| K1-3 | K1 | e | 都要 | f |\n" +
	"\n## project\n\n| 編號 | 出處 | 檢查什麼 | 適用 | 怎麼驗 |\n|---|---|---|---|---|\n| project-1 | project §範本檔 | g | 都要 | 讀檔案 |\n"

func TestRun(t *testing.T) {
	base := t.TempDir()
	r, err := Start(base, "./demou --flag", "", "v0.7.4", []byte(items), 2)
	if err != nil {
		t.Fatal(err)
	}
	if r.Meta.App != "demou" || !strings.Contains(r.Dir, filepath.Join(base, "demou")) || !strings.HasSuffix(r.Dir, "-tdp-0.7.4") {
		t.Errorf("dir %s, app %s", r.Dir, r.Meta.App)
	}
	if len(r.Meta.Batches) != 2 || r.Total() != 4 {
		t.Fatalf("batches %+v, total %d", r.Meta.Batches, r.Total())
	}
	id, path := r.NewEvidence(".txt")
	if id != "e0001" {
		t.Errorf("evidence %s", id)
	}
	os.WriteFile(path, []byte("screen"), 0o644)

	kept, problems := r.AddVerdicts([]Verdict{
		{ID: "K1-1", Verdict: Meets, Reason: "ok", Evidence: []string{id}},
		{ID: "K1-2", Verdict: NotApply, Reason: "沒有表單"},
		{ID: "K1-3", Verdict: "pass", Reason: "x"},
		{ID: "Z9-9", Verdict: Meets, Reason: "x"},
		{ID: "project-1", Verdict: Violates, Reason: ""},
	})
	if kept != 2 || len(problems) != 3 {
		t.Errorf("kept %d, problems %v", kept, problems)
	}
	// A later verdict replaces the earlier one.
	r.AddVerdicts([]Verdict{{ID: "K1-1", Verdict: Violates, Reason: "按 `j|k` 沒動"}})
	if err := r.AddGap(Gap{What: "一個 tdp 沒有的 widget", Where: "第二個畫面"}); err != nil {
		t.Fatal(err)
	}

	// The run reads the same from its directory.
	r, err = Open(r.Dir)
	if err != nil {
		t.Fatal(err)
	}
	if id, _ := r.NewEvidence(".png"); id != "e0002" {
		t.Errorf("evidence numbering after reopening: %s", id)
	}
	missing, _ := r.Missing()
	if len(missing["b1"])+len(missing["b2"]) != 2 {
		t.Errorf("missing %v", missing)
	}
	sum, err := r.Finish()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sum, "tdp-0.7.4 0 gene") || !strings.Contains(sum, "沒有 judge 判到") {
		t.Errorf("summary %q", sum)
	}
	report, _ := os.ReadFile(filepath.Join(r.Dir, ReportFile))
	for _, want := range []string{"**不通過**　tdp-0.7.4 0 gene", "違反（1）：在等 tdp 或 terminu-go 的 0 條", "app 自己要修的 1 條：K1-1", "| K1-1 | 違反 | 按 `j\\|k` 沒動 |", "一個 tdp 沒有的 widget", "| 合計 | 4 |", "判不了（含沒判到的 2 條） | 2"} {
		if !strings.Contains(string(report), want) {
			t.Errorf("report misses %q", want)
		}
	}
}

func TestOutcome(t *testing.T) {
	repo := t.TempDir()
	r, err := Start(t.TempDir(), "demo", repo, "v0.7.4", []byte(items), 1)
	if err != nil {
		t.Fatal(err)
	}
	if got, _ := r.Outcome(nil); got != Passed {
		t.Errorf("no violation: %s", got)
	}
	if got, unlisted := r.Outcome([]string{"K1-1"}); got != Failed || len(unlisted) != 1 {
		t.Errorf("no fix file: %s %v", got, unlisted)
	}
	os.MkdirAll(filepath.Join(repo, "docs"), 0o755)
	os.WriteFile(filepath.Join(repo, "docs", "demo-terminu-fix.md"), []byte("## tdp K1-1 a thing\n## tdp K1-30 another\n"), 0o644)
	if got, unlisted := r.Outcome([]string{"K1-1"}); got != Failed || unlisted != nil {
		t.Errorf("listed, still failed: %s %v", got, unlisted)
	}
	// K1-3 is not K1-30.
	if got, unlisted := r.Outcome([]string{"K1-1", "K1-3"}); got != Failed || len(unlisted) != 1 || unlisted[0] != "K1-3" {
		t.Errorf("one unlisted: %s %v", got, unlisted)
	}
}

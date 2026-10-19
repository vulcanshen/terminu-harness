//go:build unix

package server

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vulcanshen/terminu-harness/internal/run"
)

const carryItems = "# items\n\n## rules — K Core key\n\n| 編號 | 出處 | 檢查什麼 | 適用 | 怎麼驗 |\n|---|---|---|---|---|\n" +
	"| K1-1 | K1 | a | 都要 | 按 `Space` |\n| E7-1 | E7 | b | 都要 | 讀檔案：README |\n"

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir, "-c", "user.email=t@t", "-c", "user.name=t"}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// judged runs a judging of the repo that keeps every item as met.
func judged(t *testing.T, base, repo string) *run.Run {
	t.Helper()
	r, err := run.Start(base, "demo", repo, "v0.7.4", []byte(carryItems), 1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := carryOver(r, base, repo, "v0.7.4"); err != nil {
		t.Fatal(err)
	}
	return r
}

func TestCarryOver(t *testing.T) {
	base, repo := t.TempDir(), t.TempDir()
	os.WriteFile(filepath.Join(repo, "go.mod"), []byte("module demo\n\nrequire github.com/vulcanshen/terminu-go v0.1.13\n"), 0o644)
	os.WriteFile(filepath.Join(repo, "README.md"), []byte("demo\n"), 0o644)
	os.WriteFile(filepath.Join(repo, "main.go"), []byte("package main\n"), 0o644)
	git(t, repo, "init", "-q")
	git(t, repo, "add", "-A")
	git(t, repo, "commit", "-q", "-m", "one")

	first := judged(t, base, repo)
	if first.Meta.Commit == "" || first.Meta.Tgo != "v0.1.13" || first.Meta.Carried != 0 {
		t.Fatalf("first run meta %+v", first.Meta)
	}
	first.AddVerdicts([]run.Verdict{{ID: "K1-1", Verdict: run.Meets, Reason: "menu opens"}, {ID: "E7-1", Verdict: run.Meets, Reason: "README ok"}})
	if _, err := first.Finish(); err != nil {
		t.Fatal(err)
	}

	// Only the README changed: the screen item is carried, the file item judged again.
	os.WriteFile(filepath.Join(repo, "README.md"), []byte("demo, changed\n"), 0o644)
	git(t, repo, "commit", "-q", "-am", "two")
	second := judged(t, base, repo)
	got, _ := second.Judged()
	if v, ok := got["K1-1"]; !ok || !strings.HasPrefix(v.Reason, "沿用 ") {
		t.Errorf("K1-1 not carried: %+v", got)
	}
	if _, ok := got["E7-1"]; ok {
		t.Error("E7-1 carried though the README changed")
	}
	second.AddVerdicts([]run.Verdict{{ID: "E7-1", Verdict: run.Meets, Reason: "README still ok"}})
	second.Finish()

	// Code changed, not yet committed: the file item is carried, the screen item judged again.
	os.WriteFile(filepath.Join(repo, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644)
	third := judged(t, base, repo)
	got, _ = third.Judged()
	if _, ok := got["K1-1"]; ok {
		t.Error("K1-1 carried though the code changed")
	}
	if _, ok := got["E7-1"]; !ok {
		t.Error("E7-1 not carried though only code changed")
	}
}

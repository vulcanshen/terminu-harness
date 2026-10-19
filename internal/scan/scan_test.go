package scan

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const appSource = `package ui

import terminu "github.com/vulcanshen/terminu-go"

func build(a *terminu.App) {
	a.AddScreen(&terminu.Screen{})
	l := &terminu.List{}
	p := &terminu.Panel{Body: l, Details: nil}
	_ = p
	a.Text(terminu.TextInput{Title: "Name"})
	a.Confirm(terminu.Ask{}, nil)
	a.Toast(terminu.Info, "done")
}
`

func repo(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, body := range files {
		os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0o755)
		os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644)
	}
	return dir
}

func TestScan(t *testing.T) {
	r, err := Scan(repo(t, map[string]string{
		"internal/ui/ui.go":      appSource,
		"internal/ui/ui_test.go": "package ui\n// terminu.DatePicker in a test does not count\n",
	}), "v0.1.13")
	if err != nil {
		t.Fatal(err)
	}
	for f, want := range map[string]bool{
		"明細 panel": true, "text": true, "confirm": true, "toast": true, "input": true,
		"多畫面": false, "datetime picker": false, "finder": false, "表單": false, "viewer": false,
		"模式": false, "模式（在 popup 裡）": false, "有預覽的 finder": false,
	} {
		if r.Has[f] != want {
			t.Errorf("%s: %v, want %v", f, r.Has[f], want)
		}
	}
	if r.Evidence["text"] != filepath.Join("internal", "ui", "ui.go") {
		t.Errorf("evidence %q", r.Evidence["text"])
	}
	for _, c := range []struct {
		applies string
		na      bool
	}{
		{"都要", false},
		{"datetime picker", true},
		{"finder、file-picker", true},
		{"表單、input", false}, // it has an input
		{"多 panel", false},  // the scan can't tell
		{"text＋表單", true},   // no form
		{"confirm＋datetime picker", true},
		{"模式（在 panel 上）", true},
		{"發過版", false}, // not a git repo: the scan can't tell
	} {
		if got, _ := r.NotApplicable(c.applies); got != c.na {
			t.Errorf("%q: %v, want %v", c.applies, got, c.na)
		}
	}
	if none, _ := Scan(t.TempDir(), ""); none != nil {
		t.Error("an app without terminu-go was scanned")
	}
}

func TestCustomCodeMakesItUnsure(t *testing.T) {
	r, _ := Scan(repo(t, map[string]string{
		"internal/ui/ui.go":   appSource,
		"internal/ui/grid.go": "package ui\ntype grid struct{}\nfunc (g *grid) Draw(w, h int, focused bool) []string { return nil }\n",
	}), "v0.1.13")
	if got, _ := r.NotApplicable("datetime picker"); got {
		t.Error("an app with a body of its own was judged not applicable")
	}
}

func TestModesPreviewsAndReleases(t *testing.T) {
	dir := repo(t, map[string]string{"ui.go": "package ui\nvar v = &terminu.Viewer{}\nvar f = &terminu.FilePicker{}\n"})
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", dir, "-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false"}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	git("init", "-q")
	git("add", ".")
	git("commit", "-qm", "start")
	r, _ := Scan(dir, "v0.1.13")
	for f, want := range map[string]bool{
		"模式（在 popup 裡）": true, "模式": true, "模式（在 panel 上）": false,
		"有預覽的 finder": true, "有檔案預覽": true, "發過版": false,
	} {
		if r.Has[f] != want {
			t.Errorf("%s: %v, want %v", f, r.Has[f], want)
		}
	}
	if na, missing := r.NotApplicable("發過版"); !na || Why(missing[0]) != "repo 沒有 tag" {
		t.Errorf("an app with no tag: %v %v", na, missing)
	}
	git("tag", "v0.1.0")
	if r, _ = Scan(dir, "v0.1.13"); !r.Has["發過版"] {
		t.Error("an app with a tag was taken as never released")
	}
}

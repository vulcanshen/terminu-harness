package scan

import (
	"os"
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

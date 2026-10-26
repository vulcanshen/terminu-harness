package tdp

import (
	"os"
	"testing"
)

func TestHeadingIDs(t *testing.T) {
	md := []byte("# panel（框）\n" +
		"## panel.1 框\n" +
		"### panel.2 被指到的 panel\n" +
		"## P5 明文放行\n" +
		"### terms.3 Panel\n" +
		"### K12 導覽字母保留給移動\n" +
		"## K Core key\n" +
		"## 術語\n" +
		"```\n" +
		"## v0.1.2 — tdp v0.6.11\n" +
		"## project.99 in a code block\n" +
		"```\n" +
		"#### color-picker.4 色票\n")
	got := headingIDs(md)
	for _, id := range []string{"panel.1", "panel.2", "P5", "terms.3", "K12", "color-picker.4"} {
		if _, ok := got[id]; !ok {
			t.Errorf("%s not found in %v", id, got)
		}
	}
	for _, id := range []string{"K", "術語", "v0.1.2", "project.99"} {
		if _, ok := got[id]; ok {
			t.Errorf("%s should not be an ID", id)
		}
	}
	if got["panel.2"] != "被指到的 panel" {
		t.Errorf("title of panel.2 = %q", got["panel.2"])
	}
}

func TestFileName(t *testing.T) {
	for p, want := range map[string]string{
		"components/dialog/confirm-zh_TW.md": "confirm",
		"components/input/README-zh_TW.md":   "input",
		"components/README-zh_TW.md":         "components",
		"project/README-zh_TW.md":            "project",
		"principle/rules-zh_TW.md":           "",
	} {
		if got := fileName(p); got != want {
			t.Errorf("fileName(%s) = %q, want %q", p, got, want)
		}
	}
}

// TestIDsOfTDP reads the IDs of a real tdp; it needs a clone of terminu.
func TestIDsOfTDP(t *testing.T) {
	repo := os.Getenv("TERMINU_HARNESS__TDP_REPO")
	if repo == "" {
		t.Skip("no terminu clone (TERMINU_HARNESS__TDP_REPO)")
	}
	ids, err := Source{Repo: repo}.IDs("v0.10.1")
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"K4", "P5", "terms.3", "panel.5", "confirm.2", "input.10", "project.3", "confirm", "input", "components", "project"} {
		if !ids.Has(id) {
			t.Errorf("%s not in tdp v0.10.1", id)
		}
	}
	for _, id := range []string{"panel.99", "rules", "principle", "items"} {
		if ids.Has(id) {
			t.Errorf("%s should not be citable", id)
		}
	}
}

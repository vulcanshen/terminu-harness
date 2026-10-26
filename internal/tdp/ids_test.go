package tdp

import (
	"os"
	"strings"
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

func TestSectionText(t *testing.T) {
	md := []byte("# panel\n" +
		"## panel.1 框\n" +
		"框的規定\n" +
		"```\n" +
		"## not a heading\n" +
		"```\n" +
		"#### 例\n" +
		"例子\n" +
		"### panel.2 被指到的 panel\n" +
		"被指到\n" +
		"### panel.3 模式名\n" +
		"## panel.4 膠囊串\n" +
		"膠囊\n")
	text, children := sectionText(md, "panel.1")
	want := "## panel.1 框\n框的規定\n```\n## not a heading\n```\n#### 例\n例子"
	if text != want {
		t.Errorf("panel.1 text:\n%s\nwant:\n%s", text, want)
	}
	if strings.Join(children, " ") != "panel.2 panel.3" {
		t.Errorf("children of panel.1: %v", children)
	}
	if text, _ := sectionText(md, "panel.2"); text != "### panel.2 被指到的 panel\n被指到" {
		t.Errorf("panel.2 text: %q", text)
	}
	if text, _ := sectionText(md, "panel.4"); text != "## panel.4 膠囊串\n膠囊" {
		t.Errorf("panel.4 text: %q", text)
	}
}

func TestItemRows(t *testing.T) {
	items := []byte("| 編號 | 出處 | 檢查什麼 |\n|---|---|---|\n" +
		"| panel-47 | panel.5 | cursor |\n" +
		"| panel-59 | panel.6 | 排序 |\n" +
		"| panel-90 | panel.5、panel.12 | 兩節 |\n" +
		"| panel-91 | panel.5（圖） | 圖 |\n" +
		"| dp-1 | datetime-picker.2（日期）（圖） | 日期 |\n" +
		"| en-1 | panel.12 (diagram) | english |\n" +
		"| K12-1 | K12 | 字母 |\n" +
		"| F3-1 | F3、K4 | 規則 |\n")
	rows := func(keep func(string) bool) string {
		var ids []string
		for _, r := range itemRows(items, keep) {
			ids = append(ids, strings.TrimSpace(strings.Split(r, "|")[1]))
		}
		return strings.Join(ids, " ")
	}
	is := func(id string) func(string) bool { return func(src string) bool { return src == id } }
	for id, want := range map[string]string{
		"panel.5":           "panel-47 panel-90 panel-91",
		"panel.12":          "panel-90 en-1",
		"datetime-picker.2": "dp-1",
		"K4":                "F3-1",
		"K12":               "K12-1",
	} {
		if got := rows(is(id)); got != want {
			t.Errorf("rows of %s: %s, want %s", id, got, want)
		}
	}
	if got := rows(func(src string) bool { return strings.HasPrefix(src, "panel.") }); got != "panel-47 panel-59 panel-90 panel-91 en-1" {
		t.Errorf("rows of panel: %s", got)
	}
}

// TestSectionsOfTDP reads sections of a real tdp; it needs a clone of terminu.
func TestSectionsOfTDP(t *testing.T) {
	repo := os.Getenv("TERMINU_HARNESS__TDP_REPO")
	if repo == "" {
		t.Skip("no terminu clone (TERMINU_HARNESS__TDP_REPO)")
	}
	secs, err := Source{Repo: repo}.Sections("v0.10.1", []string{"tdp panel.5", "K12", "confirm"})
	if err != nil {
		t.Fatal(err)
	}
	if len(secs) != 3 {
		t.Fatalf("%d sections", len(secs))
	}
	if !strings.HasPrefix(secs[0].Text, "## panel.5 一般清單") || strings.Contains(secs[0].Text, "### panel.6") || len(secs[0].Items) == 0 {
		t.Errorf("panel.5: %d items, children %v, text starts %q", len(secs[0].Items), secs[0].Children, secs[0].Text[:40])
	}
	if !strings.Contains(strings.Join(secs[0].Children, " "), "panel.6") {
		t.Errorf("children of panel.5: %v", secs[0].Children)
	}
	if !strings.HasPrefix(secs[1].Text, "### K12 ") || len(secs[1].Items) == 0 {
		t.Errorf("K12: %d items", len(secs[1].Items))
	}
	if !strings.Contains(secs[2].Text, "## confirm.3 按鍵") || len(secs[2].Items) < 10 {
		t.Errorf("confirm: %d items", len(secs[2].Items))
	}
	if _, err := (Source{Repo: repo}).Sections("v0.10.1", []string{"panel.99"}); err == nil {
		t.Error("panel.99 found")
	}
}

package tdp

import (
	"os"
	"strings"
	"testing"
)

// TestHarnessIDs: harness's project/ has what tdp's had up to v1.1.10: the
// sections project.1–8, the page as project, the templates by path, and the
// project- items; the checklist itself is not a citable file.
func TestHarnessIDs(t *testing.T) {
	ids := HarnessIDs()
	for _, id := range []string{"project.1", "project.2", "project.7", "project.8", "project", "project/install.sh",
		"project/family.mk", "project/version.go", "project/gitignore", "project-1", "project-13", "project-37", "project-39"} {
		if !ids.Has(id) {
			t.Errorf("%s not in harness's project/", id)
		}
	}
	for _, id := range []string{"project-17", "project.9", "items", "project/items-zh_TW.md", "project/README-zh_TW.md", "K4"} {
		if ids.Has(id) {
			t.Errorf("%s in harness's project/", id)
		}
	}
	if n := len(ids.Items); n != 38 {
		t.Errorf("%d project- items, want 38", n)
	}
	if !ids.Project() {
		t.Errorf("harness's IDs have no project")
	}
}

func TestCited(t *testing.T) {
	for c, want := range map[string]struct {
		harness bool
		id      string
	}{
		"harness project.3":           {true, "project.3"},
		" harness project/install.sh": {true, "project/install.sh"},
		"tdp K4":                      {false, "K4"},
		"K4":                          {false, "K4"},
		"tdp project.3":               {false, "project.3"},
	} {
		if h, id := Cited(c); h != want.harness || id != want.id {
			t.Errorf("Cited(%q) = %v, %q", c, h, id)
		}
	}
}

// TestSectionsOfHarness: harness-prefixed IDs come from harness's project/,
// marked so, an item followed by its source; the rest from tdp.
func TestSectionsOfHarness(t *testing.T) {
	secs, err := Source{}.Sections("v9.9.9", []string{"harness project-13", "harness project/install.sh", "harness project"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, s := range secs {
		if !s.Harness {
			t.Errorf("%s not marked harness's", s.ID)
		}
		got = append(got, s.ID)
	}
	if strings.Join(got, " ") != "project-13 project/install.sh project project.3" {
		t.Errorf("sections %q", got)
	}
	if !strings.HasPrefix(secs[0].Text, "| project-13 | project.3 |") {
		t.Errorf("item row %q", secs[0].Text)
	}
	if !strings.Contains(secs[1].Text, "terminu family template (harness project/install.sh)") {
		t.Errorf("install.sh:\n%s", secs[1].Text)
	}
	if len(secs[2].Items) != 38 {
		t.Errorf("the page has %d items", len(secs[2].Items))
	}
	if !strings.HasPrefix(secs[3].Text, "## project.3 目錄") || len(secs[3].Items) == 0 {
		t.Errorf("project.3: %q, %d items", secs[3].Text[:40], len(secs[3].Items))
	}
	if _, err := (Source{}).Sections("v9.9.9", []string{"harness project.99"}); err == nil || !strings.Contains(err.Error(), "harness's project/") {
		t.Errorf("project.99: %v", err)
	}
}

// TestProjectLeftTDP: from tdp v1.2.0-beta.1 on, tdp has no project/; a tdp
// project ID says where it went. Before, tdp's is read as it was.
func TestProjectLeftTDP(t *testing.T) {
	repo := os.Getenv("TERMINU_HARNESS__TDP_REPO")
	if repo == "" {
		t.Skip("no terminu clone (TERMINU_HARNESS__TDP_REPO)")
	}
	src := Source{Repo: repo}
	_, err := src.Sections("v1.2.0-beta.1", []string{"tdp project.3"})
	if err == nil || !strings.Contains(err.Error(), `cited as "harness project.3"`) {
		t.Errorf("tdp project.3 at v1.2.0-beta.1: %v", err)
	}
	ids, err := src.IDs("v1.2.0-beta.1")
	if err != nil {
		t.Fatal(err)
	}
	if ids.Project() || ids.Has("project/install.sh") || ids.Has("project-13") {
		t.Errorf("tdp v1.2.0-beta.1 has project")
	}
	old, err := src.IDs("v1.1.10")
	if err != nil {
		t.Fatal(err)
	}
	if !old.Project() || !old.Has("project/install.sh") || !old.Has("project-13") {
		t.Errorf("tdp v1.1.10 has no project")
	}
	secs, err := src.Sections("v1.1.10", []string{"tdp project.3", "harness project.3"})
	if err != nil {
		t.Fatal(err)
	}
	if len(secs) != 2 || secs[0].Harness || !secs[1].Harness || secs[0].Path != "project/README-zh_TW.md" {
		t.Errorf("both at v1.1.10: %+v", secs)
	}
}

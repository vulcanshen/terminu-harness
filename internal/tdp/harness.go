package tdp

import (
	"io/fs"
	"sort"
	"strings"

	harness "github.com/vulcanshen/terminu-harness"
)

// HarnessPrefix is how an app cites terminu project, harness's project/:
// harness project.3, harness project, harness project/install.sh, harness
// project-13 (from tdp v1.2.0-beta.1 on; tdp project.3 before).
const HarnessPrefix = "harness "

// HarnessItemsFile is the checklist of harness's project/.
const HarnessItemsFile = "project/items-zh_TW.md"

// harnessTree is harness's project/, as the plugin carries it.
func harnessTree() tree {
	var paths []string
	_ = fs.WalkDir(harness.Project, "project", func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			paths = append(paths, p)
		}
		return err
	})
	sort.Strings(paths)
	return tree{at: "harness's project/", paths: paths, file: harness.Project.ReadFile, items: HarnessItemsFile}
}

// HarnessIDs is what an app can cite in harness's project/, without the
// harness prefix: sections (project.3), the page (project), templates
// (project/install.sh) and checklist items (project-13).
func HarnessIDs() IDs {
	ids, err := harnessTree().ids()
	if err != nil {
		panic(err) // the files are in the binary
	}
	return ids
}

// Cited splits a citation into where it points and the ID: "harness
// project.3" is harness's project.3, "tdp K4" and "K4" tdp's K4.
func Cited(c string) (inHarness bool, id string) {
	c = strings.TrimSpace(c)
	if id, ok := strings.CutPrefix(c, HarnessPrefix); ok {
		return true, strings.TrimSpace(id)
	}
	return false, strings.TrimSpace(strings.TrimPrefix(c, "tdp "))
}

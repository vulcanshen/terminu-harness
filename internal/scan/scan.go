// Package scan reads an app's Go source for the terminu-go components it
// uses, so the checklist items about features it does not have can be told
// not applicable before a judge looks (user, 2026-10-19: split by the code).
package scan

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Result is what the source shows: for each feature of the checklist's
// 適用 table the scan can tell, whether the app has it and how it is known;
// and the app's own bodies and popups, which may do any feature their own way.
type Result struct {
	Tgo      string            `json:"tgo"`
	Has      map[string]bool   `json:"has"`
	Evidence map[string]string `json:"evidence"`
	Custom   []string          `json:"custom,omitempty"`
}

// features are the 適用 names the scan can tell, by the terminu-go code that
// makes them. One the app has none of these for, it does not have.
var features = map[string]string{
	"多畫面":             `AddScreen\(`, // counted: two or more
	"panel tab":       `Tabs:\s*\[\]string`,
	"一列一個值的 panel":    `terminu\.Values\b`,
	"分組清單":            `\bGroup:`,
	"明細 panel":        `\bDetails:`,
	"位置列":             `\.Location\s*=|\bLocation:`,
	"排序":              `SortBy\(|terminu\.SortKey`,
	"標記":              `Marks:\s*true`,
	"panel filter":    `Filterable:\s*true`,
	"串流內容":            `Stream:\s*true|Follow:\s*true`,
	"loading":         `\bLoading\b|LoadingFrame|LoadSelect\(|LoadChecks\(`,
	"模式（在 panel 上）":   `terminu\.Mode\b|\.EnterMode\(`,
	"停用":              `\bDisabled:`,
	"長時間工作":           `\.Progress\b|SubmitLong`,
	"PTY":             `terminu\.Terminal\b`,
	"terminal popup":  `OpenTerminal\(`,
	"表單":              `terminu\.Form\b`,
	"viewer":          `terminu\.Viewer\b`,
	"toast":           `\.Toast\(`,
	"error popup":     `\.Error\(\s*[^)\s]|\.ErrorAbout\(`, // not err.Error()
	"confirm":         `\.Confirm\(|terminu\.Ask\{|\bAsk:\s|\bAskFor:\s`,
	"text":            `terminu\.TextInput\b`,
	"password":        `terminu\.PasswordInput\b`,
	"number":          `terminu\.NumberInput\b`,
	"select":          `terminu\.Select\b`,
	"checkbox":        `terminu\.Checks\b|Multi:\s*true`,
	"slider":          `terminu\.Slider\b`,
	"textarea":        `terminu\.Textarea\b`,
	"color picker":    `terminu\.ColorPicker\b`,
	"datetime picker": `terminu\.DatePicker\b`,
	"file-picker":     `terminu\.FilePicker\b`,
	"reorder":         `terminu\.Reorder\b`,
	"mouse":           `WithMouse|MouseMsg`,
}

// custom is code that draws a body, a popup or a whole screen (App.Canvas)
// of the app's own.
var custom = regexp.MustCompile(`func \([^)]*\) Draw\([^)]*\) \[\]string|\.Popup\(|\bCanvas:\s|\.Canvas\s*=`)

// Scan reads the repo; nil when the app is not written with terminu-go.
func Scan(repo, tgo string) (*Result, error) {
	if tgo == "" {
		return nil, nil
	}
	var src strings.Builder
	files := map[int]string{} // offset in src → file, to name evidence
	r := &Result{Tgo: tgo, Has: map[string]bool{}, Evidence: map[string]string{}}
	err := filepath.WalkDir(repo, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() && (name == "vendor" || name == "testdata" || strings.HasPrefix(name, ".")) && p != repo {
			return filepath.SkipDir
		}
		if d.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(repo, p)
		files[src.Len()] = rel
		src.Write(b)
		src.WriteByte('\n')
		if custom.Match(b) {
			r.Custom = append(r.Custom, rel)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	text := src.String()
	where := func(at int) string {
		best, file := -1, ""
		for off, f := range files {
			if off <= at && off > best {
				best, file = off, f
			}
		}
		return file
	}
	for f, pat := range features {
		locs := regexp.MustCompile(pat).FindAllStringIndex(text, -1)
		has := len(locs) > 0
		if f == "多畫面" {
			has = len(locs) >= 2
		}
		r.Has[f] = has
		if has {
			r.Evidence[f] = where(locs[0][0])
		}
	}
	// A finder is a finder or one of its kinds; an input, any input. A plain
	// finder and a file-picker show a preview, maybe of files; select and
	// checkbox popups none. A color picker has a select (its R, G and B).
	r.Has["select"] = r.Has["select"] || r.Has["color picker"]
	finder := regexp.MustCompile(`terminu\.Finder\b`).MatchString(text)
	r.Has["finder"] = r.Has["select"] || r.Has["checkbox"] || r.Has["file-picker"] || finder
	r.Has["有預覽的 finder"] = r.Has["file-picker"] || finder
	r.Has["有檔案預覽"] = r.Has["有預覽的 finder"]
	r.Has["input"] = false
	for _, f := range []string{"text", "password", "number", "select", "checkbox", "slider", "textarea", "color picker", "datetime picker", "file-picker", "reorder"} {
		r.Has["input"] = r.Has["input"] || r.Has[f]
	}
	// Modes (Rules K11): the app's own on a panel, and those components
	// bring: a terminal's history and selection, wherever it is; a viewer's
	// selection and a reorder's Drag in a popup.
	r.Has["模式（在 panel 上）"] = r.Has["模式（在 panel 上）"] || r.Has["PTY"]
	r.Has["模式（在 popup 裡）"] = r.Has["viewer"] || r.Has["reorder"] || r.Has["PTY"] || r.Has["terminal popup"]
	r.Has["模式"] = r.Has["模式（在 panel 上）"] || r.Has["模式（在 popup 裡）"]
	if b, err := os.ReadFile(filepath.Join(repo, ".goreleaser.yaml")); err == nil {
		r.Has["出 Windows 版"] = regexp.MustCompile(`(?m)^\s*-\s*windows\b`).Match(b)
	}
	// Released means a tag and a GitHub release: no tag, surely not.
	if out, err := exec.Command("git", "-C", repo, "tag", "--list").Output(); err == nil {
		r.Has["發過版"] = strings.TrimSpace(string(out)) != ""
	}
	sort.Strings(r.Custom)
	return r, nil
}

// Why says how the scan knows the app lacks feature f.
func Why(f string) string {
	switch f {
	case "發過版":
		return "repo 沒有 tag"
	case "出 Windows 版":
		return ".goreleaser.yaml 沒有 windows"
	}
	return "原始碼沒用到 terminu-go 做它的元件"
}

// NotApplicable reports an item whose 適用 names only features the scan is
// sure the app has none of: every alternative (split by 、) needs one (of
// those joined by ＋) the app lacks. 都要, a feature the scan can't tell, or
// an app with bodies or popups of its own, and it is not sure.
func (r *Result) NotApplicable(applies string) (bool, []string) {
	applies = strings.TrimSpace(applies)
	if r == nil || len(r.Custom) > 0 || applies == "" || applies == "都要" {
		return false, nil
	}
	var missing []string
	for _, alt := range strings.Split(applies, "、") {
		lacks := ""
		for _, f := range strings.Split(alt, "＋") {
			f = strings.TrimSpace(f)
			has, known := r.Has[f]
			if !known {
				return false, nil
			}
			if !has && lacks == "" {
				lacks = f
			}
		}
		if lacks == "" {
			return false, nil // this alternative the app has
		}
		missing = append(missing, lacks)
	}
	return true, missing
}

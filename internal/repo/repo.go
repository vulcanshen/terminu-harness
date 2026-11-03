// Package repo checks a family member's repository against terminu project,
// harness's project/: the templates of project.2 as harness has them, but for
// what an app may change in each, and the files, Makefile and CHANGELOG
// checked by reading them alone. It came from terminu-go's tdptest with the
// templates, at tdp v1.2.0-beta.1.
package repo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	harness "github.com/vulcanshen/terminu-harness"
)

// Family is the account of the family members: terminu project binds them
// alone (project.1).
const Family = "github.com/vulcanshen/"

// libModule is terminu-go's own module: it is not an app (tdp E8).
const libModule = "github.com/vulcanshen/terminu-go"

// Problem is something in the repository to fix.
type Problem struct {
	File string // the file it is about, relative to the repository
	Msg  string
}

// Check checks the repository at root. An app made after tdp, its module not
// under github.com/vulcanshen, keeps its repository as it likes (project.1):
// nothing is checked, and so for the lib.
func Check(root string) []Problem {
	r := &repo{root: root}
	if !r.family() {
		return r.problems
	}
	r.files()
	r.claudeFiles()
	r.templates()
	r.makefile()
	r.changelog()
	return r.problems
}

// repo is an app's repository as Check reads it.
type repo struct {
	root, name string
	problems   []Problem
}

func (r *repo) add(file, format string, a ...any) {
	r.problems = append(r.problems, Problem{file, fmt.Sprintf(format, a...)})
}

func (r *repo) read(rel string) (string, bool) {
	b, err := os.ReadFile(filepath.Join(r.root, filepath.FromSlash(rel)))
	if err != nil {
		return "", false
	}
	return strings.ReplaceAll(string(b), "\r\n", "\n"), true
}

func (r *repo) has(rel string) bool {
	_, err := os.Stat(filepath.Join(r.root, filepath.FromSlash(rel)))
	return err == nil
}

var (
	moduleRE = regexp.MustCompile(`(?m)^module\s+(\S+)`)
	majorRE  = regexp.MustCompile(`/v\d+$`)
)

// family reads go.mod's module path and says whether the repository is a
// family member's, the module under github.com/vulcanshen; the app's name is
// the part after it. A member's module is github.com/vulcanshen/<app>, from
// v2 on ending in /vN, the major version of the latest release its CHANGELOG
// names (project.3, project-13).
func (r *repo) family() bool {
	mod := module(r.root)
	if !isFamily(mod) {
		return false
	}
	r.name, _, _ = strings.Cut(strings.TrimPrefix(majorRE.ReplaceAllString(mod, ""), Family), "/")
	want := Family + r.name
	major, latest := r.release()
	if major >= 2 {
		want += "/v" + strconv.Itoa(major)
	}
	if mod != want {
		why := "CHANGELOG.md 還沒有發過版"
		if latest != "" {
			why = "CHANGELOG.md 寫的最新一版是 " + latest
		}
		r.add("go.mod", "module 是 %s，不是 %s：%s；v2 起 module 結尾帶大版號 /vN，go install %s/cmd/%s@latest 才裝得到最新的 release（harness project.3、project-13）",
			mod, want, why, want, r.name)
	}
	return true
}

// IsFamily reports whether the repository at root is a family member's,
// which Check reads: its go.mod's module under github.com/vulcanshen, and not
// terminu-go itself.
func IsFamily(root string) bool { return isFamily(module(root)) }

func isFamily(mod string) bool {
	return mod != "" && mod != libModule && strings.HasPrefix(mod, Family)
}

// module is the module path go.mod at root names, "" when there is none.
func module(root string) string {
	b, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return ""
	}
	if m := moduleRE.FindSubmatch(b); m != nil {
		return string(m[1])
	}
	return ""
}

// release is the latest release CHANGELOG.md names, its first heading with
// a version, and that version's major; "" and 0 when it names none.
func (r *repo) release() (int, string) {
	s, _ := r.read("CHANGELOG.md")
	m := regexp.MustCompile(`(?m)^## (v(\d+)\.\d+\.\d+\S*)`).FindStringSubmatch(s)
	if m == nil {
		return 0, ""
	}
	n, _ := strconv.Atoi(m[2])
	return n, m[1]
}

// files are the ones project.3 puts in every app.
func (r *repo) files() {
	for _, f := range []struct{ path, item string }{
		{"cmd/" + r.name + "/main.go", "project-12：main package"},
		{"go.sum", "project-15"},
		{"LICENSE", "project-16"},
		{"docs/icon.svg", "project-18：splash 畫的家族 mark 的正本"},
	} {
		if !r.has(f.path) {
			r.add(f.path, "沒有這個檔（harness project.3、%s）", f.item)
		}
	}
}

// claudeFiles checks that git tracks no CLAUDE.md and, under .claude/, only
// harness's .claude/rules/terminu.md (project.3, project-37). Outside a git
// work tree, or without git, it is not checked.
func (r *repo) claudeFiles() {
	out, err := exec.Command("git", "-C", r.root, "ls-files", "-z").Output()
	if err != nil {
		return
	}
	for _, f := range strings.Split(string(out), "\x00") {
		if f == ".claude/rules/terminu.md" {
			continue
		}
		if f == "CLAUDE.md" || strings.HasSuffix(f, "/CLAUDE.md") || strings.HasPrefix(f, ".claude/") {
			r.add(f, "git 追蹤了這個檔：Claude Code 的檔不進 git，只有 harness 的 .claude/rules/terminu.md 例外；git rm --cached %s 留著檔、不再追蹤（harness project.3、project-37）", f)
		}
	}
}

// line is a line of a file and where it was, as the diff shows it.
type line struct {
	text string
	n    int
}

func lines(s string) []line {
	s = strings.TrimSuffix(s, "\n")
	var out []line
	for i, l := range strings.Split(s, "\n") {
		out = append(out, line{l, i + 1})
	}
	return out
}

// template is one of project.2's templates and what an app may change in it,
// as its last column says.
type template struct {
	src, dst, item string
	may            string // what may differ, for the message
	// named: the app's name is written as <app> before comparing.
	named bool
	// fix makes the app's lines and the template's comparable: what the app
	// may change is made the same on both sides.
	fix func(app, ours []line) ([]line, []line)
}

func (r *repo) templates() {
	install, _ := r.read("install.sh")
	windows := regexp.MustCompile(`(?m)^WINDOWS="yes"$`).MatchString(install)
	known := install != ""
	for _, tm := range []template{
		{src: "family.mk", dst: "family.mk", item: "project-1", may: "裡面什麼都不能改"},
		{src: "release.yml", dst: ".github/workflows/release.yml", item: "project-2、project-8",
			may: "只有出 Windows 版的 app 在測試的 matrix 加的 windows-latest 可以不一樣",
			fix: func(app, ours []line) ([]line, []line) {
				switch {
				case !known: // either
					return matrix(app, withWindows, noWindows), ours
				case windows:
					return app, matrix(ours, noWindows, withWindows)
				}
				return app, ours
			}},
		{src: "goreleaser.yaml", dst: ".goreleaser.yaml", item: "project-3、project-9", named: true,
			may: "只有 <app>、每個標著 app: 的註解底下那一個鍵的值（註解講的鍵範本裡還沒有時，在註解那裡加上的那個鍵）、出 Windows 版（install.sh 的 WINDOWS=\"yes\"）時標著 Windows 的那幾列與它們底下的註解列拿掉的行首「# 」可以不一樣；註解與 Windows 的標記照原樣",
			fix: func(app, ours []line) ([]line, []line) {
				if !known {
					windows = sameWindows(app)
				}
				return appRegions(app), appRegions(openWindows(ours, windows))
			}},
		{src: "install.sh", dst: "install.sh", item: "project-4、project-10", named: true,
			may: `只有 <app>、WINDOWS（出 Windows 版寫 "yes"）與標著 app's own 的段落可以不一樣`,
			fix: func(app, ours []line) ([]line, []line) {
				for i, l := range app {
					if l.text == `WINDOWS="yes"` {
						app[i].text = `WINDOWS="no"`
					}
				}
				return ownBlocks(app), ownBlocks(ours)
			}},
		{src: "uninstall.sh", dst: "uninstall.sh", item: "project-5", named: true,
			may: "只有 <app> 與標著 app's own 的段落可以不一樣；註解裡的 <APP> 照原樣",
			fix: func(app, ours []line) ([]line, []line) { return ownBlocks(app), ownBlocks(ours) }},
		{src: "version.go", dst: "internal/version/version.go", item: "project-6", named: true, may: "只有 <app> 可以不一樣"},
		{src: "gitignore", dst: ".gitignore", item: "project-7", named: true,
			may: "只有 <app> 可以不一樣，app 自己的行接在最後",
			fix: func(app, ours []line) ([]line, []line) {
				if len(app) > len(ours) {
					app = app[:len(ours)]
				}
				return app, ours
			}},
	} {
		r.template(tm)
	}
	if windows {
		for _, f := range []string{"install.ps1", "uninstall.ps1"} {
			if !r.has(f) {
				r.add(f, "install.sh 寫了有 Windows 版（WINDOWS=\"yes\"），但是沒有這個檔（harness project.3、project-19）")
			}
		}
	}
}

func (r *repo) template(tm template) {
	b, err := harness.Project.ReadFile("project/" + tm.src)
	if err != nil {
		panic(err)
	}
	again := fmt.Sprintf("用 tdp_section 讀 harness project/%s，重新複製到 %s", tm.src, tm.dst)
	s, ok := r.read(tm.dst)
	if !ok {
		r.add(tm.dst, "沒有這個檔（harness project.2、%s）：%s", tm.item, again)
		return
	}
	if tm.named {
		s = r.unname(s)
	}
	app, ours := lines(s), lines(strings.ReplaceAll(string(b), "\r\n", "\n"))
	if tm.fix != nil {
		app, ours = tm.fix(app, ours)
	}
	if d := diff(ours, app, "project/"+tm.src, tm.dst); d != "" {
		named := ""
		if tm.named {
			named = "，app 的名字寫成 <app>"
		}
		r.add(tm.dst, "跟 harness 的 project/%s 不一樣（harness project.2、%s）：%s；%s。\n%s（- harness 的，+ app 的%s）",
			tm.src, tm.item, again, tm.may, d, named)
	}
}

// unname writes the app's name as <app>, as in the templates. The <APP> in
// uninstall.sh's comment is not <app>: it stays as it is (project.2).
func (r *repo) unname(s string) string {
	return regexp.MustCompile(`\b`+regexp.QuoteMeta(r.name)+`\b`).ReplaceAllLiteralString(s, "<app>")
}

// The test matrix of release.yml, as the template has it and with
// windows-latest added (project.2).
const (
	noWindows   = "os: [ubuntu-latest, macos-latest]"
	withWindows = "os: [ubuntu-latest, macos-latest, windows-latest]"
)

// matrix writes the test matrix from as to in a copy of ls.
func matrix(ls []line, from, to string) []line {
	out := append([]line(nil), ls...)
	for i, l := range out {
		if strings.TrimSpace(l.text) == from {
			out[i].text = strings.Replace(l.text, from, to, 1)
		}
	}
	return out
}

// ownBlocks leaves out what is between "# --- app's own: …" and "# --- end
// of the app's own ---": the app's (install.sh, uninstall.sh).
func ownBlocks(ls []line) []line {
	var out []line
	in := false
	for _, l := range ls {
		t := strings.TrimSpace(l.text)
		switch {
		case strings.HasPrefix(t, "# --- app's own:"):
			in = true
		case t == "# --- end of the app's own ---":
			in = false
		case in:
			continue
		}
		out = append(out, l)
	}
	return out
}

var (
	// appMark is a comment marked "app:" (.goreleaser.yaml); one that names
	// its key first ("app: dependencies, …") is for a key the template may
	// not have yet.
	appMark = regexp.MustCompile(`^# app: (?:([a-z_]+),)?`)
	yamlKey = regexp.MustCompile(`^(\s*)([A-Za-z_][A-Za-z0-9_]*):`)
)

// appRegions leaves out what each comment marked "app:" lets the app write
// (.goreleaser.yaml, project.2): the value of the key right below the
// comment, a list or a whole block included, the key itself kept; or, for a
// comment naming a key, that key, which the app adds there when the template
// has none. The comment, all its lines, and any other key stay.
func appRegions(ls []line) []line {
	var out []line
	for i := 0; i < len(ls); i++ {
		out = append(out, ls[i])
		m := appMark.FindStringSubmatch(strings.TrimSpace(ls[i].text))
		if m == nil {
			continue
		}
		for i+1 < len(ls) && strings.HasPrefix(strings.TrimSpace(ls[i+1].text), "#") {
			i++
			out = append(out, ls[i])
		}
		if i+1 == len(ls) {
			break
		}
		k := yamlKey.FindStringSubmatch(ls[i+1].text)
		if k == nil || m[1] != "" && k[2] != m[1] {
			continue // not the app's: test: under "app: dependencies, …"
		}
		i++
		if m[1] == "" {
			out = append(out, line{k[1] + k[2] + ":", ls[i].n})
		}
		i = valueEnd(ls, i, len(k[1]))
	}
	return out
}

// valueEnd is the last line of the value of the key at ls[i], indented in:
// the lines below indented more, a list at the key's own indent, and the
// blank lines between them.
func valueEnd(ls []line, i, in int) int {
	end := i
	for j := i + 1; j < len(ls); j++ {
		t := strings.TrimSpace(ls[j].text)
		switch {
		case t == "":
		case indent(ls[j].text) > in || indent(ls[j].text) == in && strings.HasPrefix(t, "- "):
			end = j
		default:
			return end
		}
	}
	return end
}

func indent(s string) int { return len(s) - len(strings.TrimLeft(s, " ")) }

var windowsMark = regexp.MustCompile(`\s+# Windows$`)

// openWindows uncomments the template's lines marked Windows and the
// commented lines of the block right under each, for an app with a Windows
// build: only the leading "# " goes, the mark stays (project.2).
func openWindows(ls []line, windows bool) []line {
	if !windows {
		return ls
	}
	out := append([]line(nil), ls...)
	for i := 0; i < len(out); i++ {
		if !windowsMark.MatchString(out[i].text) || !strings.HasPrefix(strings.TrimSpace(out[i].text), "#") {
			continue
		}
		in := indent(out[i].text)
		out[i].text = uncomment(out[i].text, in)
		for i+1 < len(out) && indent(out[i+1].text) == in && strings.HasPrefix(strings.TrimSpace(out[i+1].text), "#") &&
			!windowsMark.MatchString(out[i+1].text) {
			i++
			out[i].text = uncomment(out[i].text, in)
		}
	}
	return out
}

func uncomment(s string, in int) string {
	rest := strings.TrimPrefix(s[in:], "#")
	rest = strings.TrimPrefix(rest, " ")
	if rest == "" {
		return ""
	}
	return s[:in] + rest
}

// sameWindows, with no install.sh to say, takes .goreleaser.yaml's Windows
// lines as uncommented when its builds have windows.
func sameWindows(app []line) bool {
	for _, l := range app {
		if t := strings.TrimSpace(windowsMark.ReplaceAllString(l.text, "")); t == "- windows" {
			return true
		}
	}
	return false
}

var (
	assignment = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*\s*(:=|::=|\?=|\+=|=)`)
	target     = regexp.MustCompile(`^([A-Za-z0-9_./-]+(?:\s+[A-Za-z0-9_./-]+)*)\s*:(?:[^=]|$)`)
)

// makefile checks the app's Makefile against project.5: its first lines, the
// variables, include family.mk, then the app's own section, last.
func (r *repo) makefile() {
	s, ok := r.read("Makefile")
	if !ok {
		r.add("Makefile", "沒有這個檔（harness project.5、project-23）")
		return
	}
	ls := strings.Split(strings.TrimSuffix(s, "\n"), "\n")
	head := []string{"# " + r.name + " — build / run / test / package（家族共用的 target 在 family.mk）", "#", "# 跑 `make`（或 `make help`）列出所有指令。"}
	for i, want := range head {
		if i >= len(ls) || ls[i] != want {
			got := ""
			if i < len(ls) {
				got = ls[i]
			}
			r.add("Makefile", "第 %d 行是 %q，不是 %q：開頭幾行照 harness project.5（project-26）", i+1, got, want)
			break
		}
	}
	inc := -1
	for i, l := range ls {
		if strings.TrimSpace(l) == "include family.mk" {
			inc = i
			break
		}
	}
	if inc < 0 {
		r.add("Makefile", "沒有 include family.mk（harness project.5、project-23）")
		return
	}
	binary := false
	for i := 0; i < inc; i++ {
		t := strings.TrimSpace(ls[i])
		switch {
		case t == "" || strings.HasPrefix(t, "#") || i > 0 && strings.HasSuffix(ls[i-1], "\\"):
		case assignment.MatchString(t):
			binary = binary || regexp.MustCompile(`^BINARY\s*:=\s*`+regexp.QuoteMeta(r.name)+`$`).MatchString(t)
		default:
			r.add("Makefile", "第 %d 行在 include family.mk 之前，卻不是變數：%q（harness project.5、project-24）", i+1, ls[i])
		}
	}
	if !binary {
		r.add("Makefile", "include family.mk 之前沒有設 BINARY := %s（harness project.5、project-25）", r.name)
	}
	own, section := -1, "##@ "+r.name+" 自己的"
	for i := inc + 1; i < len(ls); i++ {
		t := strings.TrimSpace(ls[i])
		switch {
		case own < 0 && t == section:
			own = i
		case own < 0 && t != "" && !strings.HasPrefix(t, "#"):
			r.add("Makefile", "第 %d 行在 include family.mk 之後、%q 之前：%q；app 自己的 target 寫在那一段，放最後（harness project.5、project-23）", i+1, section, ls[i])
		case own >= 0 && strings.HasPrefix(t, "##@"):
			r.add("Makefile", "第 %d 行在 %q 之後又開了一段：%q；它是最後一段（harness project.5、project-23）", i+1, section, ls[i])
		}
	}
	if own < 0 {
		r.add("Makefile", "include family.mk 之後沒有 %q 這一段（harness project.5、project-23）", section)
	}
	family, _ := harness.Project.ReadFile("project/family.mk")
	ours := targets(strings.Split(string(family), "\n"))
	for i, l := range ls {
		for name := range targets([]string{l}) {
			if ours[name] {
				r.add("Makefile", "第 %d 行做了 target %s，family.mk 已經有了（harness project.5、project-27）", i+1, name)
			}
		}
	}
}

// targets are the target names the rule lines among ls make, special ones
// (.PHONY) left out.
func targets(ls []string) map[string]bool {
	out := map[string]bool{}
	for _, l := range ls {
		m := target.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		for _, name := range strings.Fields(m[1]) {
			if !strings.HasPrefix(name, ".") {
				out[name] = true
			}
		}
	}
	return out
}

// version is a tag's version, a pre-release one too (v1.1.0-beta.1).
const version = `v\d+\.\d+\.\d+(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?`

var heading = regexp.MustCompile(`^## ` + version + ` — tdp ` + version + `$`)

// changelog checks CHANGELOG.md against project.6: its first lines, and each
// heading naming a tdp written as a tag, " — tdp ", the version.
func (r *repo) changelog() {
	s, ok := r.read("CHANGELOG.md")
	if !ok {
		r.add("CHANGELOG.md", "沒有這個檔（harness project.6、project-28）")
		return
	}
	ls := strings.Split(s, "\n")
	head := []string{"# Changelog", "", "Each release names the version of terminu design (tdp) it follows."}
	for i, want := range head {
		if i >= len(ls) || ls[i] != want {
			got := ""
			if i < len(ls) {
				got = ls[i]
			}
			r.add("CHANGELOG.md", "第 %d 行是 %q，不是 %q：開頭照 harness project.6（project-28）", i+1, got, want)
			break
		}
	}
	for i, l := range ls {
		switch {
		case !strings.HasPrefix(l, "## "):
		case strings.HasPrefix(l, "## Unreleased") && l != "## Unreleased":
			r.add("CHANGELOG.md", "第 %d 行是 %q：還沒發版的改動寫在只有 ## Unreleased 的標題底下（harness project.6、project-30）", i+1, l)
		case strings.Contains(l, "tdp") && !heading.MatchString(l):
			r.add("CHANGELOG.md", "第 %d 行是 %q，不是 ## vX.Y.Z — tdp vA.B.C：push 的 tag，接著它跟的 tdp（harness project.6、project-29）", i+1, l)
		}
	}
}

// diff is a's lines against b's, the ones that differ with two around them,
// "" when they are the same.
func diff(a, b []line, an, bn string) string {
	n, m := len(a), len(b)
	// lcs[i][j] is the longest common run of a[i:] and b[j:].
	lcs := make([][]int, n+1)
	for i := range lcs {
		lcs[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if a[i].text == b[j].text {
				lcs[i][j] = lcs[i+1][j+1] + 1
			} else {
				lcs[i][j] = max(lcs[i+1][j], lcs[i][j+1])
			}
		}
	}
	type op struct {
		kind byte // ' ', '-', '+'
		l    line
	}
	var ops []op
	i, j := 0, 0
	for i < n || j < m {
		switch {
		case i < n && j < m && a[i].text == b[j].text:
			ops = append(ops, op{' ', a[i]})
			i, j = i+1, j+1
		case i < n && (j == m || lcs[i+1][j] >= lcs[i][j+1]):
			ops = append(ops, op{'-', a[i]})
			i++
		default:
			ops = append(ops, op{'+', b[j]})
			j++
		}
	}
	show := make([]bool, len(ops))
	for k, o := range ops {
		if o.kind != ' ' {
			for x := max(k-2, 0); x <= min(k+2, len(ops)-1); x++ {
				show[x] = true
			}
		}
	}
	var out strings.Builder
	rows := 0
	for k, o := range ops {
		if !show[k] {
			continue
		}
		if k == 0 || !show[k-1] {
			where := an
			if o.kind == '+' {
				where = bn
			}
			fmt.Fprintf(&out, "@@ %s 第 %d 行 @@\n", where, o.l.n)
		}
		if rows++; rows > 60 {
			out.WriteString("…\n")
			break
		}
		fmt.Fprintf(&out, "%c %s\n", o.kind, o.l.text)
	}
	return out.String()
}

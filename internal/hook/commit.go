// Package hook holds the harness's Claude Code hooks.
package hook

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/vulcanshen/terminu-harness/internal/assets"
)

// Trailer is the line a commit message carries when it changes the app's
// code but not what the app does.
const Trailer = "Terminu-Work: unchanged"

// Input is what Claude Code sends a PreToolUse hook on stdin.
type Input struct {
	Cwd       string `json:"cwd"`
	ToolName  string `json:"tool_name"`
	ToolInput struct {
		Command string `json:"command"`
	} `json:"tool_input"`
}

var (
	commitRE = regexp.MustCompile(`(?:^|[\s;&|(])git(?:\s+-C\s+(\S+))?\s+commit(?:\s|$)`)
	addRE    = regexp.MustCompile(`(?:^|[\s;&|(])git(?:\s+-C\s+\S+)?\s+add(?:\s|$)`)
	allRE    = regexp.MustCompile(`\scommit\b[^;&|]*\s(?:--all|-[a-zA-Z]*a[a-zA-Z]*)(?:\s|$)`)
	fileRE   = regexp.MustCompile(`\scommit\b[^;&|]*\s(?:-F|--file)[\s=]+(\S+)`)
)

// Git is what Commit asks git; Files runs git with args in the repo and
// returns the paths it prints.
type Git struct {
	Root  func(dir string) (string, error)
	Files func(root string, args ...string) ([]string, error)
}

// DefaultGit runs git.
var DefaultGit = Git{
	Root: func(dir string) (string, error) {
		out, err := exec.Command("git", "-C", dir, "rev-parse", "--show-toplevel").Output()
		return strings.TrimSpace(string(out)), err
	},
	Files: func(root string, args ...string) ([]string, error) {
		out, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
		return strings.Fields(string(out)), err
	},
}

// Commit decides whether a Bash command may run: a git commit in a repo the
// harness keeps (go.mod requires terminu-go and harness-assets/ is there) that
// changes the app's code without harness-assets/work.yaml, and whose message
// does not carry Trailer, is blocked. It returns why when it blocks.
func Commit(in Input, git Git) (string, bool) {
	cmd := in.ToolInput.Command
	m := commitRE.FindStringSubmatch(cmd)
	if m == nil {
		return "", false
	}
	if strings.Contains(cmd, Trailer) {
		return "", false
	}
	if f := fileRE.FindStringSubmatch(cmd); f != nil {
		p := f[1]
		if !filepath.IsAbs(p) {
			p = filepath.Join(in.Cwd, p)
		}
		if b, err := os.ReadFile(strings.Trim(p, `"'`)); err == nil && strings.Contains(string(b), Trailer) {
			return "", false
		}
	}
	dir := in.Cwd
	if m[1] != "" {
		dir = m[1]
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(in.Cwd, dir)
		}
	}
	root, err := git.Root(dir)
	if err != nil || !kept(root) {
		return "", false
	}

	// What goes into the commit: what is staged, and, when the same command
	// stages more first (git add ... && git commit) or commits with -a, what
	// is changed or new in the working tree too.
	files, _ := git.Files(root, "diff", "--cached", "--name-only")
	if addRE.MatchString(cmd) || allRE.MatchString(cmd) {
		changed, _ := git.Files(root, "diff", "--name-only")
		added, _ := git.Files(root, "ls-files", "--others", "--exclude-standard")
		files = append(append(files, changed...), added...)
	}
	var code []string
	work := false
	for _, f := range files {
		switch {
		case f == assets.Dir+"/"+assets.WorkFile:
			work = true
		case isCode(f):
			code = append(code, f)
		}
	}
	if len(code) == 0 || work {
		return "", false
	}
	if len(code) > 5 {
		code = append(code[:5], "…")
	}
	return fmt.Sprintf("這個 commit 改了 app 的程式（%s），%s/%s 沒有跟著改。\n"+
		"- 行為有變（畫面、按鍵、訊息、讀寫的檔、外部指令）：更新 work.yaml 那個 task 的 implementation 與 test，跟程式一起 commit。\n"+
		"- 行為沒變（重構、改註解、改測試以外的程式結構）：commit 訊息最後加一行 `%s`，再 commit 一次。",
		strings.Join(code, "、"), assets.Dir, assets.WorkFile, Trailer), true
}

// kept is a repo the harness keeps.
func kept(root string) bool {
	if fi, err := os.Stat(filepath.Join(root, assets.Dir)); err != nil || !fi.IsDir() {
		return false
	}
	b, err := os.ReadFile(filepath.Join(root, "go.mod"))
	return err == nil && strings.Contains(string(b), assets.TGO+" ")
}

// isCode is a file of the app's program: Go outside tests, and what cmd/ and
// internal/ hold besides tests and test data.
func isCode(f string) bool {
	if strings.HasSuffix(f, "_test.go") || strings.Contains("/"+f, "/testdata/") {
		return false
	}
	return strings.HasSuffix(f, ".go") || strings.HasPrefix(f, "cmd/") || strings.HasPrefix(f, "internal/")
}

// Run reads the hook input from r and, when the command is to be blocked,
// writes why to w and returns 2, the exit code that makes Claude Code block
// the call and show Claude the reason.
func Run(r io.Reader, w io.Writer) int {
	var in Input
	if err := json.NewDecoder(r).Decode(&in); err != nil || in.ToolName != "Bash" {
		return 0
	}
	if why, block := Commit(in, DefaultGit); block {
		fmt.Fprintln(w, why)
		return 2
	}
	return 0
}

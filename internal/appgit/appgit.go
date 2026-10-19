// Package appgit reads what judging needs of an app's git repo: its commit,
// its terminu-go version, and the files changed since an earlier commit.
package appgit

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Commit is the repo's HEAD, "" when it is not a git repo.
func Commit(dir string) string {
	out, err := exec.Command("git", "-C", dir, "rev-parse", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// Tgo is the terminu-go version go.mod requires, "" for none.
func Tgo(dir string) string {
	f, err := os.Open(filepath.Join(dir, "go.mod"))
	if err != nil {
		return ""
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(strings.TrimPrefix(strings.TrimSpace(sc.Text()), "require "))
		if len(fields) >= 2 && fields[0] == "github.com/vulcanshen/terminu-go" {
			return fields[1]
		}
	}
	return ""
}

// Changed is every file changed since the commit: committed after it, and
// not committed yet (untracked ones too).
func Changed(dir, since string) ([]string, error) {
	out, err := exec.Command("git", "-C", dir, "diff", "--name-only", since, "HEAD").Output()
	if err != nil {
		return nil, err
	}
	files := strings.Fields(string(out))
	status, err := exec.Command("git", "-C", dir, "status", "--porcelain").Output()
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(string(status), "\n") {
		if len(line) > 3 {
			files = append(files, strings.TrimSpace(line[3:]))
		}
	}
	return files, nil
}

// IsCode reports a file the app is built from: Go source, the module files,
// and what cmd/ and internal/ hold (embedded files too).
func IsCode(path string) bool {
	path = filepath.ToSlash(path)
	return strings.HasSuffix(path, ".go") || path == "go.mod" || path == "go.sum" ||
		strings.HasPrefix(path, "cmd/") || strings.HasPrefix(path, "internal/")
}

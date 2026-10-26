package hook

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// app is a repo the harness keeps.
func app(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module x\n\nrequire github.com/vulcanshen/terminu-go v0.5.1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "harness-assets"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

// fakeGit has staged, changed and new files.
func fakeGit(root string, staged, changed, added []string) Git {
	return Git{
		Root: func(string) (string, error) { return root, nil },
		Files: func(_ string, args ...string) ([]string, error) {
			switch strings.Join(args, " ") {
			case "diff --cached --name-only":
				return staged, nil
			case "diff --name-only":
				return changed, nil
			default:
				return added, nil
			}
		},
	}
}

func input(cwd, command string) Input {
	var in Input
	in.Cwd, in.ToolName, in.ToolInput.Command = cwd, "Bash", command
	return in
}

func TestCommit(t *testing.T) {
	root := app(t)
	code := []string{"internal/ui/list.go"}
	work := []string{"internal/ui/list.go", "harness-assets/work.yaml"}
	for _, c := range []struct {
		name                   string
		command                string
		staged, changed, added []string
		block                  bool
	}{
		{"not a commit", "go test ./...", code, nil, nil, false},
		{"commit-tree is not a commit", "git commit-tree abc", code, nil, nil, false},
		{"code without work.yaml", `git commit -m "Sort by year"`, code, nil, nil, true},
		{"code with work.yaml", `git commit -m "Sort by year"`, work, nil, nil, false},
		{"the trailer", "git commit -m \"Split the list\n\nTerminu-Work: unchanged\"", code, nil, nil, false},
		{"only tests", `git commit -m t`, []string{"internal/ui/list_test.go", "internal/ui/testdata/a.yaml"}, nil, nil, false},
		{"only docs", `git commit -m d`, []string{"README.md", "CHANGELOG.md"}, nil, nil, false},
		{"embedded file in internal", `git commit -m e`, []string{"internal/ui/help.txt"}, nil, nil, true},
		{"add then commit", "git add -A && git commit -m x", nil, code, nil, true},
		{"add then commit with work.yaml", "git add -A && git commit -m x", nil, code, []string{"harness-assets/work.yaml"}, false},
		{"commit -am", `git commit -am "x"`, nil, code, nil, true},
		{"unstaged code is not committed", `git commit -m x`, []string{"README.md"}, code, nil, false},
		{"git -C", `git -C sub commit -m x`, code, nil, nil, true},
	} {
		why, block := Commit(input(root, c.command), fakeGit(root, c.staged, c.changed, c.added))
		if block != c.block {
			t.Errorf("%s: block = %v (%s)", c.name, block, why)
		}
		if block && !strings.Contains(why, Trailer) {
			t.Errorf("%s: the reason does not say how to pass: %s", c.name, why)
		}
	}
}

func TestCommitTrailerInFile(t *testing.T) {
	root := app(t)
	msg := filepath.Join(root, "msg.txt")
	if err := os.WriteFile(msg, []byte("Split the list\n\nTerminu-Work: unchanged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git := fakeGit(root, []string{"internal/ui/list.go"}, nil, nil)
	if why, block := Commit(input(root, "git commit -F msg.txt"), git); block {
		t.Errorf("blocked with the trailer in -F's file: %s", why)
	}
}

func TestCommitNotKept(t *testing.T) {
	plain := t.TempDir()
	if err := os.WriteFile(filepath.Join(plain, "go.mod"), []byte("module x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git := fakeGit(plain, []string{"main.go"}, nil, nil)
	if _, block := Commit(input(plain, `git commit -m x`), git); block {
		t.Error("blocked in a repo the harness does not keep")
	}
}

func TestRun(t *testing.T) {
	var w bytes.Buffer
	if code := Run(strings.NewReader("not json"), &w); code != 0 {
		t.Errorf("bad input: exit %d", code)
	}
	if code := Run(strings.NewReader(`{"tool_name":"Read","tool_input":{}}`), &w); code != 0 {
		t.Errorf("another tool: exit %d", code)
	}
}

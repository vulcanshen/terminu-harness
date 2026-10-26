package tdp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// IDs is what an app can cite in tdp at a version: every ID written in a
// heading of principle, components and project (rules K4, principle P5, terms
// terms.3, sections panel.5), and every component and project file by name
// (confirm, input), which stands for the whole file.
type IDs struct {
	Sections map[string]string // ID -> the heading's title
	Files    map[string]string // name -> the file's path
}

// Has reports whether id is a section or a file in tdp.
func (ids IDs) Has(id string) bool {
	_, s := ids.Sections[id]
	_, f := ids.Files[id]
	return s || f
}

// Files lists the paths of every file in tdp at the tag.
func (s Source) Files(tag string) ([]string, error) {
	tag = Tag(tag)
	if _, ok := semver(tag); !ok {
		return nil, fmt.Errorf("%q is not a vX.Y.Z version", tag)
	}
	cached := ""
	if s.Cache != "" {
		cached = filepath.Join(s.Cache, tag, ".tree")
		if b, err := os.ReadFile(cached); err == nil {
			return strings.Fields(string(b)), nil
		}
	}
	var paths []string
	if s.Repo != "" {
		out, err := exec.Command("git", "-C", s.Repo, "ls-tree", "-r", "--name-only", tag).Output()
		if err != nil {
			return nil, fmt.Errorf("files of %s in %s: %w", tag, s.Repo, err)
		}
		paths = strings.Fields(string(out))
	} else {
		b, err := get("https://api.github.com/repos/" + GitHub + "/git/trees/" + tag + "?recursive=1")
		if err != nil {
			return nil, err
		}
		var tree struct {
			Tree []struct{ Path, Type string }
		}
		if err := json.Unmarshal(b, &tree); err != nil {
			return nil, fmt.Errorf("files of %s at %s: %w", GitHub, tag, err)
		}
		for _, t := range tree.Tree {
			if t.Type == "blob" {
				paths = append(paths, t.Path)
			}
		}
	}
	sort.Strings(paths)
	if cached != "" {
		if os.MkdirAll(filepath.Dir(cached), 0o755) == nil {
			_ = os.WriteFile(cached, []byte(strings.Join(paths, "\n")+"\n"), 0o644)
		}
	}
	return paths, nil
}

// IDs reads the IDs in tdp at the tag from the headings of the Traditional
// Chinese files, the original.
func (s Source) IDs(tag string) (IDs, error) {
	paths, err := s.Files(tag)
	if err != nil {
		return IDs{}, err
	}
	ids := IDs{Sections: map[string]string{}, Files: map[string]string{}}
	for _, p := range paths {
		if !citable(p) {
			continue
		}
		b, err := s.File(tag, p)
		if err != nil {
			return IDs{}, err
		}
		for id, title := range headingIDs(b) {
			ids.Sections[id] = title
		}
		if name := fileName(p); name != "" {
			ids.Files[name] = p
		}
	}
	return ids, nil
}

// citable is a file whose headings carry IDs.
func citable(p string) bool {
	if !strings.HasSuffix(p, "-zh_TW.md") {
		return false
	}
	for _, dir := range []string{"principle/", "components/", "project/"} {
		if strings.HasPrefix(p, dir) {
			return true
		}
	}
	return false
}

// fileName is the name a component or project file is cited by: the file's
// own name, or its folder's for a README (components, input, project).
func fileName(p string) string {
	if !strings.HasPrefix(p, "components/") && !strings.HasPrefix(p, "project/") {
		return ""
	}
	base := strings.TrimSuffix(path.Base(p), "-zh_TW.md")
	if base == "README" {
		return path.Base(path.Dir(p))
	}
	return base
}

var (
	headingRE = regexp.MustCompile(`^#{2,4} (\S+)\s+(.*)$`)
	idRE      = regexp.MustCompile(`^(?:[A-Z][0-9]+|[a-z][a-z0-9-]*\.[0-9]+)$`)
)

// headingIDs are the IDs at the start of the headings, outside code blocks.
func headingIDs(b []byte) map[string]string {
	ids := map[string]string{}
	inCode := false
	sc := bufio.NewScanner(bytes.NewReader(b))
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "```") {
			inCode = !inCode
			continue
		}
		if inCode {
			continue
		}
		m := headingRE.FindStringSubmatch(line)
		if m != nil && idRE.MatchString(m[1]) {
			ids[m[1]] = m[2]
		}
	}
	return ids
}
